package vpanel

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"strings"
)

type Manager struct {
	stopCh      chan bool
	containerCh chan *Container
	containers  map[string]*Container
}

func ContainerTemplates() ([]string, error) {
	names, err := Config.TemplateDirEntries()
	for i, name := range names {
		names[i] = strings.TrimPrefix(name, "lxc-")
	}
	return names, err
}

func (m *Manager) CreateContainer(metadata ContainerMetadata) error {
	if metadata.IsValid(); metadata.Err != nil {
		return metadata.Err
	}
	container, err := newContainer(metadata, m)
	if err == nil {
		metadata.save()
		container.Create()
		m.containerCh <- container
	}
	return err
}

func (m *Manager) getContainer(id string, f func(*Container)) error {
	if container, ok := m.containers[id]; ok {
		f(container)
		return nil
	}
	return errors.New("Container with id " + id + " does not exist")
}

func (m *Manager) StartContainer(id string) error {
	return m.getContainer(id, func(container *Container) { container.Start() })
}

func (m *Manager) StopContainer(id string) error {
	return m.getContainer(id, func(container *Container) { container.Stop() })
}

func (m *Manager) FreezeContainer(id string) error {
	return m.getContainer(id, func(container *Container) { container.Freeze() })
}

func (m *Manager) managerLoop() {
	for {
		select {
		case c := <-m.containerCh:
			m.containers[c.Metadata.Id.String()] = c
		case <-m.stopCh:
			m.stopCh <- true
			return
		}
	}
}

func (m *Manager) Start() {
	go m.managerLoop()

	names, err := Config.ContainerDirEntries()
	if err != nil {
		Logger.Warnf("Failed to load container directory names: %v", err)
		return
	}

	for _, name := range names {
		if id := uuid.Parse(name); id != nil {
			/* Load metadata */
			metadata, err := loadMetadata(id)
			if err != nil {
				Logger.Warnf("Failed to load container %s metadata file: %v", name, err)
				continue
			} else {
				container, err := loadContainer(metadata, m)

				if err != nil {
					Logger.Warnf("Failed to load container %s: %v", name, err)
					continue
				}
				m.containerCh <- container
			}
		} else {
			Logger.Warnf("Failed to load container %s: The name does not appear to be a UUID", name)
		}
	}
}

func (m *Manager) Stop() {
	m.stopCh <- true
	<-m.stopCh
}

func NewManager() *Manager {
	manager := new(Manager)
	manager.stopCh = make(chan bool)
	manager.containerCh = make(chan *Container)
	manager.containers = make(map[string]*Container)
	return manager
}
