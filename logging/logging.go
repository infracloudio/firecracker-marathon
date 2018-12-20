package logging

import (
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

var logger *logrus.Logger

func init() {
	var logLevel logrus.Level
	//defaultLevel := logrus.InfoLevel

	// currently need debugging log level
	logLevel = logrus.DebugLevel

	logger = logrus.New()
	logger.SetLevel(logLevel)
	logger.SetOutput(os.Stdout)

	envType, envTypeSet := os.LookupEnv("GO_ENV")
	if envTypeSet && envType == "production" {
		// Add the syslog hooks
		hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_DEBUG, "")
		if err == nil {
			logger.AddHook(hook)
		}
	}
}

// NewLogger returns a common logger instance to use in an application
func NewLogger() *logrus.Logger {
	return logger
}

// WithFunction returns a logger with a function field
func WithFunction(functionName string) *logrus.Entry {
	return logger.WithField("function", functionName)
}
