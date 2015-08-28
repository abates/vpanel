package vpanel

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/cloudfoundry/gosigar"
	"gopkg.in/lxc/go-lxc.v2"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ContainerState string

var (
	ContainerError    = ContainerState("Error")
	ContainerCreating = ContainerState("Creating")
	ContainerStopping = ContainerState("Stopping")
	ContainerStopped  = ContainerState("Stopped")
	ContainerStarting = ContainerState("Starting")
	ContainerRunning  = ContainerState("Running")
	ContainerFreezing = ContainerState("Freezing")
	ContainerFrozen   = ContainerState("Frozen")
)

type Action int

const (
	ActionCreate Action = iota
	ActionStart
	ActionStop
	ActionFreeze
)

type ContainerMetadata struct {
	Id          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Hostname    string         `json:"hostname"`
	IPAddresses []net.IP       `json:"ip"`
	Template    string         `json:"template"`
	AutoStart   bool           `json:"autostart"`
	State       ContainerState `json:"state"`
	Err         error          `json:"-"`
}

func (c *ContainerMetadata) Validate(expression bool, err error) {
	if c.Err != nil {
		return
	} else if !expression {
		c.Err = err
	}
}

func (c *ContainerMetadata) ValidateName() {
	c.Validate(len(c.Name) > 0, errors.New("Name cannot be an empty string"))
}

func (c *ContainerMetadata) ValidateHostname() {
	c.Validate(len(c.Hostname) > 0, errors.New("Hostname cannot be an empty string"))

	/* TODO: make sure this is working for unicode code points instead of byte length */
	for _, part := range strings.Split(c.Hostname, ".") {
		c.Validate(len(part) < 64, errors.New("Hostname label "+part+" is greater than 63 characters"))
	}

	c.Validate(len(c.Hostname) < 254, errors.New("Hostname cannot be more than 253 characters"))
}

func (c *ContainerMetadata) ValidateTemplate() {
	templates, err := ContainerTemplates()
	c.Validate(err != nil, err)

	found := false
	for _, template := range templates {
		if c.Template == template {
			found = true
			break
		}
	}
	c.Validate(found, errors.New(c.Template+" is not a valid template"))
}

func (c *ContainerMetadata) IsValid() bool {
	c.ValidateName()
	c.ValidateHostname()
	c.ValidateTemplate()
	return c.Err == nil
}

