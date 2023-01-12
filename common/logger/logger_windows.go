//go:build windows

package logger

import (
	"strings"

	logrus_rollingfile_hook "github.com/gfremex/logrus-rollingfile-hook"
	"github.com/sirupsen/logrus"
)

var logger logrus.Logger

type SyslogOption struct {
	Network string
	Addr    string
}

type RollingFileOption struct {
	BaseDir string `mapstructure:"base_dir"`
}

// Options contains logger options
type Options struct {
	Service     string
	Level       string
	Syslog      *SyslogOption
	RollingFile *RollingFileOption
}

// Setup setups logger
func Setup(options Options) {
	logger = *logrus.New()
	logger.ReportCaller = true
	logger.CallerSkip = 1

	service := "default"
	if options.Service != "" {
		service = options.Service
	}
	if options.Level != "" {
		level := logrus.InfoLevel
		switch strings.ToLower(options.Level) {
		case "panic":
			level = logrus.PanicLevel
		case "fatal":
			level = logrus.FatalLevel
		case "error":
			level = logrus.ErrorLevel
		case "warn":
			level = logrus.WarnLevel
		case "info":
		case "debug":
			level = logrus.DebugLevel
		case "trace":
			level = logrus.TraceLevel
		}

		logger.SetLevel(level)
	}

	logger.SetFormatter(&TextFormatter{
		Service: service,
	})

	if options.Syslog != nil {
		SetupSyslog(service, options.Syslog)
	}

	if options.RollingFile != nil {
		SetupRollingFile(service, options.RollingFile)
	}
}

func SetupSyslog(service string, option *SyslogOption) {
}

func SetupRollingFile(service string, option *RollingFileOption) {
	filenamePattern := strings.TrimSuffix(option.BaseDir, "/") + "/" + service + "-%Y-%m-%d.log"
	hook, err := logrus_rollingfile_hook.NewTimeBasedRollingFileHook("rolling-file", logrus.AllLevels, &TextFormatter{Service: service}, filenamePattern)
	if err != nil {
		logger.Error(err)
		return
	} else {
		logger.AddHook(hook)
	}
}

// Panic logs PanicLevel message
func Panic(args ...interface{}) {
	logger.Panicln(args...)
}

// Panicf logs PanicLevel message by format
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// Fatal log FatalLevel message
func Fatal(args ...interface{}) {
	logger.Fatalln(args...)
}

// Fatalf log FatalLevel message by format
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Error logs ErrorLevel message
func Error(args ...interface{}) {
	logger.Errorln(args...)
}

// Errorf logs ErrorLevel message by format
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Warn logs WarnLevel message
func Warn(args ...interface{}) {
	logger.Warnln(args...)
}

// Warnf logs WarnLevel message
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Info logs InfoLevel message
func Info(args ...interface{}) {
	logger.Infoln(args...)
}

// Infof logs InfoLevel message by format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Debug logs DebugLevel message
func Debug(args ...interface{}) {
	logger.Debugln(args...)
}

// Debugf logs DebugLevel message by format
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Trace logs TraceLevel message
func Trace(args ...interface{}) {
	logger.Traceln(args...)
}

// Tracef logs TraceLevel message by format
func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func GetLogger() *logrus.Logger {
	return &logger
}
