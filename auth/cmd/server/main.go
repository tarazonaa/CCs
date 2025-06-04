package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/models"
	"auth-service/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var globalDB *gorm.DB

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()
	db := config.InitDatabase(cfg)
	globalDB = db

	oauth2Service := services.NewOAuth2Service(db, cfg)
	oauth2Handler := handlers.NewOAuth2Handler(oauth2Service, db, cfg)
	authHandler := handlers.NewAuthHandler(db)

	router := setupRouter(oauth2Handler, authHandler)

	server := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("OAuth 2.0 Authorization Server starting on %s:%s", cfg.Host, cfg.Port)
	log.Printf("Provision Key: %s", cfg.ProvisionKey)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(oauth2Handler *handlers.OAuth2Handler, authHandler *handlers.AuthHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "oauth2-authorization-server",
		})
	})

	oauth2Group := router.Group("/oauth2")
	{
		oauth2Group.GET("/authorize", oauth2Handler.OAuth2Authorize)
		oauth2Group.POST("/authorize", oauth2Handler.OAuth2Authorize)
		oauth2Group.POST("/token", oauth2Handler.OAuth2Token)
		oauth2Group.Any("/tokens", oauth2Handler.OAuth2Tokens)
		oauth2Group.Any("/tokens/:token_id", oauth2Handler.OAuth2TokenByID)
	}

	authGroup := router.Group("/auth")
	{
		authGroup.GET("/authorize", authHandler.ShowAuthorizationPage)
	}

	apiGroup := router.Group("/api/v1")
	apiGroup.Use(oauth2Handler.ValidateToken())
	{
		apiGroup.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"user_id":   c.GetString("authenticated_userid"),
				"client_id": c.GetString("client_id"),
				"scope":     c.GetString("scope"),
				"message":   "This is a protected resource",
			})
		})
	}

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/clients", listClients)
		adminGroup.POST("/clients", createClient)
		adminGroup.GET("/clients/:client_id", getClient)
		adminGroup.PUT("/clients/:client_id", updateClient)
		adminGroup.DELETE("/clients/:client_id", deleteClient)
		adminGroup.POST("/consumers", createConsumer)
		adminGroup.GET("/consumers", listConsumers)
		adminGroup.GET("/consumers/:consumer_id", getConsumer)
	}

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://127.0.0.1:3000": true,
			"http://localhost:3001": true,
		}

		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func createConsumer(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		CustomID string `json:"custom_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	consumer := &models.Consumer{
		Username: req.Username,
		CustomID: req.CustomID,
	}

	if err := globalDB.Create(consumer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create consumer"})
		return
	}

	c.JSON(http.StatusCreated, consumer)
}

func listConsumers(c *gin.Context) {
	var consumers []models.Consumer
	if err := globalDB.Find(&consumers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch consumers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(consumers),
		"data":  consumers,
	})
}

func getConsumer(c *gin.Context) {
	consumerID := c.Param("consumer_id")

	var consumer models.Consumer
	if err := globalDB.Where("id = ?", consumerID).First(&consumer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "consumer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch consumer"})
		return
	}

	c.JSON(http.StatusOK, consumer)
}

func listClients(c *gin.Context) {
	var clients []models.OAuth2Credential
	if err := globalDB.Preload("Consumer").Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch clients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(clients),
		"data":  clients,
	})
}

func createClient(c *gin.Context) {
	var req struct {
		Name         string   `json:"name" binding:"required"`
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris" binding:"required"`
		ConsumerID   string   `json:"consumer_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	client := &models.OAuth2Credential{
		Name:         req.Name,
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURIs: req.RedirectURIs,
		ConsumerID:   req.ConsumerID,
	}

	if err := globalDB.Create(client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create client"})
		return
	}

	c.JSON(http.StatusCreated, client)
}

func getClient(c *gin.Context) {
	clientID := c.Param("client_id")

	var client models.OAuth2Credential
	if err := globalDB.Preload("Consumer").Where("client_id = ?", clientID).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch client"})
		return
	}

	c.JSON(http.StatusOK, client)
}

func updateClient(c *gin.Context) {
	clientID := c.Param("client_id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var client models.OAuth2Credential
	if err := globalDB.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch client"})
		return
	}

	if err := globalDB.Model(&client).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update client"})
		return
	}

	c.JSON(http.StatusOK, client)
}

func deleteClient(c *gin.Context) {
	clientID := c.Param("client_id")

	result := globalDB.Where("client_id = ?", clientID).Delete(&models.OAuth2Credential{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete client"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
