package handlers

import (
	"context"
	"like-service/internal/models"
	"like-service/internal/service"
	pb "like-service/pkg/proto"
)

type LikeHandler struct {
	pb.UnimplementedLikeServiceServer
	service service.LikeService
}

func NewLikeHandler(service service.LikeService) *LikeHandler {
	return &LikeHandler{service: service}
}

func (h *LikeHandler) CreateLikeTweet(ctx context.Context, req *pb.CreateLikeTweetRequest) (*pb.CreateLikeTweetResponse, error) {
	like, err := h.service.CreateLikeTweet(ctx, &models.Like{
		UserID:         int64(req.UserId),
		LikeIdentifier: string(models.LikeIdentifierTweet),
		LikedID:        int64(req.TweetId),
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateLikeTweetResponse{
		Success: true,
		Message: "Like created successfully",
		Like: &pb.Like{
			Id:             int32(like.ID),
			UserId:         int32(like.UserID),
			LikeIdentifier: string(like.LikeIdentifier),
			LikedId:        int32(like.LikedID),
			CreatedAt:      like.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (h *LikeHandler) GetLikedTweet(ctx context.Context, req *pb.GetLikedTweetRequest) (*pb.GetLikedTweetResponse, error) {
	likes, err := h.service.GetLikedTweets(ctx, &models.Like{
		UserID:         int64(req.UserId),
		LikeIdentifier: string(models.LikeIdentifierTweet),
		LikedID:        int64(req.TweetId),
	})
	if err != nil {
		return nil, err
	}
	likesProto := make([]*pb.Like, len(likes))
	for i, like := range likes {
		likesProto[i] = &pb.Like{
			Id:             int32(like.ID),
			UserId:         int32(like.UserID),
			LikeIdentifier: string(like.LikeIdentifier),
			LikedId:        int32(like.LikedID),
			CreatedAt:      like.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &pb.GetLikedTweetResponse{
		Success: true,
		Message: "Like retrieved successfully",
		Likes:   likesProto,
	}, nil
}

func (h *LikeHandler) DeleteLikeTweet(ctx context.Context, req *pb.DeleteLikeTweetRequest) (*pb.DeleteLikeTweetResponse, error) {
	err := h.service.DeleteLikeTweet(ctx, &models.Like{
		UserID:         int64(req.UserId),
		LikeIdentifier: string(models.LikeIdentifierTweet),
		LikedID:        int64(req.TweetId),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteLikeTweetResponse{
		Success: true,
		Message: "Like deleted successfully",
	}, nil
}

func (h *LikeHandler) GetLikesTweet(ctx context.Context, req *pb.GetLikesTweetRequest) (*pb.GetLikesTweetResponse, error) {
	likes, err := h.service.GetLikesTweet(ctx, &models.Like{
		LikeIdentifier: string(models.LikeIdentifierTweet),
		LikedID:        int64(req.TweetId),
	})
	if err != nil {
		return nil, err
	}
	likesProto := make([]*pb.Like, len(likes))
	for i, like := range likes {
		likesProto[i] = &pb.Like{
			Id:             int32(like.ID),
			UserId:         int32(like.UserID),
			LikeIdentifier: string(like.LikeIdentifier),
			LikedId:        int32(like.LikedID),
			CreatedAt:      like.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &pb.GetLikesTweetResponse{
		Success: true,
		Message: "Likes retrieved successfully",
		Likes:   likesProto,
	}, nil
}

func (h *LikeHandler) GetLikeTweetByUser(ctx context.Context, req *pb.GetLikeTweetByUserRequest) (*pb.GetLikeTweetByUserResponse, error) {
	likes, err := h.service.GetLikeTweetByUser(ctx, &models.Like{
		UserID: int64(req.UserId),
		LikeIdentifier: string(models.LikeIdentifierTweet),
	})
	if err != nil {
		return nil, err
	}
	likesProto := make([]*pb.Like, len(likes))
	for i, like := range likes {
		likesProto[i] = &pb.Like{
			Id: int32(like.ID),
			UserId: int32(like.UserID),
			LikeIdentifier: string(like.LikeIdentifier),
			LikedId: int32(like.LikedID),
			CreatedAt: like.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &pb.GetLikeTweetByUserResponse{
		Success: true,
		Message: "Likes retrieved successfully",
		Likes:   likesProto,
	}, nil
}	

func (h *LikeHandler) CreateLikeComment(ctx context.Context, req *pb.CreateLikeCommentRequest) (*pb.CreateLikeCommentResponse, error) {
	like, err := h.service.CreateLikeComment(ctx, &models.Like{
		UserID:         int64(req.UserId),
		LikeIdentifier: string(models.LikeIdentifierComment),
		LikedID:        int64(req.CommentId),
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateLikeCommentResponse{
		Success: true,
		Message: "Like created successfully",
		Like: &pb.Like{
			Id:             int32(like.ID),
			UserId:         int32(like.UserID),
			LikeIdentifier: string(like.LikeIdentifier),
			LikedId:        int32(like.LikedID),
			CreatedAt:      like.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (h *LikeHandler) GetLikesComment(ctx context.Context, req *pb.GetLikesCommentRequest) (*pb.GetLikesCommentResponse, error) {
	likes, err := h.service.GetLikesComment(ctx, &models.Like{
		LikeIdentifier: string(models.LikeIdentifierComment),
		LikedID:        int64(req.CommentId),
	})
	if err != nil {
		return nil, err
	}
	likesProto := make([]*pb.Like, len(likes))
	for i, like := range likes {
		likesProto[i] = &pb.Like{	
			Id:             int32(like.ID),
			UserId:         int32(like.UserID),
			LikeIdentifier: string(like.LikeIdentifier),
			LikedId:        int32(like.LikedID),
			CreatedAt:      like.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &pb.GetLikesCommentResponse{
		Success: true,
		Message: "Likes retrieved successfully",
		Likes:   likesProto,
	}, nil
}

func (h *LikeHandler) DeleteLikeComment(ctx context.Context, req *pb.DeleteLikeCommentRequest) (*pb.DeleteLikeCommentResponse, error) {
	err := h.service.DeleteLikeComment(ctx, &models.Like{
		LikeIdentifier: string(models.LikeIdentifierComment),
		LikedID:        int64(req.CommentId),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteLikeCommentResponse{
		Success: true,
		Message: "Like deleted successfully",
	}, nil
}  

func (h *LikeHandler) GetLikeCommentByUser(ctx context.Context, req *pb.GetLikeCommentByUserRequest) (*pb.GetLikeCommentByUserResponse, error) {
	likes, err := h.service.GetLikeCommentByUser(ctx, &models.Like{
		UserID: int64(req.UserId),
	})
	if err != nil {
		return nil, err
	}
	likesProto := make([]*pb.Like, len(likes))
	for i, like := range likes {
		likesProto[i] = &pb.Like{
			Id: int32(like.ID),
			UserId: int32(like.UserID),
			LikeIdentifier: string(like.LikeIdentifier),
			LikedId: int32(like.LikedID),
			CreatedAt: like.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &pb.GetLikeCommentByUserResponse{
		Success: true,
		Message: "Likes retrieved successfully",
		Likes:   likesProto,
	}, nil
}