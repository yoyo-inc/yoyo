package logger

import (
	"github.com/sirupsen/logrus"
	logrusSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

var logger logrus.Logger

type SyslogOption struct {
	Network  string
	Addr     string
	Priority syslog.Priority
}

// Options contains logger options
type Options struct {
	Service string
	Syslog  *SyslogOption
}

// Setup setups logger
func Setup(options Options) {
	logger = *logrus.New()

	service := "default"
	if options.Service != "" {
		service = options.Service
	}

	logger.SetFormatter(&TextFormatter{
		Service: service,
	})

	if options.Syslog != nil {
		SetupSyslog(service, options.Syslog)
	}
}

func SetupSyslog(service string, option *SyslogOption) {
	hook, err := logrusSyslog.NewSyslogHook(option.Network, option.Addr, option.Priority, service)
	if err != nil {
		logger.Error(err)
		return
	} else {
		logger.AddHook(hook)
	}
}

// Info logs INFO_LEVEL message
func Info(args ...interface{}) {
	logger.Infoln(args...)
}

// Infof logs INFO_LEVEL message by format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn logs WARN_LEVEL message
func Warn(args ...interface{}) {
	logger.Warnln(args...)
}

// Warnf logs WARN_LEVEL message
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error logs ERROR_LEVEL message
func Error(args ...interface{}) {
	logger.Errorln(args...)
}

// Errorf logs ERROR_LEVEL message by format
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Panic logs PANIC_LEVEL message
func Panic(args ...interface{}) {
	logger.Panicln(args...)
}

// Panicf logs PANIC_LEVEL_LEVEL message by format
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func GetLogger() *logrus.Logger {
	return &logger
}
