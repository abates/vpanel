package server

import (
	"fmt"
	"os"
	"path/filepath"
)

type config struct {
	Dir *os.File
}

var Config config

func getEnv(key, def string) string {
	v := strings.TrimSpace(os.GetEnv(key))
	if len(v) == 0 {
		v = def
	}
	return v
}

type logger struct{}

func (l *logger) Debug(v ...interface{}) {
	log.Print("DEBUG", v)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v))
}

func (l *logger) Info(v ...interface{}) {
	log.Print("INFO", v)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v))
}

func (l *logger) Warn(v ...interface{}) {
	log.Print("WARN", v)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v))
}

func (l *logger) Error(v ...interface{}) {
	log.Print("ERROR", v)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v))
}

func (l *logger) Fatal(v ...interface{}) {
	log.Print("FATAL", v)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Fatal(fmt.Sprintf(format, v))
}

var Logger logger

func setDir() {
	dir, err := filepath.Abs(getEnv("VPANEL_DIR", "/var/lib/vpanel"))
	if err != nil {
		Logger.Fatalf("Failed to resolve vpanel directory: %v", err)
	}

	config.dir, err = os.Open(dir)
	if err != nil {
		if os.IsNotExist(err) {
			Logger.Fatalf("Failed to open %s: directory does not exist", dir)
		} else {
			Logger.Fatalf("Failed to open %s: %v", dir, err)
		}
	}

	if fi, err := config.dir.Stat(); err != nil {
		Logger.Fatalf("Failed to stat %s: %v", dir, err)
	}
	if !fi.IsDir() {
		Logger.Fatalf("%s is not a directory!", dir)
	}
}

func init() {
	Logger = logger{}
	Config = config{}
	setDir()
}
