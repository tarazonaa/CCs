package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID          uuid.UUID `json:"user_id" gorm:"not null;type:uuid"`
	SentImageID     uuid.UUID `json:"sent_image_id" gorm:"uniqueIndex;not null;type:uuid"`
	ReceivedImageID uuid.UUID `json:"received_image_id" gorm:"uniqueIndex;not null;type:uuid"`
	CreatedAt       time.Time `json:"created_at"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (i *Image) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
