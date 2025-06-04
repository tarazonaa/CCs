package handlers

import (
	"auth-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) ShowAuthorizationPage(c *gin.Context) {
	clientID := c.Query("client_id")
	scope := c.Query("scope")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id required"})
		return
	}

	// Get client details
	var app models.OAuth2Credential
	if err := h.db.Preload("Consumer").Where("client_id = ?", clientID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	// Return HTML page or JSON (depending on Accept header)
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, gin.H{
			"client_id":   app.ClientID,
			"client_name": app.Name,
			"scope":       scope,
			"consumer":    app.Consumer,
		})
	} else {
		// In production, render HTML template
		c.HTML(http.StatusOK, "authorize.html", gin.H{
			"ClientID":   app.ClientID,
			"ClientName": app.Name,
			"Scope":      scope,
			"Consumer":   app.Consumer,
		})
	}
}


func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email	 string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var existing models.User
	if err := h.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
	}

	user := &models.User{
		Email: req.Email,
		Username: req.Username,
	}

	// Should add error handling?
	user.HashPassword(req.Password)

	if err := h.db.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": user.ID,
		"email": user.Email,
		"username": user.Username,
	})
}
