package models

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TweetID   int64     `json:"tweet_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes     []int64   `json:"likes"`
	LikesCount int32 `json:"likes_count"`
}

