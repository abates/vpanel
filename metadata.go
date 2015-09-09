package vpanel

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"net"
	"os"
	"strings"
)

type ValidationError error

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
