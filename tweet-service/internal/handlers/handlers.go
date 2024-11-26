package handlers

import (
	"context"
	"tweet-service/internal/methods"
	"tweet-service/internal/models"
	"tweet-service/internal/service"
	pb "tweet-service/pkg/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type TweetHandler struct {
	pb.UnimplementedTweetServiceServer
	service service.TweetService
}

func NewTweetHandler(service service.TweetService) *TweetHandler {
	return &TweetHandler{service: service}
}

func (h *TweetHandler) CreateTweet(ctx context.Context, req *pb.CreateTweetRequest) (*pb.CreateTweetResponse, error) {
	userClient := methods.ConnectUser()
	user, err := userClient.GetUser(ctx, &pb.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	tweet := &models.Tweet{
		Content:  req.Content,
		UserID:   int64(req.UserId),
		Media:    req.Media,
		Username: user.User.Username,
	}

	err = h.service.CreateTweet(ctx, tweet)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTweetResponse{
		Success: true,
		Message: "Tweet created successfully",
	}, nil
}

func (h *TweetHandler) UpdateTweet(ctx context.Context, req *pb.UpdateTweetRequest) (*pb.UpdateTweetResponse, error) {
	tweet := &models.Tweet{
		ID:      int64(req.TweetId),
		Content: req.Content,
	}
	err := h.service.UpdateTweet(ctx, tweet)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTweetResponse{
		Success: true,
		Message: "Tweet updated successfully",
	}, nil
}

func (h *TweetHandler) DeleteTweet(ctx context.Context, req *pb.DeleteTweetRequest) (*pb.DeleteTweetResponse, error) {
	err := h.service.DeleteTweet(ctx, req.TweetId)
	if err != nil {
		return nil, err
	}
	userClient := methods.ConnectUser()
	_, err = userClient.RemoveTweet(ctx, &pb.RemoveTweetRequest{TweetId: req.TweetId, UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteTweetResponse{
		Success: true,
		Message: "Tweet deleted successfully",
	}, nil
}

func (h *TweetHandler) GetTweetByID(ctx context.Context, req *pb.GetTweetByIDRequest) (*pb.GetTweetByIDResponse, error) {
	tweet, err := h.service.GetTweetByID(ctx, req.TweetId, req.UserIdToGet)
	if err != nil {
		return nil, err
	}

	return &pb.GetTweetByIDResponse{
		Success: true,
		Message: "Tweet fetched successfully",
		Tweet: &pb.Tweet{
			Id:        int32(tweet.ID),
			Content:   tweet.Content,
			UserId:    int32(tweet.UserID),
			Username:  tweet.Username,
			CreatedAt: timestamppb.New(tweet.CreatedAt),
			Media:     tweet.Media,
			Likes:     tweet.Likes,
			Comments:  tweet.Comments,
			Shares:    tweet.Shares,
			Saves:     tweet.Saves,
		},
	}, nil
}

func (h *TweetHandler) GetSavedTweets(ctx context.Context, req *pb.GetSavedTweetsRequest) (*pb.GetSavedTweetsResponse, error) {
	tweets, err := h.service.GetSavedTweets(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	protoTweets := make([]*pb.Tweet, len(tweets))
	for i, tweet := range tweets {
		protoTweets[i] = &pb.Tweet{
			Id:        int32(tweet.ID),
			Content:   tweet.Content,
			UserId:    int32(tweet.UserID),
			Username:  tweet.Username,
			CreatedAt: timestamppb.New(tweet.CreatedAt),
			Media:     tweet.Media,
			Likes:     tweet.Likes,
			Comments:  tweet.Comments,
			Shares:    tweet.Shares,
			Saves:     tweet.Saves,
		}
	}

	return &pb.GetSavedTweetsResponse{
		Success: true,
		Message: "Saved tweets fetched successfully",
		Tweets:  protoTweets,
	}, nil
}

func (h *TweetHandler) GetLikedTweets(ctx context.Context, req *pb.GetLikedTweetsRequest) (*pb.GetLikedTweetsResponse, error) {
	tweets, err := h.service.GetLikedTweets(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	protoTweets := make([]*pb.Tweet, len(tweets))
	for i, tweet := range tweets {
		protoTweets[i] = &pb.Tweet{
			Id:        int32(tweet.ID),
			Content:   tweet.Content,
			UserId:    int32(tweet.UserID),
			Username:  tweet.Username,
			CreatedAt: timestamppb.New(tweet.CreatedAt),
			Media:     tweet.Media,
			Likes:     tweet.Likes,
			Comments:  tweet.Comments,
			Shares:    tweet.Shares,
			Saves:     tweet.Saves,
		}
	}

	return &pb.GetLikedTweetsResponse{
		Success: true,
		Message: "Liked tweets fetched successfully",
		Tweets:  protoTweets,
	}, nil
}

func (h *TweetHandler) GetTweetsByUser(ctx context.Context, req *pb.GetTweetsByUserRequest) (*pb.GetTweetsByUserResponse, error) {
	tweets, err := h.service.GetTweetsByUser(ctx, req.UserId, req.UserIdToGet)
	if err != nil {
		return nil, err
	}

	protoTweets := make([]*pb.Tweet, len(tweets))
	for i, tweet := range tweets {
		protoTweets[i] = &pb.Tweet{
			Id:        int32(tweet.ID),
			Content:   tweet.Content,
			UserId:    int32(tweet.UserID),
			Username:  tweet.Username,
			CreatedAt: timestamppb.New(tweet.CreatedAt),
			Media:     tweet.Media,
			Likes:     tweet.Likes,
			Comments:  tweet.Comments,
			Shares:    tweet.Shares,
			Saves:     tweet.Saves,
		}
	}

	return &pb.GetTweetsByUserResponse{
		Success: true,
		Message: "User tweets fetched successfully",
		Tweets:  protoTweets,
	}, nil
}

func (h *TweetHandler) AddLike(ctx context.Context, req *pb.AddLikeRequest) (*pb.AddLikeResponse, error) {
	err := h.service.AddLike(ctx, req.UserId, req.TweetId)
	if err != nil {
		return nil, err
	}

	return &pb.AddLikeResponse{
		Success: true,
		Message: "Tweet liked successfully",
	}, nil
}

func (h *TweetHandler) RemoveLike(ctx context.Context, req *pb.RemoveLikeRequest) (*pb.RemoveLikeResponse, error) {
	err := h.service.RemoveLike(ctx, req.UserId, req.TweetId)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveLikeResponse{
		Success: true,
		Message: "Like removed successfully",
	}, nil
}

func (h *TweetHandler) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	err := h.service.AddComment(ctx, req.UserId, req.TweetId, req.CommentId)
	if err != nil {
		return nil, err
	}

	return &pb.AddCommentResponse{
		Success: true,
		Message: "Comment added successfully",
	}, nil
}

func (h *TweetHandler) RemoveComment(ctx context.Context, req *pb.RemoveCommentRequest) (*pb.RemoveCommentResponse, error) {
	err := h.service.RemoveComment(ctx, req.UserId, req.TweetId, req.CommentId)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveCommentResponse{
		Success: true,
		Message: "Comment removed successfully",
	}, nil
}

func (h *TweetHandler) AddShare(ctx context.Context, req *pb.AddShareRequest) (*pb.AddShareResponse, error) {
	err := h.service.AddShare(ctx, req.UserId, req.TweetId)
	if err != nil {
		return nil, err
	}

	return &pb.AddShareResponse{
		Success: true,
		Message: "Tweet shared successfully",
	}, nil
}

func (h *TweetHandler) SaveTweet(ctx context.Context, req *pb.SaveTweetRequest) (*pb.SaveTweetResponse, error) {
	err := h.service.SaveTweet(ctx, req.UserId, req.TweetId)
	if err != nil {
		return nil, err
	}

	return &pb.SaveTweetResponse{
		Success: true,
		Message: "Tweet saved successfully",
	}, nil
}

func (h *TweetHandler) RemoveSave(ctx context.Context, req *pb.RemoveSaveRequest) (*pb.RemoveSaveResponse, error) {
	err := h.service.RemoveSave(ctx, req.UserId, req.TweetId)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveSaveResponse{
		Success: true,
		Message: "Save removed successfully",
	}, nil
}
