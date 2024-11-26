package tweet_handlers

import (
	proto "api-gateway/protos/tweet-proto"
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

// TweetHandler goda tweet service uchun handler
type TweetHandler struct {
	TweetService  proto.TweetServiceClient
	minioClient   *minio.Client
	bucketName    string
	minioEndpoint string
}

// NewTweetHandler yangi TweetHandler yaratadi
func NewTweetHandler(tweetService proto.TweetServiceClient, minioEndpoint, accessKeyID, secretAccessKey, bucketName string) (*TweetHandler, error) {
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return &TweetHandler{
		TweetService:  tweetService,
		minioClient:   minioClient,
		bucketName:    bucketName,
		minioEndpoint: minioEndpoint,
	}, nil
}

// CreateTweet - yangi tweet yaratish
// @Summary Create a tweet
// @Description Create a new tweet for a user
// @Tags tweets
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param content formData string true "Tweet content"
// @Param media formData file false "Media files"
// @Success 200 {object} proto.CreateTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/{user_id} [post]
func (h *TweetHandler) CreateTweet(c *gin.Context) {
	userIdStr := c.Param("user_id")
	content := c.PostForm("content")

	var media []string
	if mediaFiles, err := c.MultipartForm(); err == nil {
		for _, file := range mediaFiles.File["media"] {
			mediaName := file.Filename
			if err := c.SaveUploadedFile(file, mediaName); err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save media file"})
				return
			}

			_, err := h.minioClient.FPutObject(c.Request.Context(), h.bucketName, mediaName, mediaName, minio.PutObjectOptions{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to upload media to MinIO"})
				return
			}

			mediaURL := "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + mediaName
			media = append(media, mediaURL)
		}
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	req := &proto.CreateTweetRequest{
		UserId:  int32(userId),
		Content: content,
		Media:   media,
	}

	resp, err := h.TweetService.CreateTweet(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTweetsByUser - foydalanuvchi bo'yicha tweetlarni olish
// @Summary Get tweets by user
// @Description Get all tweets for a specific user
// @Tags tweets
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.GetTweetsByUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/user/{user_id} [get]
func (h *TweetHandler) GetTweetsByUser(c *gin.Context) {
	var req proto.GetTweetsByUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.GetTweetsByUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i := range resp.Tweets {
		for j := range resp.Tweets[i].Media {
			mediaName := resp.Tweets[i].Media[j]
			resp.Tweets[i].Media[j] = "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + mediaName
		}
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateTweet - tweetni yangilash
// @Summary Update a tweet
// @Description Update an existing tweet
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} proto.UpdateTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/{tweet_id} [put]
func (h *TweetHandler) UpdateTweet(c *gin.Context) {
	var req proto.UpdateTweetRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.UpdateTweet(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTweet - tweetni o'chirish
// @Summary Delete a tweet
// @Description Delete a specific tweet by ID
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} proto.DeleteTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/{tweet_id} [delete]
func (h *TweetHandler) DeleteTweet(c *gin.Context) {
	var req proto.DeleteTweetRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	tweet, err := h.TweetService.GetTweetByID(c.Request.Context(), &proto.GetTweetByIDRequest{
		TweetId: req.TweetId,
		UserIdToGet: req.UserId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get tweet"})
		return
	}

	for _, media := range tweet.Tweet.Media {
		if err := h.minioClient.RemoveObject(c.Request.Context(), h.bucketName, media, minio.RemoveObjectOptions{}); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete media from MinIO"})
			return
		}
	}

	resp, err := h.TweetService.DeleteTweet(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTweetByID - tweetni ID bo'yicha olish
// @Summary Get a tweet by ID
// @Description Get a specific tweet by ID
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Success 200 {object} proto.GetTweetByIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/{tweet_id} [get]
func (h *TweetHandler) GetTweetByID(c *gin.Context) {
	var req proto.GetTweetByIDRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.GetTweetByID(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if resp.Tweet != nil {
		for j := range resp.Tweet.Media {
			mediaName := resp.Tweet.Media[j]
			resp.Tweet.Media[j] = "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + mediaName
		}
	}

	c.JSON(http.StatusOK, resp)
}

// GetSavedTweets - saqlangan tweetlarni olish
// @Summary Get saved tweets
// @Description Get all saved tweets for a user
// @Tags tweets
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.GetSavedTweetsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/saved/{user_id} [get]
func (h *TweetHandler) GetSavedTweets(c *gin.Context) {
	var req proto.GetSavedTweetsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.GetSavedTweets(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i := range resp.Tweets {
		for j := range resp.Tweets[i].Media {
			mediaName := resp.Tweets[i].Media[j]
			resp.Tweets[i].Media[j] = "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + mediaName
		}
	}

	c.JSON(http.StatusOK, resp)
}

// GetLikedTweets - yoqtirilgan tweetlarni olish
// @Summary Get liked tweets
// @Description Get all tweets liked by a user
// @Tags tweets
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.GetLikedTweetsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/liked/{user_id} [get]
func (h *TweetHandler) GetLikedTweets(c *gin.Context) {
	var req proto.GetLikedTweetsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.GetLikedTweets(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i := range resp.Tweets {
		for j := range resp.Tweets[i].Media {
			mediaName := resp.Tweets[i].Media[j]
			resp.Tweets[i].Media[j] = "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + mediaName
		}
	}

	c.JSON(http.StatusOK, resp)
}

// AddLike - tweetga yoqtirish
// @Summary Like a tweet
// @Description Add a like to a specific tweet
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.AddLikeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/like/{tweet_id}/{user_id} [post]
func (h *TweetHandler) AddLike(c *gin.Context) {
	var req proto.AddLikeRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.AddLike(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RemoveLike - tweetdan yoqtirish
// @Summary Unlike a tweet
// @Description Remove a like from a specific tweet
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.RemoveLikeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/unlike/{tweet_id}/{user_id} [delete]
func (h *TweetHandler) RemoveLike(c *gin.Context) {
	var req proto.RemoveLikeRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.RemoveLike(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddComment - tweetga komment qo'shish
// @Summary Add a comment to a tweet
// @Description Add a comment to a specific tweet
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Param user_id path int true "User ID"
// @Param content formData string true "Comment content"
// @Success 200 {object} proto.AddCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/comment/{tweet_id}/{user_id} [post]
func (h *TweetHandler) AddComment(c *gin.Context) {
	var req proto.AddCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.AddComment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RemoveComment - tweetdan kommentni o'chirish
// @Summary Remove a comment from a tweet
// @Description Remove a specific comment from a tweet
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} proto.RemoveCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/comment/{tweet_id}/{comment_id} [delete]
func (h *TweetHandler) RemoveComment(c *gin.Context) {
	var req proto.RemoveCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.RemoveComment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddShare - tweetni ulash
// @Summary Share a tweet
// @Description Share a specific tweet
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.AddShareResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/share/{tweet_id}/{user_id} [post]
func (h *TweetHandler) AddShare(c *gin.Context) {
	var req proto.AddShareRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.AddShare(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SaveTweet - tweetni saqlash
// @Summary Save a tweet
// @Description Save a specific tweet for later
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.SaveTweetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/save/{tweet_id}/{user_id} [post]
func (h *TweetHandler) SaveTweet(c *gin.Context) {
	var req proto.SaveTweetRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.SaveTweet(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RemoveSave - tweetni saqlashdan o'chirish
// @Summary Remove saved tweet
// @Description Remove a specific tweet from saved tweets
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet_id path int true "Tweet ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} proto.RemoveSaveResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets/remove-save/{tweet_id}/{user_id} [delete]
func (h *TweetHandler) RemoveSave(c *gin.Context) {
	var req proto.RemoveSaveRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	resp, err := h.TweetService.RemoveSave(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
