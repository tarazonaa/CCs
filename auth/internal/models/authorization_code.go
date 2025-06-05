package models

import (
	"auth-service/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AuthorizationCode struct {
	ID                  string         `json:"id" gorm:"primaryKey"`
	Code                string         `json:"code" gorm:"uniqueIndex;not null"`
	ClientID            string         `json:"client_id" gorm:"not null"`
	UserID              uuid.UUID      `json:"user_id" gorm:"not null"`
	RedirectURI         string         `json:"redirect_uri" gorm:"not null"`
	Scopes              pq.StringArray `json:"scopes" gorm:"type:text[]"`
	CodeChallenge       string         `json:"-"`
	CodeChallengeMethod string         `json:"-"`
	ExpiresAt           time.Time      `json:"expires_at"`
	IsUsed              bool           `json:"is_used" gorm:"default:false"`
	CreatedAt           time.Time      `json:"created_at"`

	// Relations
	Client OAuth2Credential `json:"client,omitempty" gorm:"foreignKey:ClientID;references:ClientID"`
	User   User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (ac *AuthorizationCode) BeforeCreate(tx *gorm.DB) error {
	if ac.ID == "" {
		ac.ID = uuid.New().String()
	}
	if ac.Code == "" {
		ac.Code = uuid.New().String()
	}
	if ac.ExpiresAt.IsZero() {
		ac.ExpiresAt = utils.GetCurrentTS().Add(10 * time.Minute)
	}
	return nil
}

func (ac *AuthorizationCode) IsExpired() bool {
	return utils.GetCurrentTS().After(ac.ExpiresAt)
}

func (ac *AuthorizationCode) IsValid() bool {
	return !ac.IsExpired() && !ac.IsUsed
}
