package comment_handlers

import (
	proto "api-gateway/protos/comment-proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Type aliases for swag
type CreateCommentResponse = proto.CreateCommentResponse
type GetCommentResponse = proto.GetCommentResponse
type DeleteCommentResponse = proto.DeleteCommentResponse
type GetCommentsByTweetIDResponse = proto.GetCommentsByTweetIDResponse
type LikeCommentResponse = proto.LikeCommentResponse
type RemoveLikeFromCommentResponse = proto.RemoveLikeFromCommentResponse

// CommentHandler goda comment service uchun handler
type CommentHandler struct {
	CommentService proto.CommentServiceClient
}

// NewCommentHandler yangi CommentHandler yaratadi
func NewCommentHandler(commentService proto.CommentServiceClient) *CommentHandler {
	return &CommentHandler{CommentService: commentService}
}

// CreateComment - yangi komment yaratish
// @Summary Create a comment
// @Description Create a new comment for a tweet
// @Tags comments
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param tweet_id path int true "Tweet ID"
// @Param content formData string true "Comment content"
// @Success 200 {object} proto.CreateCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{tweet_id}/{user_id} [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userIdStr := c.Param("user_id")
	tweetIdStr := c.Param("tweet_id")
	content := c.PostForm("content")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	tweetId, err := strconv.ParseInt(tweetIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tweet ID"})
		return
	}

	req := &proto.CreateCommentRequest{
		UserId:  userId,
		TweetId: tweetId,
		Content: content,
	}

	resp, err := h.CommentService.CreateComment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetComment - kommentni olish
// @Summary Get a comment
// @Description Get a comment by ID
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.GetCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{comment_id}/{user_id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	commentIdStr := c.Param("comment_id")
	userIdStr := c.Param("user_id")

	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid comment ID"})
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	req := &proto.GetCommentRequest{
		Id:     commentId,
		UserId: userId,
	}

	resp, err := h.CommentService.GetComment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteComment - kommentni o'chirish
// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} proto.DeleteCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{comment_id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	var req proto.DeleteCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.CommentService.DeleteComment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCommentsByTweetID - tweet bo'yicha kommentlarni olish
// @Summary Get comments by tweet ID
// @Description Get all comments for a specific tweet
// @Tags comments
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} proto.GetCommentsByTweetIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/tweet/{tweet_id} [get]
func (h *CommentHandler) GetCommentsByTweetID(c *gin.Context) {
	tweetIdStr := c.Param("tweet_id")

	tweetId, err := strconv.ParseInt(tweetIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tweet ID"})
		return
	}

	req := &proto.GetCommentsByTweetIDRequest{
		TweetId: tweetId,
	}

	resp, err := h.CommentService.GetCommentsByTweetID(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// LikeComment - kommentni yoqtirish
// @Summary Like a comment
// @Description Like a specific comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} proto.LikeCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/like/{comment_id} [post]
func (h *CommentHandler) LikeComment(c *gin.Context) {
	var req proto.LikeCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.CommentService.LikeComment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RemoveLikeFromComment - kommentdan yoqtirish
// @Summary Remove like from a comment
// @Description Remove like from a specific comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} proto.RemoveLikeFromCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/unlike/{comment_id} [post]
func (h *CommentHandler) RemoveLikeFromComment(c *gin.Context) {
	var req proto.RemoveLikeFromCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.CommentService.RemoveLikeFromComment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Define a new type for error responses
type ErrorResponse struct {
	Error string `json:"error"`
}
