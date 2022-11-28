package services

import (
	"context"
	"fmt"
	"log/syslog"

	"github.com/segmentio/kafka-go"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/models"
)

var (
	syslogClients = make(map[int]*syslog.Writer)
	kafkaClients  = make(map[int]*kafka.Conn)
)

func PushAlert(push models.AlertPush, message string) bool {
	switch push.Type {
	case "syslog":
		return PushBySyslog(push, message)
	case "kafka":
		return PushByKafka(push, message)
	}

	return false

}

func PushBySyslog(push models.AlertPush, message string) bool {
	if _, ok := syslogClients[push.ID]; !ok {
		syslogClient, err := syslog.Dial(push.SyslogNetwork, fmt.Sprintf("%s:%d", push.SyslogAddr, push.SyslogPort), syslog.LOG_WARNING, push.SyslogTag)
		if err != nil {
			logger.Error(err)
			return false
		}

		syslogClients[push.ID] = syslogClient
	}

	syslogClient := syslogClients[push.ID]
	_, err := fmt.Fprint(syslogClient, message)
	if err != nil {
		logger.Error(err)
		return false
	}

	return true
}

func PushByKafka(push models.AlertPush, message string) bool {
	if _, ok := kafkaClients[push.ID]; !ok {
		kafkaClient, err := kafka.DialLeader(context.Background(), push.KafkaNetwork, fmt.Sprintf("%s:%d", push.SyslogAddr, push.SyslogPort), push.KafkaTopic, 0)
		if err != nil {
			logger.Error(err)
			return false
		}

		kafkaClients[push.ID] = kafkaClient
	}

	kafkaClient := kafkaClients[push.ID]
	_, err := fmt.Fprint(kafkaClient, message)
	if err != nil {
		logger.Error(err)
		return false
	}

	return true
}
