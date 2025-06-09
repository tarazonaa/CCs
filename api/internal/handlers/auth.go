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

// ShowAuthorizationPage godoc
// @Summary      Show OAuth2 authorization page
// @Description  Returns the OAuth2 authorization page or JSON for a given client_id
// @Tags         auth
// @Accept       json
// @Produce      json,html
// @Param        client_id  query     string  true  "Client ID"
// @Param        scope      query     string  false "OAuth2 scope"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /auth/authorize [get]
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

// Register godoc
// @Summary      Register a new user
// @Description  Creates a user account with email, password, username, and name
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body   object{email=string,password=string,username=string,name=string}  true  "User registration object"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Username string `json:"username" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var existing models.User
	if err := h.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Name:     req.Name,
		Password: req.Password,
	}

	if err := h.db.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Revokes the access token and logs out the user
// @Tags         auth
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get access token from the request
	accessToken := c.GetHeader("Authorization")
	// Remove "Bearer " prefix if present
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing access token"})
		return
	}
	// Revoke all tokens for the user 
	if err := h.db.Where("access_token = ?", accessToken).Delete(&models.OAuth2Token{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke tokens"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
