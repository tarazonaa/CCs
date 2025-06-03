package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Consumer struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	CustomID  string    `json:"custom_id" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Consumer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}


type OAuth2Application struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	ClientID     string    `json:"client_id" gorm:"uniqueIndex;not null"`
	ClientSecret string    `json:"client_secret" gorm:"not null"`
	RedirectURIs []string  `json:"redirect_uris" gorm:"serializer:json;not null"`
	HashSecret   bool      `json:"hash_secret" gorm:"default:false"`
	ConsumerID   string    `json:"consumer_id" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Consumer     Consumer  `json:"consumer,omitempty" gorm:"foreignKey:ConsumerID"`
}

func (app *OAuth2Application) BeforeCreate(tx *gorm.DB) error {
	if app.ID == "" {
		app.ID = uuid.New().String()
	}
	if app.ClientID == "" {
		app.ClientID = uuid.New().String()
	}
	if app.ClientSecret == "" {
		app.ClientSecret = uuid.New().String()
	}
	return nil
}