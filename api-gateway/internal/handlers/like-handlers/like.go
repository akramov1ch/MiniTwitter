package like_handlers

import (
	proto "api-gateway/protos/like-proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// LikeHandler goda like service uchun handler
type LikeHandler struct {
	LikeService proto.LikeServiceClient
}

// NewLikeHandler yangi LikeHandler yaratadi
func NewLikeHandler(likeService proto.LikeServiceClient) *LikeHandler {
	return &LikeHandler{LikeService: likeService}
}

// CreateLikeTweet - tweetga yoqtirish
// @Summary Like a tweet
// @Description Like a specific tweet by user
// @Tags likes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} proto.CreateLikeTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/tweet/{tweet_id}/{user_id} [post]
func (h *LikeHandler) CreateLikeTweet(c *gin.Context) {
	userIdStr := c.Param("user_id")
	tweetIdStr := c.Param("tweet_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	tweetId, err := strconv.ParseInt(tweetIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tweet ID"})
		return
	}

	req := &proto.CreateLikeTweetRequest{
		UserId:  int32(userId),
		TweetId: int32(tweetId),
	}

	resp, err := h.LikeService.CreateLikeTweet(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteLikeTweet - tweetdan yoqtirish
// @Summary Unlike a tweet
// @Description Remove like from a specific tweet by user
// @Tags likes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} proto.DeleteLikeTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/tweet/{tweet_id}/{user_id} [delete]
func (h *LikeHandler) DeleteLikeTweet(c *gin.Context) {
	userIdStr := c.Param("user_id")
	tweetIdStr := c.Param("tweet_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	tweetId, err := strconv.ParseInt(tweetIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tweet ID"})
		return
	}

	req := &proto.DeleteLikeTweetRequest{
		UserId:  int32(userId),
		TweetId: int32(tweetId),
	}

	resp, err := h.LikeService.DeleteLikeTweet(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetLikesTweet - tweetga yoqtirganlarni olish
// @Summary Get likes for a tweet
// @Description Get all likes for a specific tweet
// @Tags likes
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Param user_id_to_get path int true "User ID to get likes"
// @Success 200 {object} proto.GetLikesTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/tweet/{tweet_id}/{user_id_to_get} [get]
func (h *LikeHandler) GetLikesTweet(c *gin.Context) {
	tweetIdStr := c.Param("tweet_id")
	userIdToGetStr := c.Param("user_id_to_get")

	tweetId, err := strconv.ParseInt(tweetIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tweet ID"})
		return
	}

	userIdToGet, err := strconv.ParseInt(userIdToGetStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	req := &proto.GetLikesTweetRequest{
		TweetId:     int32(tweetId),
		UserIdToGet: int32(userIdToGet),
	}

	resp, err := h.LikeService.GetLikesTweet(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetLikeTweetByUser - foydalanuvchi tomonidan yoqtirilgan tweetlarni olish
// @Summary Get liked tweets by user
// @Description Get all tweets liked by a specific user
// @Tags likes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.GetLikeTweetByUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/user/{user_id} [get]
func (h *LikeHandler) GetLikeTweetByUser(c *gin.Context) {
	userIdStr := c.Param("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	req := &proto.GetLikeTweetByUserRequest{
		UserId: int32(userId),
	}

	resp, err := h.LikeService.GetLikeTweetByUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetLikedTweet - foydalanuvchi tomonidan yoqtirilgan tweetni olish
// @Summary Get a liked tweet
// @Description Get a specific liked tweet by user
// @Tags likes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param tweet_id path int true "Tweet ID"
// @Param like_identifier path string true "Like Identifier"
// @Success 200 {object} proto.GetLikedTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/tweet/{tweet_id}/{user_id}/{like_identifier} [get]
func (h *LikeHandler) GetLikedTweet(c *gin.Context) {
	userIdStr := c.Param("user_id")
	tweetIdStr := c.Param("tweet_id")
	likeIdentifier := c.Param("like_identifier")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	tweetId, err := strconv.ParseInt(tweetIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tweet ID"})
		return
	}

	req := &proto.GetLikedTweetRequest{
		UserId:        int32(userId),
		TweetId:      int32(tweetId),
		LikeIdentifier: likeIdentifier,
	}

	resp, err := h.LikeService.GetLikedTweet(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateLikeComment - kommentga yoqtirish
// @Summary Like a comment
// @Description Like a specific comment by user
// @Tags likes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} proto.CreateLikeCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/comment/{comment_id}/{user_id} [post]
func (h *LikeHandler) CreateLikeComment(c *gin.Context) {
	userIdStr := c.Param("user_id")
	commentIdStr := c.Param("comment_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	commentId, err := strconv.ParseInt(commentIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid comment ID"})
		return
	}

	req := &proto.CreateLikeCommentRequest{
		UserId:    int32(userId),
		CommentId: int32(commentId),
	}

	resp, err := h.LikeService.CreateLikeComment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteLikeComment - kommentdan yoqtirish
// @Summary Unlike a comment
// @Description Remove like from a specific comment by user
// @Tags likes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} proto.DeleteLikeCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/comment/{comment_id}/{user_id} [delete]
func (h *LikeHandler) DeleteLikeComment(c *gin.Context) {
	userIdStr := c.Param("user_id")
	commentIdStr := c.Param("comment_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	commentId, err := strconv.ParseInt(commentIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid comment ID"})
		return
	}

	req := &proto.DeleteLikeCommentRequest{
		UserId:    int32(userId),
		CommentId: int32(commentId),
	}

	resp, err := h.LikeService.DeleteLikeComment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetLikesComment - kommentga yoqtirganlarni olish
// @Summary Get likes for a comment
// @Description Get all likes for a specific comment
// @Tags likes
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} proto.GetLikesCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/comment/{comment_id} [get]
func (h *LikeHandler) GetLikesComment(c *gin.Context) {
	commentIdStr := c.Param("comment_id")

	commentId, err := strconv.ParseInt(commentIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid comment ID"})
		return
	}

	req := &proto.GetLikesCommentRequest{
		CommentId: int32(commentId),
	}

	resp, err := h.LikeService.GetLikesComment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetLikeCommentByUser - foydalanuvchi tomonidan yoqtirilgan kommentlarni olish
// @Summary Get likes for a comment by user
// @Description Get all likes for a specific comment by a user
// @Tags likes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.GetLikeCommentByUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /likes/comment/user/{user_id} [get]
func (h *LikeHandler) GetLikeCommentByUser(c *gin.Context) {
	userIdStr := c.Param("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	req := &proto.GetLikeCommentByUserRequest{
		UserId: int32(userId),
	}

	resp, err := h.LikeService.GetLikeCommentByUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
