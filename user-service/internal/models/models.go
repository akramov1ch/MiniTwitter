package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	FollowerCount int32 `json:"follower_count"`
	FollowingCount int32 `json:"following_count"`
	TweetCount int32 `json:"tweet_count"`
	Bio string `json:"bio"`
	Tweets []int32 `json:"tweets"`
	Followers []int32 `json:"followers"`
	Following []int32 `json:"following"`
	IsPrivate bool `json:"is_private"`
	BlockedUsers []int32 `json:"blocked_users"`
}

type Avatar struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	URL string `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type Code struct {
	ID int64 `json:"id"`
	LoginIdentifier string `json:"login_identifier"`
	Code string `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

