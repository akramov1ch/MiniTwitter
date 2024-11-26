						package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"notification-service/config"
	"github.com/segmentio/kafka-go"
	"notification-service/proto"
	"notification-service/internal/websocket"
)

func PublishFollowNotification(ctx context.Context, notification *proto.FollowNotification) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{conf.KAFKA_BROKERS},
		Topic:   "follows",
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
		Key:   []byte(strconv.Itoa(int(notification.UserId))), 
		Value: valueBytes,
	}
	if err := producer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("error writing message to kafka: %v", err)
	}

	websocket.BroadcastMessage(valueBytes)

	return nil
}

func PublishLikeNotification(ctx context.Context, notification *proto.LikeNotification) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{conf.KAFKA_BROKERS},
		Topic:   "likes",
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
		Key:   []byte(strconv.Itoa(int(notification.UserId))), 
		Value: valueBytes,
	}
	if err := producer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("error writing message to kafka: %v", err)
	}

	websocket.BroadcastMessage(valueBytes)

	return nil
}

func PublishCommentNotification(ctx context.Context, notification *proto.CommentNotification) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{conf.KAFKA_BROKERS},
		Topic:   "comments",
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
		Key:   []byte(strconv.Itoa(int(notification.UserId))), 
		Value: valueBytes,
	}
	if err := producer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("error writing message to kafka: %v", err)
	}

	websocket.BroadcastMessage(valueBytes)

	return nil
}

func PublishDirectMessage(ctx context.Context, notification *proto.DirectMessage) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{conf.KAFKA_BROKERS},
		Topic:   "directs",
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
		Key:   []byte(strconv.Itoa(int(notification.ReceiverId))), 
		Value: valueBytes,
	}
	if err := producer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("error writing message to kafka: %v", err)
	}

	websocket.BroadcastMessage(valueBytes)

	return nil
}