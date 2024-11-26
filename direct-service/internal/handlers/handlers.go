package handlers

import (
	"context"
	"direct-service/internal/models"
	"direct-service/internal/service"
	pb "direct-service/pkg/proto"
)

type DirectHandler struct {
	pb.UnimplementedDirectServiceServer
	service service.DirectService
}

func NewDirectHandler(service service.DirectService) *DirectHandler {
	return &DirectHandler{service: service}
}

func (h *DirectHandler) CreateDirectMessage(ctx context.Context, req *pb.CreateDirectMessageRequest) (*pb.CreateDirectMessageResponse, error) {
	message := &models.DirectMessage{
		SenderID:   int64(req.SenderId),
		ReceiverID: int64(req.ReceiverId),
		TweetID:    int64(req.TweetId),
		Text:       req.Text,
		Media:      req.Media,
	}
	createdMessage, err := h.service.CreateDirectMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	return &pb.CreateDirectMessageResponse{
		Success: true,
		Message: "Direct message created successfully",
		DirectMessage: &pb.DirectMessage{
			Id:        int64(createdMessage.ID),
			SenderId:  int64(createdMessage.SenderID),
			ReceiverId: int64(createdMessage.ReceiverID),
			TweetId:    int64(createdMessage.TweetID),
			Text:       createdMessage.Text,
			Media:      createdMessage.Media,
			CreatedAt:  createdMessage.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (h *DirectHandler) GetDirectMessages(ctx context.Context, req *pb.GetDirectMessagesRequest) (*pb.GetDirectMessagesResponse, error) {
	directMessages, err := h.service.GetDirectMessages(ctx, int64(req.SenderId), int64(req.ReceiverId))
	if err != nil {
		return nil, err
	}
	protoDirectMessages := make([]*pb.DirectMessage, len(directMessages))
	for i, message := range directMessages {
		protoDirectMessages[i] = &pb.DirectMessage{
			Id:        int64(message.ID),
			SenderId:  int64(message.SenderID),
			ReceiverId: int64(message.ReceiverID),
			TweetId:    int64(message.TweetID),
			Text:       message.Text,
			Media:      message.Media,
			CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &pb.GetDirectMessagesResponse{
		Success: true,
		Message: "Direct messages fetched successfully",
		DirectMessages: protoDirectMessages,
	}, nil
}

func (h *DirectHandler) GetDirectMessageByID(ctx context.Context, req *pb.GetDirectMessageByIDRequest) (*pb.GetDirectMessageByIDResponse, error) {
	directMessage, err := h.service.GetDirectMessageByID(ctx, int64(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetDirectMessageByIDResponse{
		Success: true,
		Message: "Direct message fetched successfully",
		DirectMessage: &pb.DirectMessage{
			Id:        int64(directMessage.ID),
			SenderId:  int64(directMessage.SenderID),
			ReceiverId: int64(directMessage.ReceiverID),
			TweetId:    int64(directMessage.TweetID),
			Text:       directMessage.Text,
			Media:      directMessage.Media,
			CreatedAt:  directMessage.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (h *DirectHandler) DeleteDirectMessage(ctx context.Context, req *pb.DeleteDirectMessageRequest) (*pb.DeleteDirectMessageResponse, error) {
	success, err := h.service.DeleteDirectMessage(ctx, int64(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteDirectMessageResponse{
		Success: success,
		Message: "Direct message deleted successfully",
	}, nil
}
