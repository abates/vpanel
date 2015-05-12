package vpanel

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"gopkg.in/lxc/go-lxc.v2"
	"net"
	//"strconv"
	"strings"
)

type ContainerMetadata struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Hostname    string    `json:"hostname"`
	IPAddresses []net.IP  `json:"ip"`
	Template    string    `json:"template"`
	AutoStart   bool      `json:"autostart"`
}

func (c *ContainerMetadata) Validate() error {
	v := &ValidationError{}
	v.AddError(c.ValidateName())
	v.AddError(c.ValidateHostname())
	v.AddError(c.ValidateTemplate())
	if len(v.Errors) > 0 {
		return v
	}
	return nil
}

func (c *ContainerMetadata) ValidateName() error {
	if len(c.Name) == 0 {
		return errors.New("Name cannot be an empty string")
	}
	return nil
}

func (c *ContainerMetadata) ValidateHostname() error {
	if len(c.Hostname) == 0 {
		return errors.New("Hostname cannot be an empty string")
	}

	/* TODO: make sure this is working for unicode code points instead of byte length */
	for _, part := range strings.Split(c.Hostname, ".") {
		if len(part) > 63 {
			return errors.New("Hostname label " + part + " is greater than 63 characters")
		}
	}

	if len(c.Hostname) > 253 {
		return errors.New("Hostname cannot be more than 253 characters")
	}
	return nil
}

func (c *ContainerMetadata) ValidateTemplate() error {
	templates, err := Containers.Templates()
	if err != nil {
		return err
	}

	for _, template := range templates {
		if c.Template == template {
			return nil
		}
	}
	return errors.New(c.Template + " is not a valid template")
}

func NewContainerMetadata() *ContainerMetadata {
	c := ContainerMetadata{
		Id: uuid.NewRandom(),
	}
	return &c
}

type Container struct {
	lxc      *lxc.Container
	metadata *ContainerMetadata
}

/*type CtStatus int

type ContainerMemory struct {
	Limit string `json:"limit"`
	Used  string `json:"used"`
}

func (c *Container) Status() string {
	if c.err != nil {
		return "ERROR"
	}

	return c.lxc.State().String()
}

func (c *Container) Memory() ContainerMemory {
	var limit lxc.ByteSize
	containerMemory := ContainerMemory{}
	limit, c.err = c.lxc.MemoryLimit()
	containerMemory.Limit = limit.String()
	limit, c.err = c.lxc.MemoryUsage()
	containerMemory.Used = limit.String()
	return containerMemory
}

func (c *Container) Error() string {
	if c.err != nil {
		return c.err.Error()
	}
	return ""
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
	if c.err == nil {
		lxcItems := c.lxc.ConfigItem(key)
		for _, item := range lxcItems {
			items = append(items, item)
		}
	}
	return items
}

func (c *Container) Start() {
	if c.err == nil {
		state := c.lxc.State()
		if state == lxc.STOPPED {
			c.err = c.lxc.Start()
		} else if state == lxc.FROZEN {
			c.err = c.lxc.Unfreeze()
		}
	}
}

func (c *Container) Stop() {
	if c.err == nil && c.lxc.Running() {
		c.err = c.lxc.Stop()
	}
}

func (c *Container) Freeze() {
	if c.err == nil && c.lxc.Running() {
		c.err = c.lxc.Freeze()
	}
}
*/

type containers map[string]*Container

var Containers containers

func (c containers) Templates() ([]string, error) {
	names, err := Config.TemplateDirEntries()
	for i, name := range names {
		names[i] = strings.TrimPrefix(name, "lxc-")
	}
	return names, err
}

func (c containers) Create(metadata *ContainerMetadata) (*Container, error) {
	var err error
	container := &Container{
		metadata: metadata,
	}
	container.lxc, err = lxc.NewContainer(metadata.Id.String(), Config["containerDir"])
	if err == nil {
		c[metadata.Id.String()] = container
		println("Creating container")
		err = container.lxc.Create(lxc.TemplateOptions{
			Template: metadata.Template,
			Backend:  lxc.Directory,
		})
		println("Done")
	}
	return container, err
}

func (c containers) Get(id uuid.UUID) *ContainerMetadata {
	container := ContainerMetadata{}
	/*container.lxc, container.err = lxc.NewContainer(id.String(), Config.Dir.Name())
	container.Id = id
	container.Name = ""
	container.Hostname = container.ConfigItem("lxc.utsname")
	container.AutoStart, container.err = strconv.ParseBool(container.ConfigItem("lxc.start.auto"))
	container.IPAddresses = make([]net.IP, 0)
	addAddresses := func(key string) {
		for _, address := range container.ConfigItems(key) {
			ip := net.ParseIP(address)
			if ip == nil {
				Logger.Warnf("Container %s: Invalid IP Address %s", address)
			} else {
				container.IPAddresses = append(container.IPAddresses, ip)
			}
		}
	}
	addAddresses("lxc.network.ipv4")
	addAddresses("lxc.network.ipv6")*/
	return &container
}

func (c containers) All() ([]*ContainerMetadata, error) {
	containers := make([]*ContainerMetadata, 0)
	names, err := Config.ContainerDirEntries()
	for _, name := range names {
		if id := uuid.Parse(name); id != nil {
			containers = append(containers, Containers.Get(id))
		}
	}
	return containers, err
}

func init() {
	Containers = make(containers)
}
