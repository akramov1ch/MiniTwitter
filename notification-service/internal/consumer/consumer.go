package consumer

import (
	"context"
	"encoding/json"
	"log"
	"notification-service/internal/grpc"
	"notification-service/internal/models" 

	"github.com/segmentio/kafka-go"
)

func ConsumeNotifications() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "notifications",
		GroupID: "notification-group",
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		var notification models.DirectMessage
		if err := json.Unmarshal(m.Value, &notification); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		err = grpc.HandleDirectMessage(context.Background(), notification)
		if err != nil {
			log.Println("Error handling direct message:", err)
		}
	}
}
