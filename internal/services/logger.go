package services

import (
	"fantomrocks-api/internal/common"
	"github.com/op/go-logging"
	"os"
)

// define new Logger interface
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

// Get pre-configured logger with stderr output and leveled filtering.
func NewLogger(cfg *common.Config) *logging.Logger {
	// prep the backend with configured formatting, use stderr for logging
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	format := logging.MustStringFormatter(cfg.LogFormat)
	formattedBackend := logging.NewBackendFormatter(backend, format)

	// make it leveled
	leveledBackend := logging.AddModuleLevel(formattedBackend)
	level, err := logging.LogLevel(cfg.LogLevel)
	if err != nil {
		level = logging.INFO
	}
	leveledBackend.SetLevel(level, "")

	// assign the backend and return the new logger
	logging.SetBackend(leveledBackend)
	return logging.MustGetLogger(cfg.AppName)
}
