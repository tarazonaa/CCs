package handlers

import (
	"auth-service/internal/services"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type ImageHandler struct {
	imageService *services.ImageService
	MinioClient  *minio.Client
}

func NewImageHandler(imageService *services.ImageService, minioClient *minio.Client) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
		MinioClient:  minioClient,
	}
}

// CreateImage godoc
// @Summary      Create an image record and upload images
// @Description  Uploads original and inference images, creates a record for the user
// @Tags         images
// @Accept       multipart/form-data
// @Produce      json
// @Param        original_image   formData  file   true  "Original image file"
// @Param        inference_image  formData  file   true  "Inference image file"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /images [post]
func (h *ImageHandler) CreateImage(c *gin.Context) {
	userIDStr := c.GetString("authenticated_userid")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID in token"})
		return
	}

	inferenceID := uuid.New()
	receivedID := uuid.New()

	originalHeader, err := c.FormFile("original_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "original_image is required"})
		return
	}
	originalFile, _ := originalHeader.Open()
	defer originalFile.Close()

	inferenceHeader, err := c.FormFile("inference_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "inference_image is required"})
		return
	}
	inferenceFile, _ := inferenceHeader.Open()
	defer inferenceFile.Close()

	_, err = h.MinioClient.PutObject(
		c.Request.Context(),
		"cc-images",
		fmt.Sprintf("%s.png", receivedID),
		inferenceFile,
		inferenceHeader.Size,
		minio.PutObjectOptions{ContentType: inferenceHeader.Header.Get("Content-Type")},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload inference image"})
		return
	}

	_, err = h.MinioClient.PutObject(
		c.Request.Context(),
		"cc-images",
		fmt.Sprintf("%s.png", inferenceID),
		originalFile,
		originalHeader.Size,
		minio.PutObjectOptions{ContentType: originalHeader.Header.Get("Content-Type")},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload original image"})
		return
	}

	image, err := h.imageService.CreateImage(userID, &services.CreateImageRequest{
		SentImageID:     &inferenceID,
		ReceivedImageID: &receivedID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create image",
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":           "Record and images created",
		"id":                image.ID,
		"sent_image_id":     inferenceID,
		"received_image_id": receivedID,
	})
}

// GetBlobFromID godoc
// @Summary      Get image blob by ID
// @Description  Returns the image file as binary
// @Tags         images
// @Produce      image/png
// @Param        id  path  string  true  "Image ID"
// @Success      200  {file}  file
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /images/blob/{id} [get]
func (h *ImageHandler) GetBlobFromID(c *gin.Context) {
	imageID := c.Param("id")
	objectName := fmt.Sprintf("%s.png", imageID)
	obj, err := h.MinioClient.GetObject(
		c.Request.Context(),
		"cc-images",
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get image"})
		return
	}

	stat, err := obj.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))
	io.Copy(c.Writer, obj)
}

// GetAllImages godoc
// @Summary      List all images
// @Description  Lists all images, optionally filtered by user_id, with pagination
// @Tags         images
// @Produce      json
// @Param        user_id  query  string  false  "User ID"
// @Param        limit    query  int     false  "Limit"
// @Param        offset   query  int     false  "Offset"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /images [get]
func (h *ImageHandler) GetAllImages(c *gin.Context) {
	// Parse query parameters
	var userID *uuid.UUID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		parsedUserID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
		userID = &parsedUserID
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			if parsedLimit > 100 {
				parsedLimit = 100
			}
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	images, total, err := h.imageService.GetAllImages(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   images,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetImageByID godoc
// @Summary      Get image by UUID
// @Description  Returns the image metadata for the given image ID
// @Tags         images
// @Produce      json
// @Param        id  path  string  true  "Image ID"
// @Success      200  {object}  object
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /images/{id} [get]
func (h *ImageHandler) GetImageByID(c *gin.Context) {
	imageIDStr := c.Param("id")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID format"})
		return
	}

	image, err := h.imageService.GetImageByID(imageID)
	if err != nil {
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

// GetImageBySentID godoc
// @Summary      Get image by sent image ID
// @Description  Returns image metadata by sent image UUID
// @Tags         images
// @Produce      json
// @Param        sent_image_id  path  string  true  "Sent Image ID"
// @Success      200  {object}  object
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /images/sent/{sent_image_id} [get]
func (h *ImageHandler) GetImageBySentID(c *gin.Context) {
	sentImageIDStr := c.Param("sent_image_id")
	sentImageID, err := uuid.Parse(sentImageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sent_image_id format"})
		return
	}

	image, err := h.imageService.GetImageBySentID(sentImageID)
	if err != nil {
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

// GetImageByReceivedID godoc
// @Summary      Get image by received image ID
// @Description  Returns image metadata by received image UUID
// @Tags         images
// @Produce      json
// @Param        received_image_id  path  string  true  "Received Image ID"
// @Success      200  {object}  object
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /images/received/{received_image_id} [get]
func (h *ImageHandler) GetImageByReceivedID(c *gin.Context) {
	receivedImageIDStr := c.Param("received_image_id")
	receivedImageID, err := uuid.Parse(receivedImageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid received_image_id format"})
		return
	}

	image, err := h.imageService.GetImageByReceivedID(receivedImageID)
	if err != nil {
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

// GetUserImages godoc
// @Summary      Get all images for authenticated user
// @Description  Returns paginated images for the current user
// @Tags         images
// @Produce      json
// @Param        limit   query  int  false  "Limit"
// @Param        offset  query  int  false  "Offset"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /images/user [get]
func (h *ImageHandler) GetUserImages(c *gin.Context) {
	userIDStr, exists := c.Get("authenticated_userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	userIDParsed, ok := userIDStr.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("User id found: %s", userIDParsed)})
		return
	}

	userID, err := uuid.Parse(userIDParsed)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Couldn't parse uuid"})
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			if parsedLimit > 100 {
				parsedLimit = 100
			}
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	images, total, err := h.imageService.GetAllImages(&userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    images,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"user_id": userIDParsed,
	})
}

// DeleteImage godoc
// @Summary      Delete an image
// @Description  Deletes image by ID for the authenticated user
// @Tags         images
// @Param        id  path  string  true  "Image ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /images/{id} [delete]
func (h *ImageHandler) DeleteImage(c *gin.Context) {

	userIDStr := c.GetString("authenticated_userid")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID in token"})
		return
	}

	imageIDStr := c.Param("id")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID format"})
		return
	}

	err = h.imageService.DeleteImage(imageID, userID)
	if err != nil {
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
