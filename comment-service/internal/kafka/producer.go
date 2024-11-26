package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"comment-service/config"

	"github.com/segmentio/kafka-go"
)

type NotificationMessage struct {
	UserID     int32  `json:"user_id"`
	TweetID    int32  `json:"tweet_id"`
	CommenterID int32  `json:"commenter_id"`
	Action     string `json:"action"`   
	Message    string `json:"message"`
}

func PublishCommentNotification(ctx context.Context, tweetOwnerID int32, tweetID int32, commenterID int32, commentContent string) error { // o'zgartirildi
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

	notification := NotificationMessage{
		UserID:     tweetOwnerID, // tweet egasi
		TweetID:    tweetID,       // yangi comment qo'shilgan tweet
		CommenterID: commenterID,  // comment yozgan foydalanuvchi
		Action:     "add",         // harakat turi
		Message:    fmt.Sprintf("Yangi comment qo'shildi: %s", commentContent), // xabar
	}

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