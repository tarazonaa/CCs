package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/models"
	"auth-service/internal/seeds"
	"auth-service/internal/services"
	"auth-service/internal/utils"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

var globalDB *gorm.DB

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()
	db := config.InitDatabase(cfg)
	globalDB = db

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.RootUser, cfg.Minio.RootPwd, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	setupMinioClient(minioClient, &ctx)

	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	for _, bucket := range buckets {
		log.Println(bucket.Name)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Consumer{}, &models.OAuth2Token{}, &models.OAuth2Credential{}, &models.AuthorizationCode{}, &models.Image{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	if err := seeds.SeedClients(db, "clients.json"); err != nil {
		log.Fatal("Seed failed: ", err)
	}

	oauth2Service := services.NewOAuth2Service(db, cfg)
	oauth2Handler := handlers.NewOAuth2Handler(oauth2Service, db, cfg)
	authHandler := handlers.NewAuthHandler(db)
	imageService := services.NewImageService(db)
	imageHandler := handlers.NewImageHandler(imageService, minioClient)

	router := setupRouter(oauth2Handler, authHandler, imageHandler)

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

func setupMinioClient(client *minio.Client, ctx *context.Context) {
	err := client.MakeBucket(*ctx, "cc-images", minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(*ctx, "cc-images")
		if errBucketExists == nil && exists {
			log.Printf("Bucket exists.")
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created bucket")
	}
}

func setupRouter(oauth2Handler *handlers.OAuth2Handler, authHandler *handlers.AuthHandler, imageHandler *handlers.ImageHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": utils.GetCurrentTS(),
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
		oauth2Group.POST("/introspect", oauth2Handler.IntrospectToken)
	}

	authGroup := router.Group("/auth")
	{
		authGroup.GET("/authorize", authHandler.ShowAuthorizationPage)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/logout", authHandler.Logout)
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

		// Image routes
		imageGroup := apiGroup.Group("/images")
		{
			imageGroup.POST("", imageHandler.CreateImage)
			imageGroup.GET("", imageHandler.GetAllImages)
			imageGroup.GET("/:id", imageHandler.GetImageByID)
			imageGroup.DELETE("/:id", imageHandler.DeleteImage)
			imageGroup.GET("/blob/:id", imageHandler.GetBlobFromID)
			imageGroup.GET("/sent/:sent_image_id", imageHandler.GetImageBySentID)
			imageGroup.GET("/received/:received_image_id", imageHandler.GetImageByReceivedID)
		}

		apiGroup.GET("/users/:user_id/images", imageHandler.GetUserImages)
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
		c.Header("Access-Control-Allow-Origin", "*")
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
		Name         string    `json:"name" binding:"required"`
		ClientID     string    `json:"client_id"`
		ClientSecret string    `json:"client_secret"`
		RedirectURIs []string  `json:"redirect_uris" binding:"required"`
		ConsumerID   uuid.UUID `json:"consumer_id" binding:"required"`
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
