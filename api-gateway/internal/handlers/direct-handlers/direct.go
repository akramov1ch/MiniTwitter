package direct_handlers

import (
	proto "api-gateway/protos/direct-proto"
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

// DirectHandler goda direct message service uchun handler
type DirectHandler struct {
	DirectService  proto.DirectServiceClient
	minioClient    *minio.Client
	bucketName     string
	minioEndpoint  string
}

// NewDirectHandler yangi DirectHandler yaratadi
func NewDirectHandler(directService proto.DirectServiceClient, minioEndpoint, accessKeyID, secretAccessKey, bucketName string) (*DirectHandler, error) {
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return &DirectHandler{
		DirectService:  directService,
		minioClient:    minioClient,
		bucketName:     bucketName,
		minioEndpoint:  minioEndpoint,
	}, nil
}

// CreateDirectMessage - yangi to'g'ridan-to'g'ri xabar yaratish
// @Summary Create a direct message
// @Description Create a new direct message between users
// @Tags direct_messages
// @Accept json
// @Produce json
// @Param sender_id path int true "Sender ID"
// @Param receiver_id path int true "Receiver ID"
// @Param text formData string true "Message text"
// @Param media formData file false "Media files"
// @Success 200 {object} proto.CreateDirectMessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /directs/{sender_id}/{receiver_id} [post]
func (h *DirectHandler) CreateDirectMessage(c *gin.Context) {
	senderIdStr := c.Param("sender_id")
	receiverIdStr := c.Param("receiver_id")
	text := c.PostForm("text")

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

	senderId, err := strconv.ParseInt(senderIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid sender ID"})
		return
	}

	receiverId, err := strconv.ParseInt(receiverIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid receiver ID"})
		return
	}

	req := &proto.CreateDirectMessageRequest{
		SenderId:   senderId,
		ReceiverId: receiverId,
		Text:       text,
		Media:      media,
	}

	resp, err := h.DirectService.CreateDirectMessage(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDirectMessages - to'g'ridan-to'g'ri xabarlarni olish
// @Summary Get direct messages
// @Description Get all direct messages between two users
// @Tags direct_messages
// @Accept json
// @Produce json
// @Param sender_id path int true "Sender ID"
// @Param receiver_id path int true "Receiver ID"
// @Success 200 {object} proto.GetDirectMessagesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /directs/{sender_id}/{receiver_id} [get]
func (h *DirectHandler) GetDirectMessages(c *gin.Context) {
	senderIdStr := c.Param("sender_id")
	receiverIdStr := c.Param("receiver_id")

	senderId, err := strconv.ParseInt(senderIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid sender ID"})
		return
	}

	receiverId, err := strconv.ParseInt(receiverIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid receiver ID"})
		return
	}

	req := &proto.GetDirectMessagesRequest{
		SenderId:   senderId,
		ReceiverId: receiverId,
	}

	resp, err := h.DirectService.GetDirectMessages(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i := range resp.DirectMessages {
		for j := range resp.DirectMessages[i].Media {
			mediaName := resp.DirectMessages[i].Media[j]
			resp.DirectMessages[i].Media[j] = "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + mediaName
		}
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteDirectMessage - to'g'ridan-to'g'ri xabarni o'chirish
// @Summary Delete a direct message
// @Description Delete a specific direct message by ID
// @Tags direct_messages
// @Accept json
// @Produce json
// @Param id path int true "Direct Message ID"
// @Success 200 {object} proto.DeleteDirectMessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /directs/{id} [delete]
func (h *DirectHandler) DeleteDirectMessage(c *gin.Context) {
	var req proto.DeleteDirectMessageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	directMessage, err := h.DirectService.GetDirectMessageByID(c.Request.Context(), &proto.GetDirectMessageByIDRequest{Id: req.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get direct message"})
		return
	}

	for _, media := range directMessage.DirectMessage.Media {
		if err := h.minioClient.RemoveObject(c.Request.Context(), h.bucketName, media, minio.RemoveObjectOptions{}); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete media from MinIO"})
			return
		}
	}

	resp, err := h.DirectService.DeleteDirectMessage(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDirectMessageByID - to'g'ridan-to'g'ri xabarni ID bo'yicha olish
// @Summary Get a direct message by ID
// @Description Get a specific direct message by ID
// @Tags direct_messages
// @Accept json
// @Produce json
// @Param id path int true "Direct Message ID"
// @Success 200 {object} proto.GetDirectMessageByIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /directs/{id} [get]
func (h *DirectHandler) GetDirectMessageByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid direct message ID"})
		return
	}

	req := &proto.GetDirectMessageByIDRequest{
		Id: id,
	}

	resp, err := h.DirectService.GetDirectMessageByID(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Media URL'larini yaratish
	for j := range resp.DirectMessage.Media {
		mediaName := resp.DirectMessage.Media[j]
		resp.DirectMessage.Media[j] = "https://" + h.minioEndpoint + "/" + h.bucketName + "/" + mediaName
	}

	c.JSON(http.StatusOK, resp)
}
