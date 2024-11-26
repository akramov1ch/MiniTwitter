package service

import (
	"context"
	"database/sql"
	"direct-service/internal/kafka"
	"direct-service/internal/methods"
	"direct-service/internal/models"
	"direct-service/pkg/proto"
	"direct-service/utils"
	"errors"
	"github.com/go-redis/redis/v8"
)

type DirectService struct {
	db          *sql.DB
	RedisClient *redis.Client
}

func NewDirectService(db *sql.DB, redisClient *redis.Client) *DirectService {
	return &DirectService{db: db, RedisClient: redisClient}
}

func (s *DirectService) CreateDirectMessage(ctx context.Context, message *models.DirectMessage) (*models.DirectMessage, error) {
	query := `INSERT INTO directs (sender_id, receiver_id, tweet_id, text, media) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	user := methods.ConnectUser()
	sender, err := user.GetUser(ctx, &proto.GetUserRequest{UserId: int32(message.SenderID)})
	if err != nil {
		return nil, err
	}
	receiver, err := user.GetUser(ctx, &proto.GetUserRequest{UserId: int32(message.ReceiverID)})
	if err != nil {
		return nil, err
	}

	if sender.User.IsPrivate && sender.User.Id != receiver.User.Id {
		if !utils.InSlice(sender.User.Followers, int32(receiver.User.Id)) {
			return nil, errors.New("user is private and you are not following")
		}
	}
	if receiver.User.IsPrivate && sender.User.Id != receiver.User.Id {
		if !utils.InSlice(receiver.User.Followers, int32(sender.User.Id)) {
			return nil, errors.New("user is private and you are not following")
		}
	}
	if !utils.InSlice(sender.User.BlockedUsers, int32(receiver.User.Id)) || !utils.InSlice(receiver.User.BlockedUsers, int32(sender.User.Id)) {
		return nil, errors.New("user is blocked")
	}
	tweet := methods.ConnectTweet()
	if message.TweetID != 0 {
		_, err = tweet.GetTweetByID(ctx, &proto.GetTweetByIDRequest{TweetId: int32(message.TweetID)})
		if err != nil {
			return nil, err
		}
		_, err = tweet.AddShare(ctx, &proto.AddShareRequest{TweetId: int32(message.TweetID), UserId: int32(message.SenderID)})
		if err != nil {
			return nil, err
		}
	}
	err = s.db.QueryRowContext(ctx, query, message.SenderID, message.ReceiverID, message.TweetID, message.Text, message.Media).Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.TweetID, &message.Text, &message.Media, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	err = kafka.PublishDirectMessage(ctx, "direct-message", *message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *DirectService) GetDirectMessages(ctx context.Context, senderID, receiverID int64) ([]*models.DirectMessage, error) {
	query := `SELECT * FROM directs WHERE sender_id = $1 AND receiver_id = $2 OR sender_id = $2 AND receiver_id = $1`
	user := methods.ConnectUser()
	_, err := user.GetUser(ctx, &proto.GetUserRequest{UserId: int32(senderID)})
	if err != nil {
		return nil, err
	}
	_, err = user.GetUser(ctx, &proto.GetUserRequest{UserId: int32(receiverID)})
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, query, senderID, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.DirectMessage
	for rows.Next() {
		var message models.DirectMessage
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.TweetID, &message.Text, &message.Media, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *DirectService) GetDirectMessageByID(ctx context.Context, id int64) (*models.DirectMessage, error) {
	query := `SELECT * FROM directs WHERE id = $1`
	var message models.DirectMessage
	err := s.db.QueryRowContext(ctx, query, id).Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.TweetID, &message.Text, &message.Media, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (s *DirectService) DeleteDirectMessage(ctx context.Context, id int64) (bool, error) {
	query := `DELETE FROM directs WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
