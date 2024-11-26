package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"direct-service/config"
	"direct-service/internal/models"

	"github.com/segmentio/kafka-go"
)

type DirectMessageNotification struct {
	Message models.DirectMessage `json:"message"`
}

func PublishDirectMessage(ctx context.Context, topic string, directMessage models.DirectMessage) error {
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

	valueBytes, err := json.Marshal(DirectMessageNotification{Message: directMessage})
	if err != nil {
		return fmt.Errorf("error marshaling direct message to JSON: %v", err)
	}

	msg := kafka.Message{
		Key:   []byte(strconv.Itoa(int(directMessage.SenderID))),
		Value: valueBytes,
	}
	if err := producer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("error writing message to kafka: %v", err)
	}

	return nil
}