package grpc

import (
	"context"
	"log"
	"notification-service/internal/models"
	"notification-service/internal/producer"
	"notification-service/proto"
)

type NotificationServiceServer struct {
	proto.UnimplementedNotificationServiceServer
}

// SendDirectMessage qabul qiluvchi foydalanuvchiga xabar yuboradi
func (s *NotificationServiceServer) SendDirectMessage(ctx context.Context, msg *proto.DirectMessage) (*proto.EmptyResponse, error) {
	err := HandleDirectMessage(ctx, models.DirectMessage{
		SenderID:   int64(msg.SenderId),
		ReceiverID: int64(msg.ReceiverId),
		Media:      msg.Media,
		Text:       msg.Text,
	})
	if err != nil {
		log.Printf("Error publishing direct message: %v\n", err)
		return &proto.EmptyResponse{}, err
	}
	log.Printf("Received DirectMessage: %+v\n", msg)
	return &proto.EmptyResponse{}, nil
}

// SendFollowNotification qabul qiluvchi foydalanuvchiga kuzatuv bildirishnomasini yuboradi
func (s *NotificationServiceServer) SendFollowNotification(ctx context.Context, msg *proto.FollowNotification) (*proto.EmptyResponse, error) {
	err := HandleFollowNotification(ctx, models.FollowNotification{
		UserID:     int32(msg.UserId),
		FollowerID: int32(msg.FollowerId),
		Message:    msg.Message,
	})
	if err != nil {
		log.Printf("Error publishing follow notification: %v\n", err)
		return &proto.EmptyResponse{}, err
	}
	log.Printf("Received FollowNotification: %+v\n", msg)
	return &proto.EmptyResponse{}, nil
}

// SendLikeNotification qabul qiluvchi foydalanuvchiga like bildirishnomasini yuboradi
func (s *NotificationServiceServer) SendLikeNotification(ctx context.Context, msg *proto.LikeNotification) (*proto.EmptyResponse, error) {
	err := HandleLikeNotification(ctx, models.LikeNotification{
		UserID:      int32(msg.UserId),
		CommenterID: int32(msg.CommenterId),
		TweetOwnerID: int32(msg.TweetOwnerId),
		Action:      msg.Action,
		Message:     msg.Message,
	})
	if err != nil {
		log.Printf("Error publishing like notification: %v\n", err)
		return &proto.EmptyResponse{}, err
	}
	log.Printf("Received LikeNotification: %+v\n", msg)
	return &proto.EmptyResponse{}, nil
}

// SendCommentNotification qabul qiluvchi foydalanuvchiga comment bildirishnomasini yuboradi
func (s *NotificationServiceServer) SendCommentNotification(ctx context.Context, msg *proto.CommentNotification) (*proto.EmptyResponse, error) {
	err := HandleCommentNotification(ctx, models.CommentNotification{
		UserID:      int32(msg.UserId),
		TweetID:     int32(msg.TweetId),
		CommenterID: int32(msg.CommenterId),
		Action:      msg.Action,
		Message:     msg.Message,
	})
	if err != nil {
		log.Printf("Error publishing comment notification: %v\n", err)
		return &proto.EmptyResponse{}, err
	}
	log.Printf("Received CommentNotification: %+v\n", msg)
	return &proto.EmptyResponse{}, nil
}

func HandleDirectMessage(ctx context.Context, directMessage models.DirectMessage) error {
	notification := proto.DirectMessage{
		SenderId:   int64(directMessage.SenderID), // Yuboruvchi foydalanuvchi
		ReceiverId: int64(directMessage.ReceiverID), // Qabul qiluvchi foydalanuvchi
		Media:      directMessage.Media,
		Text:       directMessage.Text,
	}

	return producer.PublishDirectMessage(ctx, &notification)
}

// HandleFollowNotification kuzatuvchi bildirishnomani yuboradi
func HandleFollowNotification(ctx context.Context, followNotification models.FollowNotification) error {
	notification := proto.FollowNotification{
		UserId:     int32(followNotification.UserID), // Qabul qiluvchi foydalanuvchi
		FollowerId: int32(followNotification.FollowerID), // Yuboruvchi foydalanuvchi
		Action:     "follow",                             // Harakat turi
		Message:    followNotification.Message,          // Xabar
	}

	return producer.PublishFollowNotification(ctx, &notification)
}

// HandleLikeNotification like bildirishnomani yuboradi
func HandleLikeNotification(ctx context.Context, likeNotification models.LikeNotification) error {
	notification := proto.LikeNotification{
		UserId:      int32(likeNotification.UserID), // Qabul qiluvchi foydalanuvchi
		CommenterId: int32(likeNotification.CommenterID), // Kommentator
		TweetOwnerId: int32(likeNotification.TweetOwnerID), // Tweet egasi
		Action:      "like",                             // Harakat turi
		Message:     likeNotification.Message,          // Xabar
	}

	return producer.PublishLikeNotification(ctx, &notification)
}

// HandleCommentNotification comment bildirishnomani yuboradi
func HandleCommentNotification(ctx context.Context, commentNotification models.CommentNotification) error {
	notification := proto.CommentNotification{
		UserId:      int32(commentNotification.UserID), // Qabul qiluvchi foydalanuvchi
		TweetId:     int32(commentNotification.TweetID), // Tweet ID
		CommenterId: int32(commentNotification.CommenterID), // Kommentator
		Action:      "comment",                             // Harakat turi
		Message:     commentNotification.Message,          // Xabar
	}

	return producer.PublishCommentNotification(ctx, &notification)
}