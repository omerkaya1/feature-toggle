package internal

import "log"

type BaseLogger struct{}

func NewBaseLogger() *BaseLogger {
	return &BaseLogger{}
}

func (b BaseLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (b BaseLogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (b BaseLogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (b BaseLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
