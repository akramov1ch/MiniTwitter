package user_handlers

import (
	proto "api-gateway/protos/user-proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// UserHandler goda user service uchun handler
type UserHandler struct {
	UserService   proto.UserServiceClient
	minioClient   *minio.Client
	bucketName    string
	minioEndpoint string
}

// NewUserHandler yangi UserHandler yaratadi
func NewUserHandler(userService proto.UserServiceClient, minioEndpoint, accessKeyID, secretAccessKey, bucketName string) (*UserHandler, error) {
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return &UserHandler{
		UserService:   userService,
		minioClient:   minioClient,
		bucketName:    bucketName,
		minioEndpoint: minioEndpoint,
	}, nil
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {object} proto.RegisterResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	req := &proto.RegisterRequest{
		Username: c.PostForm("username"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	resp, err := h.UserService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// VerifyCode - kodni tasdiqlash
// @Summary Verify code
// @Description Verify the code sent to the user
// @Tags users
// @Accept json
// @Produce json
// @Param login_identifier formData string true "Login Identifier"
// @Param code formData string true "Verification Code"
// @Success 200 {object} proto.VerifyCodeResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/verify [post]
func (h *UserHandler) VerifyCode(c *gin.Context) {
	loginIdentifier := c.PostForm("login_identifier")
	code := c.PostForm("code")

	req := &proto.VerifyCodeRequest{
		LoginIdentifier: loginIdentifier,
		Code: code,
	}

	resp, err := h.UserService.VerifyCode(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Login - foydalanuvchini tizimga kiritish
// @Summary User login
// @Description Login a user with their credentials
// @Tags users
// @Accept json
// @Produce json
// @Param login_identifier formData string true "Login Identifier"
// @Param password formData string true "Password"
// @Success 200 {object} proto.LoginResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	req := &proto.LoginRequest{
		LoginIdentifier: c.PostForm("login_identifier"),
		Password:        c.PostForm("password"),
	}

	resp, err := h.UserService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUser - foydalanuvchini olish
// @Summary Get user details
// @Description Get details of a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param user_id_to_get formData int true "User ID to get"
// @Success 200 {object} proto.GetUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/get [post]
func (h *UserHandler) GetUser(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	userIdToGetStr := c.PostForm("user_id_to_get")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	userIdToGet, err := strconv.ParseInt(userIdToGetStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id_to_get"})
		return
	}

	req := &proto.GetUserRequest{
		UserId:      int32(userId),
		UserIdToGet: int32(userIdToGet),
	}

	resp, err := h.UserService.GetUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserByUsername - foydalanuvchini username bo'yicha olish
// @Summary Get user by username
// @Description Get user details by username
// @Tags users
// @Accept json
// @Produce json
// @Param username formData string true "Username"
// @Success 200 {object} proto.GetUserByUsernameResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/username [post]
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.PostForm("username")

	req := &proto.GetUserByUsernameRequest{
		Username: username,
	}

	resp, err := h.UserService.GetUserByUsername(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserByEmail - foydalanuvchini email bo'yicha olish
// @Summary Get user by email
// @Description Get user details by email
// @Tags users
// @Accept json
// @Produce json
// @Param email formData string true "Email"
// @Success 200 {object} proto.GetUserByEmailResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/email [post]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.PostForm("email")

	req := &proto.GetUserByEmailRequest{
		Email: email,
	}

	resp, err := h.UserService.GetUserByEmail(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserByPhoneNumber - foydalanuvchini telefon raqami bo'yicha olish
// @Summary Get user by phone number
// @Description Get user details by phone number
// @Tags users
// @Accept json
// @Produce json
// @Param phone formData string true "Phone Number"
// @Success 200 {object} proto.GetUserByPhoneResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/phone [post]
func (h *UserHandler) GetUserByPhoneNumber(c *gin.Context) {
	phone := c.PostForm("phone")

	req := &proto.GetUserByPhoneRequest{
		Phone: phone,
	}

	resp, err := h.UserService.GetUserByPhone(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUser - foydalanuvchini yangilash
// @Summary Update user details
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param username formData string false "Username"
// @Param email formData string false "Email"
// @Param phone formData string false "Phone"
// @Param bio formData string false "Bio"
// @Param name formData string false "Name"
// @Param is_private formData bool false "Is Private"
// @Success 200 {object} proto.UpdateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/update [post]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	username := c.PostForm("username")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	bio := c.PostForm("bio")
	name := c.PostForm("name")
	isPrivate := c.PostForm("is_private")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	isPrivateBool, err := strconv.ParseBool(isPrivate)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid is_private value"})
		return
	}

	req := &proto.UpdateUserRequest{
		UserId:    int32(userId),
		Username:  username,
		Email:     email,
		Phone:     phone,
		Bio:       bio,
		Name:      name,
		IsPrivate: isPrivateBool,
	}

	resp, err := h.UserService.UpdateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteUser - foydalanuvchini o'chirish
// @Summary Delete a user
// @Description Delete a specific user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Success 200 {object} proto.DeleteUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/delete [post]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	req := &proto.DeleteUserRequest{
		UserId: int32(userId),
	}

	resp, err := h.UserService.DeleteUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Logout - foydalanuvchini tizimdan chiqish
// @Summary User logout
// @Description Logout a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Success 200 {object} proto.LogoutResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	req := &proto.LogoutRequest{
		UserId: int32(userId),
	}

	resp, err := h.UserService.Logout(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ChangePassword - foydalanuvchi parolini o'zgartirish
// @Summary Change user password
// @Description Change the password for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param old_password formData string true "Old Password"
// @Param new_password formData string true "New Password"
// @Success 200 {object} proto.ChangePasswordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	req := &proto.ChangePasswordRequest{
		UserId:      int32(userId),
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	resp, err := h.UserService.ChangePassword(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ForgotPassword - foydalanuvchi parolini tiklash
// @Summary Forgot password
// @Description Reset the password for a user
// @Tags users
// @Accept json
// @Produce json
// @Param login_identifier formData string true "Login Identifier"
// @Param new_password formData string true "New Password"
// @Success 200 {object} proto.ForgotPasswordResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/forgot-password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	loginIdentifier := c.PostForm("login_identifier")
	newPassword := c.PostForm("new_password")

	req := &proto.ForgotPasswordRequest{
		LoginIdentifier: loginIdentifier,
		NewPassword:     newPassword,
	}

	resp, err := h.UserService.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddTweet - foydalanuvchiga tweet qo'shish
// @Summary Add a tweet to user
// @Description Add a specific tweet to a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param tweet_id formData int true "Tweet ID"
// @Success 200 {object} proto.AddTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/add-tweet [post]
func (h *UserHandler) AddTweet(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	tweetIdStr := c.PostForm("tweet_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	tweetId, err := strconv.ParseInt(tweetIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tweet_id"})
		return
	}

	req := &proto.AddTweetRequest{
		UserId:  int32(userId),
		TweetId: int32(tweetId),
	}

	resp, err := h.UserService.AddTweet(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RemoveTweet - foydalanuvchidan tweetni o'chirish
// @Summary Remove a tweet from user
// @Description Remove a specific tweet from a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param tweet_id formData int true "Tweet ID"
// @Success 200 {object} proto.RemoveTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/remove-tweet [post]
func (h *UserHandler) RemoveTweet(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	tweetIdStr := c.PostForm("tweet_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	tweetId, err := strconv.ParseInt(tweetIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tweet_id"})
		return
	}

	req := &proto.RemoveTweetRequest{
		UserId:  int32(userId),
		TweetId: int32(tweetId),
	}

	resp, err := h.UserService.RemoveTweet(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddFollower - foydalanuvchini kuzatish
// @Summary Add a follower to user
// @Description Add a specific follower to a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param follower_id formData int true "Follower ID"
// @Success 200 {object} proto.AddFollowerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/add-follower [post]
func (h *UserHandler) AddFollower(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	followerIdStr := c.PostForm("follower_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	followerId, err := strconv.ParseInt(followerIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid follower_id"})
		return
	}

	req := &proto.AddFollowerRequest{
		UserId:     int32(userId),
		FollowerId: int32(followerId),
	}

	resp, err := h.UserService.AddFollower(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RemoveFollower - foydalanuvchidan kuzatuvchini o'chirish
// @Summary Remove a follower from user
// @Description Remove a specific follower from a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param follower_id formData int true "Follower ID"
// @Success 200 {object} proto.RemoveFollowerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/remove-follower [post]
func (h *UserHandler) RemoveFollower(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	followerIdStr := c.PostForm("follower_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	followerId, err := strconv.ParseInt(followerIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid follower_id"})
		return
	}

	req := &proto.RemoveFollowerRequest{
		UserId:     int32(userId),
		FollowerId: int32(followerId),
	}

	resp, err := h.UserService.RemoveFollower(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddFollowing - foydalanuvchini kuzatish
// @Summary Add a following to user
// @Description Add a specific following to a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param following_id formData int true "Following ID"
// @Success 200 {object} proto.AddFollowingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/add-following [post]
func (h *UserHandler) AddFollowing(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	followingIdStr := c.PostForm("following_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	followingId, err := strconv.ParseInt(followingIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid following_id"})
		return
	}

	req := &proto.AddFollowingRequest{
		UserId:     int32(userId),
		FollowingId: int32(followingId),
	}

	resp, err := h.UserService.AddFollowing(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RemoveFollowing - foydalanuvchidan kuzatishni o'chirish
// @Summary Remove a following from user
// @Description Remove a specific following from a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param following_id formData int true "Following ID"
// @Success 200 {object} proto.RemoveFollowingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/remove-following [post]
func (h *UserHandler) RemoveFollowing(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	followingIdStr := c.PostForm("following_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	followingId, err := strconv.ParseInt(followingIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid following_id"})
		return
	}

	req := &proto.RemoveFollowingRequest{
		UserId:     int32(userId),
		FollowingId: int32(followingId),
	}

	resp, err := h.UserService.RemoveFollowing(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AcceptFollowRequest - kuzatuvchi so'rovini qabul qilish
// @Summary Accept a follow request
// @Description Accept a specific follow request
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param follower_id formData int true "Follower ID"
// @Success 200 {object} proto.AcceptFollowRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/accept-follow [post]
func (h *UserHandler) AcceptFollowRequest(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	followerIdStr := c.PostForm("follower_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	followerId, err := strconv.ParseInt(followerIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid follower_id"})
		return
	}

	req := &proto.AcceptFollowRequestRequest{
		UserId:     int32(userId),
		FollowerId: int32(followerId),
	}

	resp, err := h.UserService.AcceptFollowRequest(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RejectFollowRequest - kuzatuvchi so'rovini rad etish
// @Summary Reject a follow request
// @Description Reject a specific follow request
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param follower_id formData int true "Follower ID"
// @Success 200 {object} proto.RejectFollowRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/reject-follow [post]
func (h *UserHandler) RejectFollowRequest(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	followerIdStr := c.PostForm("follower_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	followerId, err := strconv.ParseInt(followerIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid follower_id"})
		return
	}

	req := &proto.RejectFollowRequestRequest{
		UserId:     int32(userId),
		FollowerId: int32(followerId),
	}

	resp, err := h.UserService.RejectFollowRequest(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// BlockUser - foydalanuvchini bloklash
// @Summary Block a user
// @Description Block a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param blocked_user_id formData int true "Blocked User ID"
// @Success 200 {object} proto.BlockUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/block [post]
func (h *UserHandler) BlockUser(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	blockedUserIdStr := c.PostForm("blocked_user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	blockedUserId, err := strconv.ParseInt(blockedUserIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid blocked_user_id"})
		return
	}

	req := &proto.BlockUserRequest{
		UserId:        int32(userId),
		BlockedUserId: int32(blockedUserId),
	}

	resp, err := h.UserService.BlockUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UnblockUser - foydalanuvchini blokdan chiqarish
// @Summary Unblock a user
// @Description Unblock a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param blocked_user_id formData int true "Blocked User ID"
// @Success 200 {object} proto.UnblockUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/unblock [post]
func (h *UserHandler) UnblockUser(c *gin.Context) {
	userIdStr := c.PostForm("user_id")
	blockedUserIdStr := c.PostForm("blocked_user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user_id"})
		return
	}

	blockedUserId, err := strconv.ParseInt(blockedUserIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid blocked_user_id"})
		return
	}

	req := &proto.UnblockUserRequest{
		UserId:        int32(userId),
		BlockedUserId: int32(blockedUserId),
	}

	resp, err := h.UserService.UnblockUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddAvatar - foydalanuvchi avatarini qo'shish
// @Summary Add avatar to user
// @Description Add an avatar for a specific user
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData int true "User ID"
// @Param avatar formData file true "Avatar File"
// @Success 200 {object} proto.AddAvatarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/add-avatar [post]
func (h *UserHandler) AddAvatar(c *gin.Context) {
	var req proto.AddAvatarRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "File not found"})
		return
	}

	objectName := file.Filename
	_, err = h.minioClient.FPutObject(c.Request.Context(), h.bucketName, objectName, file.Filename, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to upload file"})
		return
	}

	avatarURL := "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + objectName

	req.AvatarUrl = avatarURL

	resp, err := h.UserService.AddAvatar(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateAvatar - foydalanuvchi avatarini yangilash
// @Summary Update user avatar
// @Description Update the avatar for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Param avatar formData file true "New Avatar File"
// @Success 200 {object} proto.UpdateAvatarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/update-avatar [post]
func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	var req proto.UpdateAvatarRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}
	oldAvatarName := req.AvatarUrl
	if err := h.minioClient.RemoveObject(c.Request.Context(), h.bucketName, oldAvatarName, minio.RemoveObjectOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete old avatar"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "File not found"})
		return
	}

	objectName := file.Filename
	_, err = h.minioClient.FPutObject(c.Request.Context(), h.bucketName, objectName, file.Filename, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to upload new avatar"})
		return
	}

	avatarURL := "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + objectName

	req.AvatarUrl = avatarURL

	resp, err := h.UserService.UpdateAvatar(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAvatar - foydalanuvchi avatarini olish
// @Summary Get user avatar
// @Description Retrieve the avatar for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Success 200 {object} proto.GetAvatarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/get-avatar [post]
func (h *UserHandler) GetAvatar(c *gin.Context) {
	var req proto.GetAvatarRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	avatarName, err := h.UserService.GetAvatar(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get avatar name"})
		return
	}

	avatarURL := "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + avatarName.Avatar.Url

	resp := &proto.GetAvatarResponse{
		Success: true,
		Message: "Avatar retrieved successfully",
		Avatar: &proto.Avatar{
			Url: avatarURL,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteAvatar - foydalanuvchi avatarini o'chirish
// @Summary Delete user avatar
// @Description Delete the avatar for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id formData int true "User ID"
// @Success 200 {object} proto.RemoveAvatarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/delete-avatar [post]
func (h *UserHandler) DeleteAvatar(c *gin.Context) {
	var req proto.RemoveAvatarRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	avatarName, err := h.UserService.GetAvatar(c.Request.Context(), &proto.GetAvatarRequest{UserId: req.UserId}) // Avatar nomini olish
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get avatar name"})
		return
	}

	if err := h.minioClient.RemoveObject(c.Request.Context(), h.bucketName, avatarName.Avatar.Url, minio.RemoveObjectOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete avatar"})
		return
	}

	resp := &proto.RemoveAvatarResponse{
		Success: true,
		Message: "Avatar deleted successfully",
	}

	c.JSON(http.StatusOK, resp)
}
