package vpanel

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type configInterface interface {
	Set(string, interface{})
	Get(string, interface{})
}

type config struct {
	store map[string]interface{}
}

var Config config

type resolver func(string) (interface{}, error)

func getEnv(key, def string, r resolver) interface{} {
	value := strings.TrimSpace(os.Getenv(key))
	if len(value) == 0 {
		value = def
	}

	if r != nil {
		resolved, err := r(value)
		if err != nil {
			Logger.Fatalf("Failed to set %s to value %s: %s", key, value, err.Error())
		}
		return resolved
	}
	return value
}

func resolveDir(path string) (interface{}, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	dir, err := os.Open(path)
	defer dir.Close()

	if err != nil {
		return "", err
	}

	fi, err := dir.Stat()
	if err != nil {
		return "", err
	}

	if !fi.IsDir() {
		return "", errors.New(path + " is not a directory")
	}
	return path, nil
}

func (c config) dirEntries(key string) ([]string, error) {
	var path string
	c.Get(key, &path)

	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return []string{}, err
	}
	return dir.Readdirnames(0)
}

func (c config) ContainerDirEntries() ([]string, error) {
	return c.dirEntries("containerDir")
}

func (c config) TemplateDirEntries() ([]string, error) {
	return c.dirEntries("templateDir")
}

func (c config) Set(key string, value interface{}) {
	c.store[key] = value
}

func (c config) Get(key string, value interface{}) {
	r := reflect.ValueOf(value)
	e := r.Elem()
	e.Set(reflect.ValueOf(c.store[key]))
}

func parseDuration(value string) (interface{}, error) {
	duration, err := strconv.ParseInt(value, 10, 64)
	return time.Duration(duration), err
}

func init() {
	Config.store = make(map[string]interface{})
	Config.Set("containerDir", getEnv("VP_CONTAINER_DIR", "/var/lib/vpanel", resolveDir))
	Config.Set("templateDir", getEnv("VP_CONTAINER_DIR", "/usr/share/lxc/templates", resolveDir))
	Config.Set("listenPort", getEnv("VP_LISTEN_PORT", "3000", nil))
	Config.Set("hostMonitorInterval", getEnv("VP_HOST_MONITOR_INTERVAL", "60", parseDuration))
	Config.Set("containerMonitorInterval", getEnv("VP_CONTAINER_MONITOR_INTERVAL", "60", parseDuration))
}
