package service

import (
	"context"
	"database/sql"
	"errors"
	"like-service/internal/kafka"
	"like-service/internal/methods"
	"like-service/internal/models"
	"like-service/pkg/proto"
	"like-service/utils"

	"github.com/go-redis/redis/v8"
)

type LikeService struct {
	db          *sql.DB
	RedisClient *redis.Client
}

func NewLikeService(db *sql.DB, redisClient *redis.Client) *LikeService {
	return &LikeService{db: db, RedisClient: redisClient}
}

func (s *LikeService) CreateLikeTweet(ctx context.Context, like *models.Like) (models.Like, error) {
	query := `
		INSERT INTO likes (user_id, like_identifier, liked_id)
		VALUES ($1, $2, $3)
		RETURNING *
	`
	tweet := methods.ConnectTweet()
	_, err := tweet.GetTweetByID(ctx, &proto.GetTweetByIDRequest{TweetId: int32(like.LikedID)})
	if err != nil {
		return models.Like{}, err
	}
	user := methods.ConnectUser()
	userResponse, err := user.GetUser(ctx, &proto.GetUserRequest{UserId: int32(like.UserID)})
	if err != nil {
		return models.Like{}, err
	}
	if userResponse.User.IsPrivate {
		if !utils.InSlice(userResponse.User.Followers, int32(like.UserID)) {
			return models.Like{}, errors.New("user is private and not in followers")
		}
	}
	if utils.InSlice(userResponse.User.BlockedUsers, int32(like.UserID)) {
		return models.Like{}, errors.New("user is blocked")
	}
	row := s.db.QueryRowContext(ctx, query, like.UserID, like.LikeIdentifier, like.LikedID)
	var newLike models.Like
	err = row.Scan(&newLike.ID, &newLike.UserID, &newLike.LikeIdentifier, &newLike.LikedID, &newLike.CreatedAt)
	if err != nil {
		return models.Like{}, err
	}
	_, err = tweet.AddLike(ctx, &proto.AddLikeRequest{TweetId: int32(like.LikedID), UserId: int32(like.UserID)})
	if err != nil {
		return models.Like{}, err
	}

	twit, err := tweet.GetTweetByID(ctx, &proto.GetTweetByIDRequest{TweetId: int32(like.LikedID)})
	if err != nil {
		return models.Like{}, err
	}
	kafka.PublishLikeNotification(ctx, "like", int32(like.UserID), twit.Tweet.UserId, int32(like.LikedID), false)
	return newLike, nil
}

func (s *LikeService) CreateLikeComment(ctx context.Context, like *models.Like) (models.Like, error) {
	query := `
		INSERT INTO likes (user_id, like_identifier, liked_id)
		VALUES ($1, $2, $3) 
		RETURNING *
	`
	comment := methods.ConnectComment()
	_, err := comment.GetComment(ctx, &proto.GetCommentRequest{Id: int64(like.LikedID)})
	if err != nil {
		return models.Like{}, err
	}
	user := methods.ConnectUser()
	userResponse, err := user.GetUser(ctx, &proto.GetUserRequest{UserId: int32(like.UserID)})
	if err != nil {
		return models.Like{}, err
	}
	if userResponse.User.IsPrivate {
		if !utils.InSlice(userResponse.User.Followers, int32(like.UserID)) {
			return models.Like{}, errors.New("user is private and not in followers")
		}
	}
	if utils.InSlice(userResponse.User.BlockedUsers, int32(like.UserID)) {
		return models.Like{}, errors.New("user is blocked")
	}
	row := s.db.QueryRowContext(ctx, query, like.UserID, like.LikeIdentifier, like.LikedID)
	var newLike models.Like
	err = row.Scan(&newLike.ID, &newLike.UserID, &newLike.LikeIdentifier, &newLike.LikedID, &newLike.CreatedAt)
	if err != nil {
		return models.Like{}, err
	}
	_, err = comment.LikeComment(ctx, &proto.LikeCommentRequest{CommentId: int64(like.LikedID), UserId: int64(like.UserID)})
	if err != nil {
		return models.Like{}, err
	}
	coment, err := comment.GetComment(ctx, &proto.GetCommentRequest{Id: int64(like.LikedID)})
	if err != nil {
		return models.Like{}, err
	}
	kafka.PublishLikeNotification(ctx, "like", int32(like.UserID), int32(coment.Comment.UserId), int32(like.LikedID), true)
	return newLike, nil
}

