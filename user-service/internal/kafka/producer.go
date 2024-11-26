package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"user-service/config"

	"github.com/segmentio/kafka-go"
)

type NotificationMessage struct {
	UserID     int32  `json:"user_id"`
	FollowerID int32  `json:"follower_id"`
	Action     string `json:"action"` // "add", "accept", "reject", "request"
	Message    string `json:"message"`
}

func PublishNotification(ctx context.Context, topic string, notification NotificationMessage) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{conf.KAFKA_BROKERS},
		Topic:   conf.KAFKA_TOPIC,
	})

	defer func() {
		if err := producer.Close(); err != nil {
			fmt.Printf("error closing kafka writer: %v\n", err)
		}
	}()

	valueBytes, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error marshaling notification to JSON: %v", err)
	}

	msg := kafka.Message{
		Key:   []byte(strconv.Itoa(int(notification.UserID))), 
		Value: valueBytes,
	}
	if err := producer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("error writing message to kafka: %v", err)
	}

	return nil
}
