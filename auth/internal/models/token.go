// En tu archivo models/oauth2_token.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OAuth2Token struct {
	ID                  string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AccessToken         string `json:"access_token" gorm:"uniqueIndex;not null"`
	RefreshToken        string `json:"refresh_token,omitempty" gorm:"uniqueIndex"`
	TokenType           string `json:"token_type" gorm:"default:bearer"`
	ExpiresIn           int    `json:"expires_in"`
	Scope               string `json:"scope,omitempty"`
	AuthenticatedUserID string `json:"authenticated_userid,omitempty" gorm:"column:authenticated_userid"` // ← CAMBIO AQUÍ
	CredentialID        string `json:"credential_id" gorm:"not null;type:uuid"`
	CreatedAt           int64  `json:"created_at"`

	// Relación
	Credential OAuth2Credential `json:"credential,omitempty" gorm:"foreignKey:CredentialID"`
}

func (OAuth2Token) TableName() string {
	return "oauth2_tokens"
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
	if t.TokenType == "refresh" && t.RefreshToken == "" {
		t.RefreshToken = generateRandomToken()
	}
	if t.CreatedAt == 0 {
		t.CreatedAt = time.Now().Unix() * 1000
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
