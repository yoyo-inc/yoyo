package sched

import (
	"github.com/robfig/cron/v3"
	"github.com/yoyo-inc/yoyo/common/logger"
)

var C *cron.Cron

func Setup() {
	C = cron.New(cron.WithLogger(&cronLogger{}))
	C.Start()
}

type cronLogger struct{}

func (*cronLogger) Info(msg string, keysAndValues ...interface{}) {
	logger.Debugf(msg, keysAndValues...)
}

func (*cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	logger.Errorf(msg, keysAndValues...)
}
