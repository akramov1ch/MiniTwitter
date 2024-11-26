package service

import (
	"context"
	"database/sql"
	"errors"
	"tweet-service/internal/methods"
	"tweet-service/internal/models"
	"tweet-service/pkg/proto"
	"tweet-service/utils"

	"github.com/go-redis/redis/v8"
)

type TweetService struct {
	db          *sql.DB
	RedisClient *redis.Client
}

func NewTweetService(db *sql.DB, redisClient *redis.Client) *TweetService {
	return &TweetService{db: db, RedisClient: redisClient}
}

func (s *TweetService) CreateTweet(ctx context.Context, tweet *models.Tweet) error {
	query := `
		INSERT INTO tweets (user_id, content, created_at)
		VALUES ($1, $2, $3)
		RETURNING *
	`
	err := s.db.QueryRowContext(ctx, query, tweet.UserID, tweet.Content, tweet.CreatedAt).Scan(
		&tweet.ID,
		&tweet.Content,
		&tweet.UserID,
		&tweet.Username,
		&tweet.CreatedAt,
		&tweet.Media,
		&tweet.Likes,
		&tweet.Comments,
		&tweet.Shares,
		&tweet.Saves,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) GetTweetsByUser(ctx context.Context, userID, userIDToGet int32) ([]*models.Tweet, error) {
	query := `
		SELECT * FROM tweets WHERE user_id = $1
	`
	users := methods.ConnectUser()
	user, err := users.GetUser(ctx, &proto.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if user.User.IsPrivate {
		if !utils.Contains(user.User.Followers, userIDToGet) {
			return nil, errors.New("user is private and you are blocked")
		}
	}

	if !utils.Contains(user.User.BlockedUsers, userIDToGet) {
		return nil, errors.New("user is block you")
	}
	rows, err := s.db.QueryContext(ctx, query, userIDToGet)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tweets := []*models.Tweet{}
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(
			&tweet.ID,
			&tweet.Content,
			&tweet.UserID,
			&tweet.Username,
			&tweet.CreatedAt,
			&tweet.Media,
			&tweet.Likes,
			&tweet.Comments,
			&tweet.Shares,
			&tweet.Saves,
		)
		if err != nil {
			return nil, err
		}
		tweet.LikeCount = int64(len(tweet.Likes))
		tweet.CommentCount = int64(len(tweet.Comments))
		tweet.ShareCount = int64(len(tweet.Shares))
		tweet.SaveCount = int64(len(tweet.Saves))
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}

func (s *TweetService) UpdateTweet(ctx context.Context, tweet *models.Tweet) error {
	query := `
		UPDATE tweets SET content = $1 WHERE id = $2
	`
	_, err := s.db.ExecContext(ctx, query, tweet.Content, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) DeleteTweet(ctx context.Context, tweetID int32) error {
	query := `
		DELETE FROM tweets WHERE id = $1
	`
	_, err := s.db.ExecContext(ctx, query, tweetID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) AddLike(ctx context.Context, tweetID int32, userID int32) error {
	query := `SELECT * FROM tweets WHERE id = $1 RETURNING *`
	var tweet models.Tweet
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.Username, &tweet.CreatedAt, &tweet.Media, &tweet.Likes, &tweet.Comments, &tweet.Shares, &tweet.Saves)
	if err != nil {
		return err
	}

	tweet.Likes = append(tweet.Likes, userID)
	query = `UPDATE tweets SET likes = $1 WHERE id = $2`
	_, err = s.db.ExecContext(ctx, query, tweet.Likes, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) RemoveLike(ctx context.Context, tweetID int32, userID int32) error {
	query := `SELECT * FROM tweets WHERE id = $1 RETURNING *`
	var tweet models.Tweet
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.Username, &tweet.CreatedAt, &tweet.Media, &tweet.Likes, &tweet.Comments, &tweet.Shares, &tweet.Saves)
	if err != nil {
		return err
	}

	tweet.Likes = utils.RemoveElement(tweet.Likes, userID)
	query = `UPDATE tweets SET likes = $1 WHERE id = $2`
	_, err = s.db.ExecContext(ctx, query, tweet.Likes, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) AddComment(ctx context.Context, tweetID int32, userID int32, commentID int32) error {
	query := `SELECT * FROM tweets WHERE id = $1 RETURNING *`
	var tweet models.Tweet
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.Username, &tweet.CreatedAt, &tweet.Media, &tweet.Likes, &tweet.Comments, &tweet.Shares, &tweet.Saves)
	if err != nil {
		return err
	}

	tweet.Comments = append(tweet.Comments, commentID)
	query = `UPDATE tweets SET comments = $1 WHERE id = $2`
	_, err = s.db.ExecContext(ctx, query, tweet.Comments, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) RemoveComment(ctx context.Context, tweetID int32, userID int32, commentID int32) error {
	query := `SELECT * FROM tweets WHERE id = $1 RETURNING *`
	var tweet models.Tweet
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.Username, &tweet.CreatedAt, &tweet.Media, &tweet.Likes, &tweet.Comments, &tweet.Shares, &tweet.Saves)
	if err != nil {
		return err
	}

	tweet.Comments = utils.RemoveElement(tweet.Comments, commentID)
	query = `UPDATE tweets SET comments = $1 WHERE id = $2`
	_, err = s.db.ExecContext(ctx, query, tweet.Comments, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) AddShare(ctx context.Context, tweetID int32, userID int32) error {
	query := `SELECT * FROM tweets WHERE id = $1 RETURNING *`
	var tweet models.Tweet
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.Username, &tweet.CreatedAt, &tweet.Media, &tweet.Likes, &tweet.Comments, &tweet.Shares, &tweet.Saves)
	if err != nil {
		return err
	}

	tweet.Shares = append(tweet.Shares, userID)
	query = `UPDATE tweets SET shares = $1 WHERE id = $2`
	_, err = s.db.ExecContext(ctx, query, tweet.Shares, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) SaveTweet(ctx context.Context, tweetID int32, userID int32) error {
	query := `SELECT * FROM tweets WHERE id = $1 RETURNING *`
	var tweet models.Tweet
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.Username, &tweet.CreatedAt, &tweet.Media, &tweet.Likes, &tweet.Comments, &tweet.Shares, &tweet.Saves)
	if err != nil {
		return err
	}

	tweet.Saves = append(tweet.Saves, userID)
	query = `UPDATE tweets SET saves = $1 WHERE id = $2`
	_, err = s.db.ExecContext(ctx, query, tweet.Saves, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) RemoveSave(ctx context.Context, tweetID int32, userID int32) error {
	query := `SELECT * FROM tweets WHERE id = $1 RETURNING *`
	var tweet models.Tweet
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.Username, &tweet.CreatedAt, &tweet.Media, &tweet.Likes, &tweet.Comments, &tweet.Shares, &tweet.Saves)
	if err != nil {
		return err
	}

	tweet.Saves = utils.RemoveElement(tweet.Saves, userID)
	query = `UPDATE tweets SET saves = $1 WHERE id = $2`
	_, err = s.db.ExecContext(ctx, query, tweet.Saves, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TweetService) GetSavedTweets(ctx context.Context, userID int32) ([]*models.Tweet, error) {
	query := `SELECT * FROM tweets WHERE id IN (SELECT unnest(saves) FROM tweets WHERE user_id = $1)`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tweets := []*models.Tweet{}
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(
			&tweet.ID,
			&tweet.Content,
			&tweet.UserID,
			&tweet.Username,
			&tweet.CreatedAt,
			&tweet.Media,
			&tweet.Likes,
			&tweet.Comments,
			&tweet.Shares,
			&tweet.Saves,
		)
		if err != nil {
			return nil, err
		}
		tweet.LikeCount = int64(len(tweet.Likes))
		tweet.CommentCount = int64(len(tweet.Comments))
		tweet.ShareCount = int64(len(tweet.Shares))
		tweet.SaveCount = int64(len(tweet.Saves))
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}

func (s *TweetService) GetTweetByID(ctx context.Context, tweetID int32, userIDToGet int32) (*models.Tweet, error) {
	query := `SELECT * FROM tweets WHERE id = $1`
	users := methods.ConnectUser()
	var tweet models.Tweet
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.Username, &tweet.CreatedAt, &tweet.Media, &tweet.Likes, &tweet.Comments, &tweet.Shares, &tweet.Saves)
	if err != nil {
		return nil, err
	}
	user, err := users.GetUser(ctx, &proto.GetUserRequest{UserId: int32(tweet.UserID)})
	if err != nil {
		return nil, err
	}
	if user.User.IsPrivate {
		if !utils.Contains(user.User.Followers, userIDToGet) {
			return nil, errors.New("user is private and you are blocked")
		}
	}
	if !utils.Contains(user.User.BlockedUsers, userIDToGet) {
		return nil, errors.New("user is block you")
	}
	return &tweet, nil
}

func (s *TweetService) GetLikedTweets(ctx context.Context, userID int32) ([]*models.Tweet, error) {
	query := `SELECT * FROM tweets WHERE id IN (SELECT unnest(likes) FROM tweets WHERE user_id = $1)`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tweets := []*models.Tweet{}
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(
			&tweet.ID,
			&tweet.Content,
			&tweet.UserID,
			&tweet.Username,
			&tweet.CreatedAt,
			&tweet.Media,
			&tweet.Likes,
			&tweet.Comments,
			&tweet.Shares,
			&tweet.Saves,
		)
		if err != nil {
			return nil, err
		}
		tweet.LikeCount = int64(len(tweet.Likes))
		tweet.CommentCount = int64(len(tweet.Comments))
		tweet.ShareCount = int64(len(tweet.Shares))
		tweet.SaveCount = int64(len(tweet.Saves))
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}
