// internal/config/database.go
package config

import (
	"log"
	"strings"

	// "auth-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(config *Config) *gorm.DB {
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required. Please set it to a PostgreSQL connection string.")
	}

	if !strings.HasPrefix(config.DatabaseURL, "postgres://") && !strings.HasPrefix(config.DatabaseURL, "postgresql://") {
		log.Fatal("Only PostgreSQL databases are supported. DATABASE_URL must start with 'postgres://' or 'postgresql://'")
	}

	// Connect to PostgreSQL
	log.Printf("Attempting to connect to PostgreSQL: %s", maskPassword(config.DatabaseURL))
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
	}
	log.Println("Connected to PostgreSQL database successfully")

	// err = db.AutoMigrate(
	// 	&models.Consumer{},
	// 	&models.OAuth2Credential{},
	// 	&models.OAuth2Token{},
	// )
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}

func maskPassword(dbURL string) string {
	if strings.Contains(dbURL, "@") {
		parts := strings.Split(dbURL, "@")
		if len(parts) >= 2 {
			userPart := parts[0]
			if strings.Contains(userPart, ":") {
				userParts := strings.Split(userPart, ":")
				if len(userParts) >= 3 {
					// postgres://user:password@host -> postgres://user:***@host
					userParts[2] = "***"
					return strings.Join(userParts, ":") + "@" + strings.Join(parts[1:], "@")
				}
			}
		}
	}
	return dbURL
}
