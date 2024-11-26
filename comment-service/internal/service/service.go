package service

import (
	"comment-service/internal/kafka"
	"comment-service/internal/methods"
	"comment-service/internal/models"
	"comment-service/pkg/proto"
	"comment-service/utils"
	"context"
	"database/sql"
	"errors"

	"github.com/go-redis/redis/v8"
)

type CommentService struct {
	db          *sql.DB
	RedisClient *redis.Client
}

func NewCommentService(db *sql.DB, redisClient *redis.Client) *CommentService {
	return &CommentService{db: db, RedisClient: redisClient}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	query := `
		INSERT INTO comments (user_id, tweet_id, content)
		VALUES ($1, $2, $3)
		RETURNING *
	`
	tweets := methods.ConnectTweet()
	_, err := tweets.GetTweetByID(ctx, &proto.GetTweetByIDRequest{TweetId: int32(comment.TweetID)})
	if err != nil {
		return nil, err
	}
	users := methods.ConnectUser()
	user, err := users.GetUser(ctx, &proto.GetUserRequest{UserId: int32(comment.UserID)})
	if err != nil {
		return nil, err				
	}

	if user.User.IsPrivate {
		if !utils.InSlice(user.User.Followers, int32(comment.UserID)) {
			return nil, errors.New("this user is private and you are not following")
		}
	}
	if !utils.InSlice(user.User.BlockedUsers, int32(comment.UserID)) {
		return nil, errors.New("this user has blocked you")
	}
	var createdComment models.Comment
	err = s.db.QueryRowContext(ctx, query, comment.UserID, comment.TweetID, comment.Content).Scan(&createdComment.ID, &createdComment.UserID, &createdComment.TweetID, &createdComment.Content, &createdComment.CreatedAt, &createdComment.Likes)
	if err != nil {
		return nil, err
	}
	_, err = tweets.AddComment(ctx, &proto.AddCommentRequest{
		UserId: int32(comment.UserID),
		TweetId: int32(comment.TweetID),
		CommentId: int32(createdComment.ID),
	})
	if err != nil {
		return nil, err
	}
	err = kafka.PublishCommentNotification(ctx, int32(comment.UserID), int32(comment.TweetID), int32(comment.UserID), comment.Content)
	if err != nil {
		return nil, err
	}
	return &createdComment, nil
}

func (s *CommentService) GetCommentsByTweetID(ctx context.Context, tweetID int64) ([]*models.Comment, error) {
	query := `
		SELECT * FROM comments WHERE tweet_id = $1
	`
	tweets := methods.ConnectTweet()
	tweet, err := tweets.GetTweetByID(ctx, &proto.GetTweetByIDRequest{TweetId: int32(tweetID)})
	if err != nil {
		return nil, err
	}
	userID := tweet.Tweet.UserId
	users := methods.ConnectUser()
	user, err := users.GetUser(ctx, &proto.GetUserRequest{UserId: int32(userID)})
	if err != nil {
		return nil, err
	}
	if user.User.IsPrivate {
		if !utils.InSlice(user.User.Followers, int32(userID)) {
			return nil, errors.New("this user is private and you are not following")
		}
	}
	if !utils.InSlice(user.User.BlockedUsers, int32(userID)) {
		return nil, errors.New("this user has blocked you")
	}
	rows, err := s.db.QueryContext(ctx, query, tweetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.TweetID, &comment.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id int64, userID int64, tweetID int64) error {
	query := `
		SELECT * FROM comments WHERE id = $1 RETURNING *
	`
	var comment models.Comment
	err := s.db.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.UserID, &comment.TweetID, &comment.Content, &comment.CreatedAt, &comment.Likes)
	if err != nil {
		return err
	}

	tweets := methods.ConnectTweet()
	tweet, err := tweets.GetTweetByID(ctx, &proto.GetTweetByIDRequest{TweetId: int32(comment.TweetID)})
	if err != nil {
		return err
	}
	if userID != comment.UserID {
		return errors.New("this is not your comment")
	} else if userID != int64(tweet.Tweet.UserId) {
		return errors.New("this is not your tweet")
	}

	deleteQuery := `
		DELETE FROM comments WHERE id = $1
	`
	_, err = s.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}
	_, err = tweets.RemoveComment(ctx, &proto.RemoveCommentRequest{
		UserId: int32(comment.UserID),
		TweetId: int32(comment.TweetID),
		CommentId: int32(comment.ID),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentService) LikeComment(ctx context.Context, userID int64, commentID int64) (bool, error) {
	query := `
		SELECT * FROM comments WHERE id = $1 RETURNING *
	`
	var comment models.Comment
	err := s.db.QueryRowContext(ctx, query, commentID).Scan(&comment.ID, &comment.UserID, &comment.TweetID, &comment.Content, &comment.CreatedAt, &comment.Likes)
	if err != nil {
		return false, err
	}

	comment.Likes = append(comment.Likes, userID)
	comment.LikesCount++

	updateQuery := `
		UPDATE comments SET likes = $1, likes_count = $2 WHERE id = $3
	`
	_, err = s.db.ExecContext(ctx, updateQuery, comment.Likes, comment.LikesCount, commentID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *CommentService) RemoveLikeFromComment(ctx context.Context, userID int64, commentID int64) (bool, error) {
	query := `
		SELECT * FROM comments WHERE id = $1 RETURNING *
	`
	var comment models.Comment
	err := s.db.QueryRowContext(ctx, query, commentID).Scan(&comment.ID, &comment.UserID, &comment.TweetID, &comment.Content, &comment.CreatedAt, &comment.Likes)
	if err != nil {
		return false, err
	}

	comment.Likes = utils.RemoveElement(comment.Likes, userID)
	comment.LikesCount--

	updateQuery := `
		UPDATE comments SET likes = $1, likes_count = $2 WHERE id = $3
	`
	_, err = s.db.ExecContext(ctx, updateQuery, comment.Likes, comment.LikesCount, commentID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *CommentService) GetCommentByID(ctx context.Context, id int64, userID int64) (*models.Comment, error) {
	query := `
		SELECT * FROM comments WHERE id = $1 RETURNING *
	`
	var comment models.Comment
	err := s.db.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.UserID, &comment.TweetID, &comment.Content, &comment.CreatedAt, &comment.Likes)
	if err != nil {
		return nil, err
	}
	comment.LikesCount = int32(len(comment.Likes))
	users := methods.ConnectUser()
	user, err := users.GetUser(ctx, &proto.GetUserRequest{UserId: int32(userID)})
	if err != nil {
		return nil, err
	}
	if user.User.IsPrivate {
		if !utils.InSlice(user.User.Followers, int32(userID)) {
			return nil, errors.New("this user is private and you are not following")
		}
	}
	if !utils.InSlice(user.User.BlockedUsers, int32(userID)) {
		return nil, errors.New("this user has blocked you")
	}
	return &comment, nil
}
