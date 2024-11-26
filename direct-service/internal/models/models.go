package models

import "time"

type DirectMessage struct {
	ID int64 `json:"id"`
	SenderID int64 `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`
	TweetID int64 `json:"tweet_id"`
	Text string `json:"text"`
	Media []string `json:"media"`
	CreatedAt time.Time `json:"created_at"`
}

