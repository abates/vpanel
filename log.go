package vpanel

import (
	"fmt"
	"log"
	"os"
)

type logger struct{}

func (l logger) Debug(v ...interface{}) {
	l.Print(fmt.Sprintf("DEBUG %v", v...))
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v...))
}

func (l logger) Info(v ...interface{}) {
	l.Print(fmt.Sprintf("INFO %v", v...))
}

func (l logger) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

func (l logger) Warn(v ...interface{}) {
	l.Print(fmt.Sprintf("WARN %v", v...))
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v...))
}

func (l logger) Error(v ...interface{}) {
	l.Print(fmt.Sprintf("ERROR %v", v...))
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

func (l logger) Fatal(v ...interface{}) {
	l.Print(fmt.Sprintf("FATAL %v", v...))
	os.Exit(1)
}

func (l logger) Fatalf(format string, v ...interface{}) {
	l.Fatal(fmt.Sprintf(format, v...))
}

func (l logger) Print(v ...interface{}) {
	log.Print(v...)
}

var Logger logger
