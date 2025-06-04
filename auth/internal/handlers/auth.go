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
