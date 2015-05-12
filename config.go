package vpanel

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type config map[string]string

var Config config

func getEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if len(v) == 0 {
		v = def
	}
	return v
}

func resolveDir(path string) (string, error) {
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

func dirEntries(path string) ([]string, error) {
	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return []string{}, err
	}
	return dir.Readdirnames(0)
}

func (c config) ContainerDirEntries() ([]string, error) {
	return dirEntries(c["containerDir"])
}

func (c config) TemplateDirEntries() ([]string, error) {
	return dirEntries(c["templateDir"])
}

type resolver func(string) (string, error)

func (c config) set(k, v string, r resolver) {
	if r != nil {
		v, err := r(v)
		if err != nil {
			Logger.Fatalf("Failed to set %s to value %s: %s", k, v, err.Error())
		}
	}
	c[k] = v
}

func init() {
	Config = make(config)
	Config.set("containerDir", getEnv("VP_CONTAINER_DIR", "/var/lib/vpanel"), resolveDir)
	Config.set("templateDir", getEnv("VP_CONTAINER_DIR", "/usr/share/lxc/templates"), resolveDir)
	Config.set("listenPort", getEnv("VP_LISTEN_PORT", "3000"), nil)
}
