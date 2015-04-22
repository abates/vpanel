package server

import (
	"code.google.com/p/go-uuid/uuid"
	"gopkg.in/lxc/go-lxc.v2"
	"net"
	"strconv"
)

type CtStatus int

type ContainerMemory struct {
	Limit string `json:"limit"`
	Used  string `json:"used"`
}

type Container struct {
	lxc         *lxc.Container
	err         error
	Name        string   `json:"name"`
	Hostname    string   `json:"hostname"`
	IPAddresses []net.IP `json:"ip"`
	AutoStart   bool     `json:"autostart"`
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
	if c.err != nil {
		return nil
	}
	limit, c.err = c.lxc.MemoryLimit()
	if c.err != nil {
		return nil
	}
	containerMemory.Limit = limit.String()
	limit, c.err = c.lxc.MemoryUsage()
	if c.err != nil {
		return nil
	}
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
	items := make([]string)
	if c.err == nil {
		items := c.lxc.ConfigItem(key)
		for item := range items {
			items = append(items, item)
		}
	}
	return items
}

func GetContainer(id uuid.UUID) *Container {
	c := Container{}
	c.lxc, c.err = lxc.NewContainer(id.String(), Config.Dir.Name())
	c.Id = id
	c.Name = ""
	c.Hostname = c.ConfigItem("lxc.utsname")
	c.AutoStart, c.err = strconv.ParseBool(c.ConfigItem("lxc.start.auto"))
	c.IPAddresses = make([]net.IP)
	addAddresses = func(key) {
		for address := range c.ConfitItems(key) {
			ip := net.ParseIP(ip)
			if ip == nil {
				Logger.Warnf("Container %s: Invalid IP Address %s", address)
			} else {
				c.IPAddresses = append(c.IPAddresses, ip)
			}
		}
	}
	addAddresses("lxc.network.ipv4")
	addAddresses("lxc.network.ipv6")
	return &c
}

func GetContainers() ([]*Container, error) {
	containers := make([]*Container)
	if names, err := Config.Dir.Readdirnames; err != nil {
		return containers, err
	}

	for name := range names {
		if id := uuid.Parse(name); id != nil {
			containers = append(containers, GetContainer(id))
		}
	}
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
