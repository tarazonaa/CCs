package services

import (
	"auth-service/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageService struct {
	db *gorm.DB
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{
		db: db,
	}
}

type CreateImageRequest struct {
	SentImageID     *uuid.UUID `json:"sent_image_id,omitempty"`
	ReceivedImageID *uuid.UUID `json:"received_image_id,omitempty"`
}

type ImageResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	SentImageID     uuid.UUID `json:"sent_image_id"`
	ReceivedImageID uuid.UUID `json:"received_image_id"`
	CreatedAt       string    `json:"created_at"`
}

func (s *ImageService) CreateImage(userID uuid.UUID, req *CreateImageRequest) (*ImageResponse, error) {
	// Validate user exists
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to validate user")
	}

	// Create image record
	image := &models.Image{
		UserID: userID,
	}

	// Use provided IDs or generate new ones
	if req.SentImageID != nil {
		image.SentImageID = *req.SentImageID
	}
	if req.ReceivedImageID != nil {
		image.ReceivedImageID = *req.ReceivedImageID
	}

	if err := s.db.Create(image).Error; err != nil {
		return nil, errors.New("failed to create image record")
	}

	response := &ImageResponse{
		ID:              image.ID,
		UserID:          image.UserID,
		SentImageID:     image.SentImageID,
		ReceivedImageID: image.ReceivedImageID,
		CreatedAt:       image.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return response, nil
}

func (s *ImageService) GetAllImages(userID *uuid.UUID, limit, offset int) ([]ImageResponse, int64, error) {
	var images []models.Image
	var total int64

	query := s.db.Model(&models.Image{}).Preload("User")
	
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count images")
	}

	// paginated
	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&images).Error; err != nil {
		return nil, 0, errors.New("failed to fetch images")
	}

	responses := make([]ImageResponse, len(images))
	for i, img := range images {
		responses[i] = ImageResponse{
			ID:              img.ID,
			UserID:          img.UserID,
			SentImageID:     img.SentImageID,
			ReceivedImageID: img.ReceivedImageID,
			CreatedAt:       img.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return responses, total, nil
}

func (s *ImageService) GetImageByID(imageID uuid.UUID) (*ImageResponse, error) {
	var image models.Image
	if err := s.db.Preload("User").Where("id = ?", imageID).First(&image).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("image not found")
		}
		return nil, errors.New("failed to fetch image")
	}

	response := &ImageResponse{
		ID:              image.ID,
		UserID:          image.UserID,
		SentImageID:     image.SentImageID,
		ReceivedImageID: image.ReceivedImageID,
		CreatedAt:       image.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return response, nil
}

func (s *ImageService) GetImageBySentID(sentImageID uuid.UUID) (*ImageResponse, error) {
	var image models.Image
	if err := s.db.Preload("User").Where("sent_image_id = ?", sentImageID).First(&image).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("image not found")
		}
		return nil, errors.New("failed to fetch image")
	}

	return &ImageResponse{
		ID:              image.ID,
		UserID:          image.UserID,
		SentImageID:     image.SentImageID,
		ReceivedImageID: image.ReceivedImageID,
		CreatedAt:       image.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *ImageService) GetImageByReceivedID(receivedImageID uuid.UUID) (*ImageResponse, error) {
	var image models.Image
	if err := s.db.Preload("User").Where("received_image_id = ?", receivedImageID).First(&image).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("image not found")
		}
		return nil, errors.New("failed to fetch image")
	}

	return &ImageResponse{
		ID:              image.ID,
		UserID:          image.UserID,
		SentImageID:     image.SentImageID,
		ReceivedImageID: image.ReceivedImageID,
		CreatedAt:       image.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}


func (s *ImageService) DeleteImage(imageID uuid.UUID, userID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", imageID, userID).Delete(&models.Image{})
	if result.Error != nil {
		return errors.New("failed to delete image")
	}

	if result.RowsAffected == 0 {
		return errors.New("image not found")
	}

	return nil
}
