package models

import "time"

type Tweet struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Media     []string  `json:"media"`	
	Likes     []int32   `json:"likes"`
	Comments  []int32   `json:"comments"`
	Shares    []int32   `json:"shares"`
	Saves     []int32   `json:"saves"`
	MediaCount int64     `json:"media_count"`
	LikeCount int64     `json:"like_count"`
	CommentCount int64     `json:"comment_count"`
	ShareCount int64     `json:"share_count"`
	SaveCount  int64     `json:"save_count"`
	IsSaved    bool      `json:"is_saved"`
	IsLiked    bool      `json:"is_liked"`
}
