package logger

import "github.com/sirupsen/logrus"

var logger logrus.Logger

// Setup setups logger
func Setup() {
	logger = *logrus.New()

	logger.SetFormatter(&TextFormatter{
		Service: "default",
	})
}

// Info logs INFO_LEVEL message
func Info(args ...interface{}) {
	logger.Infoln(args...)
}

// Infof logs INFO_LEVEL message by format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
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
