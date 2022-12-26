package goerrors

import (
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

var instance atomic.Value

func Setup(format logrus.Formatter, level logrus.Level) error {
	log := logrus.New()
	log.Formatter = format
	log.Level = level
	instance.Store(log)
	return nil
}

func Log() *logrus.Logger {
	if logger, ok := instance.Load().(*logrus.Logger); ok && logger != nil {
		return logger
	}

	return logrus.StandardLogger()
}
