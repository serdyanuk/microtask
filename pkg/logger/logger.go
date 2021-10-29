package logger

import (
	"os"

	"github.com/serdyanuk/microtask/config"
	"github.com/sirupsen/logrus"
)

var defaultLogger = newLogger(config.Get())

type Logger struct {
	*logrus.Logger
}

func Get() *Logger {
	return defaultLogger
}

func newLogger(conifg *config.Config) *Logger {
	var log = logrus.New()

	log.SetOutput(os.Stdout)

	return &Logger{
		Logger: log,
	}
}
