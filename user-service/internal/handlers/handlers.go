package handlers

import (
	"context"
	"errors"
	"time"
	"user-service/internal/models"
	"user-service/internal/service"
	pb "user-service/pkg/proto"
	"user-service/utils"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	if req.Email == "" && req.Phone == "" {
		return nil, errors.New("email or phone is required")
	}
	
	
	message := utils.GenerateCode()
	if req.Email != "" {
		utils.SendEmail(req.Email, message, "Verification Code")
	} else if req.Phone != "" {
		utils.SendSMS(req.Phone, message)
	}

	user, err := h.service.RegisterUser(ctx, req.Username, req.Email, req.Phone, req.Name, password, message)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		User: &pb.User{
			Id:             int32(user.ID),
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Phone:          user.Phone,
			CreatedAt:      timestamppb.New(user.CreatedAt),
			Bio:            user.Bio,
			TweetCount:     int32(user.TweetCount),
			FollowerCount:  int32(user.FollowerCount),
			FollowingCount: int32(user.FollowingCount),
			Tweets:         user.Tweets,
			Followers:      user.Followers,
			Following:      user.Following,
			IsPrivate:      user.IsPrivate,
		},
	}, nil
}

func (h *UserHandler) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.service.LoginUser(ctx, req.LoginIdentifier)
	if err != nil {
		return nil, err
	}

	err = utils.VerifyPassword(user.Password, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(int(user.ID), user.Username)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		Success: true,
		Message: "Login successful",
		User: &pb.User{
			Id:             int32(user.ID),
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Phone:          user.Phone,
			CreatedAt:      timestamppb.New(user.CreatedAt),
			Bio:            user.Bio,
			TweetCount:     int32(user.TweetCount),
			FollowerCount:  int32(user.FollowerCount),
			FollowingCount: int32(user.FollowingCount),
			Tweets:         user.Tweets,
			Followers:      user.Followers,
			Following:      user.Following,
			IsPrivate:      user.IsPrivate,
		},
		Token: token,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.service.GetUser(ctx, int64(req.UserId), int64(req.UserIdToGet))
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Message: "User fetched successfully",
		User: &pb.User{
			Id:             int32(user.ID),
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Phone:          user.Phone,
			CreatedAt:      timestamppb.New(user.CreatedAt),
			Bio:            user.Bio,
			TweetCount:     int32(user.TweetCount),
			FollowerCount:  int32(user.FollowerCount),
			FollowingCount: int32(user.FollowingCount),
			Tweets:         user.Tweets,
			Followers:      user.Followers,
			Following:      user.Following,
			IsPrivate:      user.IsPrivate,
		},
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, err := h.service.UpdateUser(ctx, models.User{
		ID:        int64(req.UserId),
		Username:  req.Username,
		Email:     req.Email,
		Name:      req.Name,
		Phone:     req.Phone,
		Bio:       req.Bio,
		IsPrivate: req.IsPrivate,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{
		Message: "User updated successfully",
		User: &pb.User{
			Id:         int32(user.ID),
			Username:   user.Username,
			Email:      user.Email,
			Name:       user.Name,
			Phone:      user.Phone,
			CreatedAt:  timestamppb.New(user.CreatedAt),
			Bio:        user.Bio,
			TweetCount: int32(user.TweetCount),
			IsPrivate:  user.IsPrivate,
			Tweets:     user.Tweets,
			Followers:  user.Followers,
			Following:  user.Following,
		},
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.service.DeleteUser(ctx, int64(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}

func (h *UserHandler) FollowUser(ctx context.Context, req *pb.AddFollowingRequest) (*pb.AddFollowingResponse, error) {
	user, err := h.service.GetUser(ctx, int64(req.FollowingId), int64(req.UserId))
	if err != nil {
		return nil, err
	}												

	if user.IsPrivate {
		err = h.service.SendFollowRequest(ctx, int64(req.UserId), int64(req.FollowingId))
		if err != nil {
			return nil, err
		}
		return &pb.AddFollowingResponse{Message: "Follow request sent successfully"}, nil
	}

	err = h.service.AddFollowing(ctx, int64(req.UserId), int64(req.FollowingId))
	if err != nil {
		return nil, err
	}
	return &pb.AddFollowingResponse{Message: "User followed successfully"}, nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	err := h.service.ChangePassword(ctx, int64(req.UserId), req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &pb.ChangePasswordResponse{Message: "Password changed successfully"}, nil
}

func (h *UserHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	err := h.service.ForgotPassword(ctx, req.NewPassword, req.LoginIdentifier)
	if err != nil {
		return nil, err
	}
	return &pb.ForgotPasswordResponse{Message: "Password reset successfully"}, nil
}

func (h *UserHandler) UnfollowUser(ctx context.Context, req *pb.RemoveFollowingRequest) (*pb.RemoveFollowingResponse, error) {
	err := h.service.RemoveFollowing(ctx, int64(req.UserId), int64(req.FollowingId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveFollowingResponse{Message: "User unfollowed successfully"}, nil
}

func (h *UserHandler) AddFollower(ctx context.Context, req *pb.AddFollowerRequest) (*pb.AddFollowerResponse, error) {
	err := h.service.AddFollower(ctx, int64(req.UserId), int64(req.FollowerId))
	if err != nil {
		return nil, err
	}
	return &pb.AddFollowerResponse{Message: "Follower added successfully"}, nil
}

func (h *UserHandler) RemoveFollower(ctx context.Context, req *pb.RemoveFollowerRequest) (*pb.RemoveFollowerResponse, error) {
	err := h.service.RemoveFollower(ctx, int64(req.UserId), int64(req.FollowerId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveFollowerResponse{Message: "Follower removed successfully"}, nil
}

func (h *UserHandler) AddAvatar(ctx context.Context, req *pb.AddAvatarRequest) (*pb.AddAvatarResponse, error) {
	err := h.service.AddAvatar(ctx, int64(req.UserId), req.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return &pb.AddAvatarResponse{Message: "Avatar added successfully"}, nil
}

func (h *UserHandler) GetAvatar(ctx context.Context, req *pb.GetAvatarRequest) (*pb.GetAvatarResponse, error) {
	avatar, err := h.service.GetAvatar(ctx, int64(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.GetAvatarResponse{
		Success: true,
		Message: "Avatar fetched successfully",
		Avatar: &pb.Avatar{
			UserId: int64(req.UserId),
			Url:    avatar,
		},
	}, nil
}

func (h *UserHandler) UpdateAvatar(ctx context.Context, req *pb.UpdateAvatarRequest) (*pb.UpdateAvatarResponse, error) {
	err := h.service.UpdateAvatar(ctx, int64(req.UserId), req.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateAvatarResponse{Message: "Avatar updated successfully"}, nil
}

func (h *UserHandler) RemoveAvatar(ctx context.Context, req *pb.RemoveAvatarRequest) (*pb.RemoveAvatarResponse, error) {
	err := h.service.DeleteAvatar(ctx, int64(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveAvatarResponse{Message: "Avatar deleted successfully"}, nil
}

func (h *UserHandler) AddTweet(ctx context.Context, req *pb.AddTweetRequest) (*pb.AddTweetResponse, error) {
	err := h.service.AddTweet(ctx, int64(req.UserId), int64(req.TweetId))
	if err != nil {
		return nil, err
	}
	return &pb.AddTweetResponse{Message: "Tweet added successfully"}, nil
}

func (h *UserHandler) RemoveTweet(ctx context.Context, req *pb.RemoveTweetRequest) (*pb.RemoveTweetResponse, error) {
	err := h.service.RemoveTweet(ctx, int64(req.UserId), int64(req.TweetId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveTweetResponse{Message: "Tweet removed successfully"}, nil
}

func (h *UserHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	user, err := h.service.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByUsernameResponse{
		Message: "User fetched successfully",
		User: &pb.User{
			Id:             int32(user.ID),
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Phone:          user.Phone,
			CreatedAt:      timestamppb.New(user.CreatedAt),
			Bio:            user.Bio,
			TweetCount:     int32(user.TweetCount),
			FollowerCount:  int32(user.FollowerCount),
			FollowingCount: int32(user.FollowingCount),
			Tweets:         user.Tweets,
			Followers:      user.Followers,
			Following:      user.Following,
			IsPrivate:      user.IsPrivate,
		},
	}, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	user, err := h.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByEmailResponse{
		Message: "User fetched successfully",
		User: &pb.User{
			Id:             int32(user.ID),
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Phone:          user.Phone,
			CreatedAt:      timestamppb.New(user.CreatedAt),
			Bio:            user.Bio,
			TweetCount:     int32(user.TweetCount),
			FollowerCount:  int32(user.FollowerCount),
			FollowingCount: int32(user.FollowingCount),
			Tweets:         user.Tweets,
			Followers:      user.Followers,
			Following:      user.Following,
			IsPrivate:      user.IsPrivate,
		},
	}, nil
}

func (h *UserHandler) GetUserByPhone(ctx context.Context, req *pb.GetUserByPhoneRequest) (*pb.GetUserByPhoneResponse, error) {
	user, err := h.service.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByPhoneResponse{
		Message: "User fetched successfully",
		User: &pb.User{
			Id:             int32(user.ID),
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Phone:          user.Phone,
			CreatedAt:      timestamppb.New(user.CreatedAt),
			Bio:            user.Bio,
			TweetCount:     int32(user.TweetCount),
			FollowerCount:  int32(user.FollowerCount),
			FollowingCount: int32(user.FollowingCount),
			Tweets:         user.Tweets,
			Followers:      user.Followers,
			Following:      user.Following,
			IsPrivate:      user.IsPrivate,
		},
	}, nil
}

func (h *UserHandler) AcceptFollowRequest(ctx context.Context, req *pb.AcceptFollowRequestRequest) (*pb.AcceptFollowRequestResponse, error) {
	err := h.service.AcceptFollowRequest(ctx, int64(req.UserId), int64(req.FollowerId))
	if err != nil {
		return nil, err
	}
	return &pb.AcceptFollowRequestResponse{Message: "Follow request accepted successfully"}, nil
}

func (h *UserHandler) RejectFollowRequest(ctx context.Context, req *pb.RejectFollowRequestRequest) (*pb.RejectFollowRequestResponse, error) {
	err := h.service.RejectFollowRequest(ctx, int64(req.UserId), int64(req.FollowerId))
	if err != nil {
		return nil, err
	}
	return &pb.RejectFollowRequestResponse{Message: "Follow request rejected successfully"}, nil
}

func (h *UserHandler) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.VerifyCodeResponse, error) {
	code, err := h.service.VerifyCode(ctx, req.LoginIdentifier, req.Code)
	if err != nil {
		return nil, err
	}
	halfHour := 30 * time.Minute

	user, err := h.service.GetUserByEmail(ctx, req.LoginIdentifier)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		user, err = h.service.GetUserByPhone(ctx, req.LoginIdentifier)
		if err != nil {
			return nil, err
		}
	}

	if time.Now().After(code.CreatedAt.Add(halfHour)) {
		if user.ID != 0 {
			err = h.service.DeleteUser(ctx, user.ID)
			if err != nil {
				return nil, err
			}
		}
		return nil, errors.New("code expired")
	}

	return &pb.VerifyCodeResponse{
		Success: true,
		Message: "Code verified successfully",
	}, nil
}
