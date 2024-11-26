package models

import "time"

type DirectMessage struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	TweetID    int64     `json:"tweet_id"`
	Text       string    `json:"text"`
	Media      []string  `json:"media"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowNotification struct {
	UserID     int32  `json:"user_id"`
	FollowerID int32  `json:"follower_id"`
	Action     string `json:"action"` // "add", "accept", "reject", "request"
	Message    string `json:"message"`
}

type LikeNotification struct {
	UserID      int32  `json:"user_id"`
	CommenterID int32  `json:"commenter_id"`
	TweetOwnerID int32  `json:"tweet_owner_id"`
	Action      string `json:"action"` // "like", "unlike"
	Message     string `json:"message"`
	TweetID     int32  `json:"tweet_id"`
}

type CommentNotification struct {
	UserID      int32  `json:"user_id"`
	TweetID     int32  `json:"tweet_id"`
	CommenterID int32  `json:"commenter_id"`
	Action      string `json:"action"` // "add", "remove"
	Message     string `json:"message"`
}