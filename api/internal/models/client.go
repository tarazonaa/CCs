package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Consumer struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	CustomID  string    `json:"custom_id" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Consumer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type OAuth2Credential struct {
	ID           uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string         `json:"name" gorm:"not null"`
	ClientID     string         `json:"client_id" gorm:"uniqueIndex;not null"`
	ClientSecret string         `json:"client_secret" gorm:"not null"`
	RedirectURIs pq.StringArray `json:"redirect_uris" gorm:"type:text[]" swaggertype:"array,string"`
	ConsumerID   uuid.UUID      `json:"consumer_id" gorm:"not null;type:uuid"`
	CreatedAt    time.Time      `json:"created_at"`

	// Relación
	Consumer Consumer `json:"consumer,omitempty" gorm:"foreignKey:ConsumerID;constraint:OnDelete:CASCADE"`
}

func (OAuth2Credential) TableName() string {
	return "oauth2_credentials"
}

func (app *OAuth2Credential) BeforeCreate(tx *gorm.DB) error {
	if app.ID == uuid.Nil {
		app.ID = uuid.New()
	}
	if app.ClientID == "" {
		app.ClientID = uuid.New().String()
	}
	if app.ClientSecret == "" {
		raw := uuid.New().String()
		hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		app.ClientSecret = string(hashed)
	}
	return nil
}
