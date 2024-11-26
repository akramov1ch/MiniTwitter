package models

import "time"

type Like struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	LikeIdentifier string `json:"like_identifier"` //tweet or comment
	LikedID int64 `json:"liked_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeIdentifier string

const (
	LikeIdentifierTweet LikeIdentifier = "tweet"
	LikeIdentifierComment LikeIdentifier = "comment"
)

