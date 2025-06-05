// En tu archivo models/oauth2_token.go
package models

import (
	"auth-service/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OAuth2Token struct {
	ID                  string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AccessToken         string `json:"access_token" gorm:"uniqueIndex;not null"`
	RefreshToken        string `json:"refresh_token,omitempty" gorm:"uniqueIndex"`
	AccessTokenExpiration time.Time    `json:"access_token_expiration,omitempty"` 
	RefreshTokenExpiration time.Time    `json:"refresh_token_expiration,omitempty"`
	Scope               string `json:"scope,omitempty"`
	AuthenticatedUserID string `json:"authenticated_userid,omitempty" gorm:"column:authenticated_userid"`
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
	if t.CreatedAt == 0 {
		t.CreatedAt = utils.GetCurrentTS().Unix() * 1000
	}
	return nil
}

func (t *OAuth2Token) IsExpired() bool {
	return utils.GetCurrentTS().After(t.AccessTokenExpiration) 
}

func generateRandomToken() string {
	return uuid.New().String() + uuid.New().String()
}

func (t *OAuth2Token) IsRefreshable() bool {
	// Check if the access token is close to expiration (e.g. within the next hour)
	if t.IsExpired() {
		return utils.GetCurrentTS().Before(t.RefreshTokenExpiration)
	}

	if utils.GetCurrentTS().Before(t.AccessTokenExpiration.Add(-time.Hour)) {
		return false
	}
	
	return true
}
