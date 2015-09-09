package vpanel

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"gopkg.in/lxc/go-lxc.v2"
	"net"
	"strconv"
	"strings"
	"time"
)

type Action struct {
	respCh chan error
}

type CreateAction Action
type StartAction Action
type StopAction Action
type FreezeAction Action
type DestroyAction Action

type ContainerStats struct {
	Cpu       float64     `json:"cpu"`
	Memory    MemoryStats `json:"memory"`
	systemCpu sigar.Cpu
	procCpu   int64
}

type Container struct {
	lxc       *lxc.Container
	Metadata  ContainerMetadata
	Stats     ContainerStats
	monitorCh chan chan ContainerStats
	metaCh    chan chan ContainerMetadata
	actionCh  chan interface{}
	stopCh    chan bool
}

type InvalidStateErr error

func newContainer(metadata ContainerMetadata) (*Container, error) {
	var err error
	container := Container{
		Metadata: metadata,
	}

	var containerDir string
	Config.Get("containerDir", &containerDir)

	container.lxc, err = lxc.NewContainer(metadata.Id.String(), containerDir)
	return &container, err
}

func loadContainer(metadata ContainerMetadata) (*Container, error) {
	container, err := newContainer(metadata)
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
	for {
		select {
		case i := <-c.actionCh:
			switch action := i.(type) {
			case CreateAction:
				action.respCh <- c.create()
			case StartAction:
				c.start()
			case StopAction:
				c.stop()
			case FreezeAction:
				c.freeze()
			case DestroyAction:
				c.destroy()
			default:
				Logger.Warnf("Don't know how to handle action %v", action)
			}
		case <-c.stopCh:
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
			switch c.lxc.State() {
			case lxc.STOPPED:
				c.Metadata.State = ContainerStopped
			case lxc.STARTING:
				c.Metadata.State = ContainerStarting
			case lxc.STOPPING:
				c.Metadata.State = ContainerStopping
			case lxc.ABORTING:
			case lxc.FREEZING:
				c.Metadata.State = ContainerFreezing
			case lxc.FROZEN:
				c.Metadata.State = ContainerFrozen
			case lxc.THAWED:
			}
		case respCh := <-c.monitorCh:
			respCh <- c.Stats
		case respCh := <-c.metaCh:
			respCh <- c.Metadata
		case <-c.stopCh:
			ticker.Stop()
			return
		}
	}
}

func (c *Container) Create() error {
	action := CreateAction{make(chan error)}
	c.actionCh <- action
	return <-action.respCh
}

func (c *Container) create() error {
	c.Metadata.State = ContainerCreating
	return c.lxc.Create(lxc.TemplateOptions{
		Template: c.Metadata.Template,
		Backend:  lxc.Directory,
	})
}

func (c *Container) Start() error {
	action := StartAction{make(chan error)}
	c.actionCh <- action
	return <-action.respCh
}

func (c *Container) start() (err error) {
	state := c.lxc.State()
	if state == lxc.STOPPED || state == lxc.FROZEN {
		if state == lxc.STOPPED {
			err = c.lxc.Start()
			var interval time.Duration
			Config.Get("containerMonitorInterval", &interval)
			go c.monitorLoop(interval)
		} else if state == lxc.FROZEN {
			err = c.lxc.Unfreeze()
		}
	} else {
		err = InvalidStateErr(fmt.Errorf("cannot start when the container is %s", strings.ToLower(state.String())))
	}
	return err
}

func (c *Container) Stop() error {
	action := StopAction{make(chan error)}
	c.actionCh <- action
	return <-action.respCh
}

func (c *Container) stop() error {
	if c.lxc.Running() {
		err := c.lxc.Stop()
		c.stopCh <- true
		return err
	}

	return InvalidStateErr(fmt.Errorf("cannot stop then container unless it is running"))
}

func (c *Container) Freeze() error {
	action := FreezeAction{make(chan error)}
	c.actionCh <- action
	return <-action.respCh
}

func (c *Container) freeze() error {
	if c.lxc.Running() {
		err := c.lxc.Freeze()
		c.stopCh <- true
		return err
	}
	return InvalidStateErr(fmt.Errorf("cannot freeze the container unless it is running"))
}

func (c *Container) Destroy() error {
	action := DestroyAction{make(chan error)}
	c.actionCh <- action
	return <-action.respCh
}

func (c *Container) destroy() error {
	if c.lxc.Running() {
		return InvalidStateErr(fmt.Errorf("cannot destroy a running container!"))
	}
	c.stopCh <- true
	return nil
}
