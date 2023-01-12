//go:build windows

package services

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/models"
)

var (
	kafkaClients = make(map[int]*kafka.Conn)
)

func PushAlert(push models.AlertPush, message string) bool {
	switch push.Type {
	case "kafka":
		return PushByKafka(push, message)
	}

	return false

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
