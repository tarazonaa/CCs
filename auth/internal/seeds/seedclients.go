package seeds

import (
	"encoding/json"
	"fmt"
	"os"

	"auth-service/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RawClient struct {
	Name         string `json:"name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func SeedOAuthClients(db *gorm.DB, path string) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var rawClients []RawClient
	if err := json.Unmarshal(data, &rawClients); err != nil {
		return fmt.Errorf("Invalid json: %w", err)
	}

	for _, raw := range rawClients {
		var existing models.OAuth2Credential
		if err := db.Where("client_id = ?", raw.ClientID).First(&existing).Error; err != nil {
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(raw.ClientSecret), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("Failed to hash secret: %w", err)
		}

		var consumer models.Consumer
		if err := db.FirstOrCreate(&consumer, models.Consumer{
			Username: "ccs-global-consumer",
			CustomID: "ccs-global-id",
		}).Error; err != nil {
			return fmt.Errorf("Missing consumer to associate client: %w", err)
		}

		client := models.OAuth2Credential{
			Name:         raw.Name,
			ClientID:     raw.ClientID,
			ClientSecret: string(hashed),
			ConsumerID:   consumer.ID,
		}

		if err := db.Create(&client).Error; err != nil {
			return fmt.Errorf("Failed to create client: %w", err)
		}
	}
	return nil
}
