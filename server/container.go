package api

import (
	"code.google.com/p/go-uuid/uuid"
	"gopkg.in/lxc/go-lxc.v2"
	"net"
)

type CtStatus string

const (
	CtStatusCreating CtStatus = "creating"
	CtStatusRunning  CtStatus = "running"
	CtStatusStopped  CtStatus = "stopped"
	CtStatusFrozen   CtStatus = "frozen"
)

type ContainerMemory struct {
	Assigned uint64 `json:"assigned"`
	Used     uint64 `json:"used"`
}

type Container struct {
	lxc       lxc.Container
	Status    CtStatus        `json:"status"`
	Name      string          `json:"name"`
	Hostname  string          `json:"hostname"`
	IP        net.IP          `json:"ip"`
	Memory    ContainerMemory `json:"memory"`
	AutoStart bool            `json:"autostart"`
	Error     error           `json:"error"`
}

type Containers struct{}

func GetContainer(id uuid.UUID) *Container {
}

func GetContainers() []*Container {
}

func (c *Container) Refresh() error {
}

func (c *Container) Start() error {
}

func (c *Container) Stop() error {
}

func (c *Container) Freeze() error {
}
