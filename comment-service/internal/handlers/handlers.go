package handlers

import (
	"comment-service/internal/models"
	"comment-service/internal/service"
	pb "comment-service/pkg/proto"
	"context"
)

type CommentHandler struct {
	pb.UnimplementedCommentServiceServer
	service service.CommentService
}

func NewCommentHandler(service service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	comment, err := h.service.CreateComment(ctx, &models.Comment{
		UserID:  int64(req.UserId),
		TweetID: int64(req.TweetId),
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateCommentResponse{
		Comment: &pb.Comment{
			Id:         int64(comment.ID),
			UserId:     int64(comment.UserID),
			TweetId:    int64(comment.TweetID),
			Content:    comment.Content,
			CreatedAt:  comment.CreatedAt.Format("2006-01-02 15:04:05"),
			Likes:      comment.Likes,
			LikesCount: int32(comment.LikesCount),
		},
	}, nil
}

func (h *CommentHandler) GetCommentsByTweetID(ctx context.Context, req *pb.GetCommentsByTweetIDRequest) (*pb.GetCommentsByTweetIDResponse, error) {
	comments, err := h.service.GetCommentsByTweetID(ctx, int64(req.TweetId))
	if err != nil {
		return nil, err
	}
	var pbComments []*pb.Comment
	for _, comment := range comments {
		pbComments = append(pbComments, &pb.Comment{
			Id:         int64(comment.ID),
			UserId:     int64(comment.UserID),
			TweetId:    int64(comment.TweetID),
			Content:    comment.Content,
			CreatedAt:  comment.CreatedAt.Format("2006-01-02 15:04:05"),
			Likes:      comment.Likes,
			LikesCount: int32(comment.LikesCount),
		})
	}
	return &pb.GetCommentsByTweetIDResponse{
		Comments: pbComments,
	}, nil
}

func (h *CommentHandler) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentResponse, error) {
	success, err := h.service.LikeComment(ctx, int64(req.UserId), int64(req.CommentId))
	if err != nil {
		return nil, err
	}
	return &pb.LikeCommentResponse{
		Success: success,
		Message: "Comment liked successfully",
	}, nil
}

func (h *CommentHandler) RemoveLikeFromComment(ctx context.Context, req *pb.RemoveLikeFromCommentRequest) (*pb.RemoveLikeFromCommentResponse, error) {
	success, err := h.service.RemoveLikeFromComment(ctx, int64(req.UserId), int64(req.CommentId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveLikeFromCommentResponse{
		Success: success,
		Message: "Like removed from comment successfully",
	}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	err := h.service.DeleteComment(ctx, int64(req.UserId), int64(req.TweetId), int64(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCommentResponse{
		Success: true,
		Message: "Comment deleted successfully",
	}, nil
}

func (h *CommentHandler) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	comment, err := h.service.GetCommentByID(ctx, int64(req.Id), int64(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.GetCommentResponse{
		Comment: &pb.Comment{
			Id:         int64(comment.ID),
			UserId:     int64(comment.UserID),
			TweetId:    int64(comment.TweetID),
			Content:    comment.Content,
			CreatedAt:  comment.CreatedAt.Format("2006-01-02 15:04:05"),
			Likes:      comment.Likes,
			LikesCount: int32(comment.LikesCount),
		},
	}, nil
}
