	package config

	import (
		"os"
		"strconv"
	)

	type Config struct {
		Port string
		Host string
		DatabaseURL string
		OAuth2 OAuth2Config
		ProvisionKey string
	}

	type OAuth2Config struct {
		AccessTokenExpiration  int  `json:"token_expiration"`
		RefreshTokenExpiration int  `json:"refresh_token_expiration"`
		AuthCodeExpiration     int  `json:"auth_code_expiration"`
		
		EnableClientCredentials      bool `json:"enable_client_credentials"`
		EnableAuthorizationCode      bool `json:"enable_authorization_code"`
		EnableImplicitGrant         bool `json:"enable_implicit_grant"`
		EnablePasswordCredentials   bool `json:"enable_password_credentials"`
		
		EnablePKCE         bool `json:"enable_pkce"`
		PKCERequired       bool `json:"pkce_required"`
		
		// Token settings
		ReuseRefreshToken  bool `json:"reuse_refresh_token"`
		AcceptHTTPIfAlreadyTerminated bool `json:"accept_http_if_already_terminated"`
		
		// Global credentials
		GlobalCredentials  bool `json:"global_credentials"`
		Anonymous         string `json:"anonymous"`
		HideCredentials   bool `json:"hide_credentials"`
	}


	func LoadConfig() *Config {
		return &Config{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
			DatabaseURL: getEnv("DATABASE_URL", "oauth2.db"),
			ProvisionKey: getEnv("PROVISION_KEY", generateProvisionKey()),
			
			OAuth2: OAuth2Config{
				AccessTokenExpiration:  getEnvAsInt("ACCESS_TOKEN_EXPIRATION", 7200),   
				RefreshTokenExpiration: getEnvAsInt("REFRESH_TOKEN_EXPIRATION", 1209600),
				AuthCodeExpiration:     getEnvAsInt("AUTH_CODE_EXPIRATION", 600),
				
				EnableClientCredentials:    getEnvAsBool("ENABLE_CLIENT_CREDENTIALS", true),
				EnableAuthorizationCode:    getEnvAsBool("ENABLE_AUTHORIZATION_CODE", true),
				EnableImplicitGrant:       getEnvAsBool("ENABLE_IMPLICIT_GRANT", false),
				EnablePasswordCredentials: getEnvAsBool("ENABLE_PASSWORD_CREDENTIALS", false),
				
				EnablePKCE:               getEnvAsBool("ENABLE_PKCE", true),
				PKCERequired:             getEnvAsBool("PKCE_REQUIRED", false),
				
				ReuseRefreshToken:        getEnvAsBool("REUSE_REFRESH_TOKEN", false),
				GlobalCredentials:        getEnvAsBool("GLOBAL_CREDENTIALS", false),
				HideCredentials:          getEnvAsBool("HIDE_CREDENTIALS", false),
				
				Anonymous: getEnv("ANONYMOUS", ""),
			},
		}
	}

	// Helper functions
	func getEnv(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}

	func getEnvAsInt(key string, defaultValue int) int {
		if value := os.Getenv(key); value != "" {
			if intValue, err := strconv.Atoi(value); err == nil {
				return intValue
			}
		}
		return defaultValue
	}

	func getEnvAsBool(key string, defaultValue bool) bool {
		if value := os.Getenv(key); value != "" {
			if boolValue, err := strconv.ParseBool(value); err == nil {
				return boolValue
			}
		}
		return defaultValue
	}

	func generateProvisionKey() string {
		
		return "default-provision-key-change-in-production"
	}