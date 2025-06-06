package handlers

import (
	"auth-service/internal/services"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type MinioHandler struct {
	MinioService *services.MinioService
}

func NewMinioHandler(minioService *services.MinioService) *MinioHandler {
	return &MinioHandler{
		MinioService: minioService,
	}
}

func (h *MinioHandler) StoreImage(c *gin.Context) {
	id := c.Param("id")

	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not provided"})
		return
	}

	src, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	objectName := fmt.Sprintf("%s-%s", id, header.Filename)

	info, err := h.MinioService.MinioClient.PutObject(
		c.Request.Context(),
		"CC-Images",
		objectName,
		src,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "upload successful",
		"object_name": info.Key,
		"size":        info.Size,
	})
}

func (h *MinioHandler) GetImageByID(c *gin.Context) {
	id := c.Param("id")
	object, err := h.MinioService.MinioClient.GetObject(
		c.Request.Context(),
		"CC-Images",
		id,
		minio.GetObjectOptions{},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get image"})
		return
	}
	stat, err := object.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))
	io.Copy(c.Writer, object)
}
