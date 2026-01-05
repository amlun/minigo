package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Init configures a global logrus logger.
func Init(level string) {
	if logger != nil {
		return
	}
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{})
	switch level {
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "error":
		l.SetLevel(logrus.ErrorLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}
	logger = l
}

// L returns the global logger.
func L() *logrus.Logger {
	if logger == nil {
		Init("info")
	}
	return logger
}