func (c ContainerMetadata) save() error {
	if c.IsValid(); c.Err != nil {
		return c.Err
	}
	var containerDir string
	Config.Get("containerDir", &containerDir)
	file, err := os.OpenFile(containerDir+"/"+c.Id.String()+"/config.json", os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	return enc.Encode(&c)
}

func NewContainerMetadata() ContainerMetadata {
	c := ContainerMetadata{
		Id: uuid.NewRandom(),
	}
	return c
}

func loadMetadata(id uuid.UUID) (ContainerMetadata, error) {
	metadata := ContainerMetadata{}
	var containerDir string
	Config.Get("containerDir", &containerDir)
	file, err := os.Open(containerDir + "/" + id.String() + "/config.json")
	if err == nil {
		dec := json.NewDecoder(file)
		err = dec.Decode(&metadata)
	}

	return metadata, err
}

type ContainerStats struct {
	Cpu       float64     `json:"cpu"`
	Memory    MemoryStats `json:"memory"`
	systemCpu sigar.Cpu
	procCpu   int64
}

type Container struct {
	lxc       *lxc.Container
	Metadata  ContainerMetadata
	manager   *Manager
	Stats     ContainerStats
	monitorCh chan chan ContainerStats
	actionCh  chan Action
	stateCh   chan ContainerState
	stopCh    chan bool
	err       error
}

func newContainer(metadata ContainerMetadata, manager *Manager) (*Container, error) {
	var err error
	container := Container{
		Metadata: metadata,
		manager:  manager,
	}

	var containerDir string
	Config.Get("containerDir", &containerDir)

	container.lxc, err = lxc.NewContainer(metadata.Id.String(), containerDir)
	return &container, err
}

func loadContainer(metadata ContainerMetadata, manager *Manager) (*Container, error) {
	container, err := newContainer(metadata, manager)
	if err != nil {
		Logger.Warnf("Failed to create container: %v", err)
		return nil, err
	}

	hostname := container.ConfigItem("lxc.utsname")
	autoStart, err := strconv.ParseBool(container.ConfigItem("lxc.start.auto"))
	if err != nil {
		Logger.Warnf("Invalid value for lxc.start.auto when loading %v.  Expecting true or false", metadata.Id)
	}

	ipAddresses := make([]net.IP, 0)
	addAddresses := func(key string) {
		for _, address := range container.ConfigItems(key) {
			ip := net.ParseIP(address)
			if ip == nil {
				Logger.Warnf("Container %s: Invalid IP Address %s", address)
			} else {
				ipAddresses = append(container.Metadata.IPAddresses, ip)
			}
		}
	}
	addAddresses("lxc.network.ipv4")
	addAddresses("lxc.network.ipv6")

	metadata.Hostname = hostname
	metadata.AutoStart = autoStart
	metadata.IPAddresses = ipAddresses

	return container, nil
}

func (c *Container) ConfigItem(key string) string {
	items := c.ConfigItems(key)
	if len(items) > 0 {
		return items[0]
	}
	return ""
}

func (c *Container) ConfigItems(key string) []string {
	items := make([]string, 0)
	lxcItems := c.lxc.ConfigItem(key)
	for _, item := range lxcItems {
		items = append(items, item)
	}
	return items
}

func (c *Container) actionLoop() {
	var wg sync.WaitGroup
	for {
		select {
		case action := <-c.actionCh:
			wg.Add(1)
			go func() {
				switch action {
				case ActionCreate:
					c.create()
				case ActionStart:
					c.start()
				case ActionStop:
					c.stop()
				case ActionFreeze:
					c.freeze()
				default:
					Logger.Warnf("Don't know how to handle action %v", action)
				}
				wg.Done()
			}()
		case state := <-c.stateCh:
			c.updateState(state)
		case <-c.stopCh:
			wg.Wait()
			c.stopCh <- true
			return
		}
	}
}

func (c *Container) monitorLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			oldSystemCpu := c.Stats.systemCpu
			c.Stats.systemCpu.Get()
			procCpu, err := c.lxc.CPUStats()
			if err != nil {
				Logger.Warnf("Failed to retrieve cpu stats for container %v: %v", c.Metadata.Id, err)
			}
			c.Stats.procCpu = procCpu["user"] + procCpu["system"]
			delta := c.Stats.systemCpu.Delta(oldSystemCpu)
			c.Stats.Cpu = utilize(uint64(c.Stats.procCpu), (&delta).Total())
		case respCh := <-c.monitorCh:
			respCh <- c.Stats
		case <-c.stopCh:
			ticker.Stop()
			return
		}
	}
}

func (c *Container) updateState(state ContainerState) {
	Logger.Debugf("Updating container %s state to %v", c.Metadata.Id, state)
	if c.err == nil {
		c.Metadata.State = state
	} else {
		Logger.Warnf("Failed to update container %s state from %v to %v: %v", c.Metadata.Id, c.Metadata.State, state, c.err)
		c.Metadata.State = ContainerError
	}
}

func (c *Container) Create() {
	c.actionCh <- ActionCreate
}

func (c *Container) create() {
	c.stateCh <- ContainerCreating
	c.err = c.lxc.Create(lxc.TemplateOptions{
		Template: c.Metadata.Template,
		Backend:  lxc.Directory,
	})
	c.stateCh <- ContainerStopped
}

func (c *Container) Start() {
	c.actionCh <- ActionStart
}

func (c *Container) start() {
	state := c.lxc.State()
	if state == lxc.STOPPED || state == lxc.FROZEN {
		c.stateCh <- ContainerStarting
		if state == lxc.STOPPED {
			c.err = c.lxc.Start()
			var interval time.Duration
			Config.Get("VP_CONTAINER_MONITOR_INTERVAL", &interval)
			go c.monitorLoop(interval)
		} else if state == lxc.FROZEN {
			c.err = c.lxc.Unfreeze()
		}
		c.stateCh <- ContainerRunning
	}
}

func (c *Container) Stop() {
	c.actionCh <- ActionStop
}

func (c *Container) stop() {
	if c.lxc.Running() {
		c.stateCh <- ContainerStopping
		c.err = c.lxc.Stop()
		c.stopCh <- true
		c.stateCh <- ContainerStopped
	}
}

func (c *Container) Freeze() {
	c.actionCh <- ActionFreeze
}

func (c *Container) freeze() {
	if c.lxc.Running() {
		c.stateCh <- ContainerFreezing
		c.err = c.lxc.Freeze()
		c.stopCh <- true
		c.stateCh <- ContainerFrozen
	}
}
