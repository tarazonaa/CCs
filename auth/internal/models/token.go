package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OAuth2Token struct {
	ID                     string    `json:"id" gorm:"primaryKey"`
	AccessToken            string    `json:"access_token" gorm:"uniqueIndex;not null"`
	RefreshToken           string    `json:"refresh_token,omitempty" gorm:"uniqueIndex"`
	TokenType              string    `json:"token_type" gorm:"default:bearer"`
	ExpiresIn              int       `json:"expires_in"` 
	Scope                  string    `json:"scope,omitempty"`
	AuthenticatedUserID    string    `json:"authenticated_userid,omitempty"`
	CredentialID           string    `json:"credential_id" gorm:"not null"` 
	ServiceID              string    `json:"service_id,omitempty"`
	CreatedAt              int64     `json:"created_at"`


	Credential OAuth2Application `json:"credential,omitempty" gorm:"foreignKey:CredentialID"`
}

func (t *OAuth2Token) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	if t.AccessToken == "" {
		t.AccessToken = generateRandomToken()
	}
	if t.TokenType == "" {
		t.TokenType = "bearer"
	}
	if t.CreatedAt == 0 {
		t.CreatedAt = time.Now().Unix() * 1000 // Kong uses milliseconds
	}
	return nil
}

func (t *OAuth2Token) IsExpired() bool {
	if t.ExpiresIn == 0 {
		return false 
	}
	expirationTime := time.Unix(t.CreatedAt/1000, 0).Add(time.Duration(t.ExpiresIn) * time.Second)
	return time.Now().After(expirationTime)
}

func generateRandomToken() string {
	return uuid.New().String() + uuid.New().String()
}