func (s *LikeService) DeleteLikeTweet(ctx context.Context, like *models.Like) error {
	query := `
		DELETE FROM likes
		WHERE user_id = $1 AND like_identifier = $2 AND liked_id = $3
	`
	_, err := s.db.ExecContext(ctx, query, like.UserID, like.LikeIdentifier, like.LikedID)
	if err != nil {
		return err
	}
	tweet := methods.ConnectTweet()
	_, err = tweet.RemoveLike(ctx, &proto.RemoveLikeRequest{TweetId: int32(like.LikedID), UserId: int32(like.UserID)})
	if err != nil {
		return err
	}
	return nil
}

func (s *LikeService) DeleteLikeComment(ctx context.Context, like *models.Like) error {
	query := `
		DELETE FROM likes
		WHERE user_id = $1 AND like_identifier = $2 AND liked_id = $3
	`
	_, err := s.db.ExecContext(ctx, query, like.UserID, like.LikeIdentifier, like.LikedID)
	if err != nil {
		return err
	}
	comment := methods.ConnectComment()
	_, err = comment.RemoveLikeFromComment(ctx, &proto.RemoveLikeFromCommentRequest{CommentId: int64(like.LikedID), UserId: int64(like.UserID)})
	if err != nil {
		return err
	}
	return nil
}

func (s *LikeService) GetLikesTweet(ctx context.Context, like *models.Like) ([]models.Like, error) {
	query := `
		SELECT * FROM likes
		WHERE like_identifier = $1 AND liked_id = $2
	`
	var likes []models.Like
	err := s.db.QueryRowContext(ctx, query, like.LikeIdentifier, like.LikedID).Scan(&likes)
	if err != nil {
		return []models.Like{}, err
	}					
	return likes, nil
}

func (s *LikeService) GetLikesComment(ctx context.Context, like *models.Like) ([]models.Like, error) {
	query := `
		SELECT * FROM likes
		WHERE like_identifier = $1 AND liked_id = $2
	`
	var likes []models.Like
	err := s.db.QueryRowContext(ctx, query, like.LikeIdentifier, like.LikedID).Scan(&likes)
	if err != nil {
		return []models.Like{}, err
	}
	return likes, nil
}

func (s *LikeService) GetLikedTweets(ctx context.Context, like *models.Like) ([]models.Like, error) {
	query := `
		SELECT * FROM likes
		WHERE user_id = $1 AND liked_id = $2 AND like_identifier = $3
		RETURNING *
	`
	var likes []models.Like
	err := s.db.QueryRowContext(ctx, query, like.UserID, like.LikedID, like.LikeIdentifier).Scan(&likes)
	if err != nil {
		return []models.Like{}, err
	}
	return likes, nil
}

func (s *LikeService) GetLikeTweetByUser(ctx context.Context, like *models.Like) ([]models.Like, error) {
	query := `
		SELECT * FROM likes
		WHERE user_id = $1 AND like_identifier = $2 AND liked_id = $3
	`
	var likes []models.Like
	err := s.db.QueryRowContext(ctx, query, like.UserID, like.LikeIdentifier, like.LikedID).Scan(&likes)
	if err != nil {
		return []models.Like{}, err
	}
	return likes, nil
}

func (s *LikeService) GetLikeCommentByUser(ctx context.Context, like *models.Like) ([]models.Like, error) {
	query := `
		SELECT * FROM likes
		WHERE user_id = $1 AND like_identifier = $2 AND liked_id = $3
	`
	var likes []models.Like
	err := s.db.QueryRowContext(ctx, query, like.UserID, like.LikeIdentifier, like.LikedID).Scan(&likes)
	if err != nil {
		return []models.Like{}, err
	}
	return likes, nil
}
