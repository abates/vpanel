package vpanel

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"strings"
	"sync"
)

type ContainerNotFoundError error
type DuplicateContainerIdError error

type Manager struct {
	stopCh          chan bool
	containersMutex sync.RWMutex
	containers      map[string]*Container
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
	container, err := newContainer(metadata)
	if err == nil {
		metadata.save()
		err = container.Create()
		if err == nil {
			err = m.AddContainer(container)
		}
	}
	return err
}

func (m *Manager) AddContainer(container *Container) error {
	m.containersMutex.Lock()
	defer m.containersMutex.Unlock()
	id := container.Metadata.Id.String()
	if _, found := m.containers[id]; !found {
		m.containers[id] = container
		return nil
	}
	return DuplicateContainerIdError(fmt.Errorf("A container with id %s alredy exists", id))
}

func (m *Manager) GetContainer(id string) (*Container, error) {
	m.containersMutex.RLock()
	defer m.containersMutex.RUnlock()
	if container, found := m.containers[id]; found {
		return container, nil
	}
	return nil, ContainerNotFoundError(fmt.Errorf("Container %s was not found", id))
}

func (m *Manager) DestroyContainer(id string) error {
	m.containersMutex.Lock()
	defer m.containersMutex.Unlock()
	if container, found := m.containers[id]; found {
		err := container.destroy()
		if err == nil {
			delete(m.containers, id)
		}
		return err
	}
	return ContainerNotFoundError(fmt.Errorf("Container %s was not found", id))
}

func NewManager() *Manager {
	manager := new(Manager)
	manager.stopCh = make(chan bool)
	manager.containers = make(map[string]*Container)

	names, err := Config.ContainerDirEntries()
	if err != nil {
		Logger.Warnf("Failed to load container directory names: %v", err)
		return nil
	}

	for _, name := range names {
		if id := uuid.Parse(name); id != nil {
			/* Load metadata */
			metadata, err := loadMetadata(id)
			if err != nil {
				Logger.Warnf("Failed to load container %s metadata file: %v", name, err)
				continue
			} else {
				container, err := loadContainer(metadata)

				if err != nil {
					Logger.Warnf("Failed to load container %s: %v", name, err)
					continue
				}
				manager.AddContainer(container)
			}
		} else {
			Logger.Warnf("Failed to load container %s: The name does not appear to be a UUID", name)
		}
	}

	return manager
}
