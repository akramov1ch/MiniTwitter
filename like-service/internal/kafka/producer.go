package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"like-service/config"

	"github.com/segmentio/kafka-go"
)

type LikeNotificationMessage struct {
	UserID      int32  `json:"user_id"`      
	CommenterID int32  `json:"commenter_id"` 
	TweetOwnerID int32  `json:"tweet_owner_id"` 
	Action      string `json:"action"`        
	Message     string `json:"message"`
	TweetID     int32  `json:"tweet_id"`      
}

func PublishLikeNotification(ctx context.Context, action string, userID int32, targetID, id int32, isComment bool) error { // yangi funksiya
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

	var notification LikeNotificationMessage
	if isComment {
		notification = LikeNotificationMessage{
			UserID:     userID,       
			CommenterID: targetID,     
			Action:     action,        
			Message:    fmt.Sprintf("Commentga like bosildi: %d", id ), 
		}
	} else {
		notification = LikeNotificationMessage{
			UserID:     userID,       
			TweetOwnerID: targetID,   
			Action:     action,        			
			Message:    fmt.Sprintf("Tweetga like bosildi: %d", id), 
		}
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