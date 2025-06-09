package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/models"
	"auth-service/internal/services"
	"auth-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OAuth2Handler struct {
	oauth2Service *services.OAuth2Service
	db            *gorm.DB
	config        *config.Config
}

func NewOAuth2Handler(oauth2Service *services.OAuth2Service, db *gorm.DB, cfg *config.Config) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2Service: oauth2Service,
		db:            db,
		config:        cfg,
	}
}

// IntrospectToken godoc
// @Summary      Introspect token
// @Description  Checks if a token is valid and returns token/user info
// @Tags         oauth2
// @Accept       json
// @Produce      json
// @Param        token  body  object{token=string}  true  "Token to introspect"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /oauth2/introspect [post]
func (h *OAuth2Handler) IntrospectToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding: "required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"active": false, "error": "missing token"})
		return
	}

	var token models.OAuth2Token
	if err := h.db.Preload("Credential").Where("access_token = ?", req.Token).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"active": false})
		return
	}

	if token.IsExpired() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"active":  false,
			"refresh": true,
		})
		return
	}

	if token.IsRefreshable() {
		c.JSON(http.StatusOK, gin.H{
			"active":         true,
			"should_refresh": true,
			"exp":            token.AccessTokenExpiration.Unix(),
		})
		return
	}

	// Get the user from DB
	userUUID, err := uuid.Parse(token.AuthenticatedUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"active": false,
		})
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", userUUID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"active": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active":        true,
		"username":      user.Username,
		"email":         user.Email,
		"client_id":     token.Credential.ClientID,
		"refresh_token": token.RefreshToken,
		"scope":         token.Scope,
		"exp":           token.AccessTokenExpiration,
	})
}

// OAuth2Authorize godoc
// @Summary      OAuth2 authorization endpoint
// @Description  Handles OAuth2 authorization requests
// @Tags         oauth2
// @Accept       application/x-www-form-urlencoded
// @Produce      json
// @Param        response_type  formData  string  true  "Response type"
// @Param        client_id      formData  string  true  "Client ID"
// @Param        redirect_uri   formData  string  true  "Redirect URI"
// @Param        scope          formData  string  false "Scope"
// @Param        state          formData  string  false "State"
// @Success      302  {string}  string  "Redirects to client"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /oauth2/authorize [get]
func (h *OAuth2Handler) OAuth2Authorize(c *gin.Context) {
	var req services.AuthorizeRequest

	if err := c.ShouldBind(&req); err != nil {
		h.sendErrorRedirect(c, "invalid_request", "Invalid request parameters", req.RedirectURI, req.State)
		return
	}

	if req.ResponseType == "" || req.ClientID == "" {
		h.sendErrorRedirect(c, "invalid_request", "Missing required parameters", req.RedirectURI, req.State)
		return
	}

	response, err := h.oauth2Service.Authorize(&req)
	if err != nil {
		h.sendErrorRedirect(c, "invalid_client", err.Error(), req.RedirectURI, req.State)
		return
	}

	c.Redirect(http.StatusFound, response.RedirectURI)
}

// OAuth2Token godoc
// @Summary      OAuth2 token endpoint
// @Description  Issues OAuth2 tokens (access/refresh) for a client
// @Tags         oauth2
// @Accept       application/x-www-form-urlencoded,application/json
// @Produce      json
// @Param        grant_type     formData  string  true  "Grant type"
// @Param        client_id      formData  string  false "Client ID"
// @Param        client_secret  formData  string  false "Client Secret"
// @Param        code           formData  string  false "Authorization code"
// @Param        redirect_uri   formData  string  false "Redirect URI"
// @Param        refresh_token  formData  string  false "Refresh token"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /oauth2/token [post]
func (h *OAuth2Handler) OAuth2Token(c *gin.Context) {
	var req services.TokenRequest

	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&req); err != nil {
			h.sendTokenError(c, "invalid_request", "Invalid JSON request", http.StatusBadRequest)
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			h.sendTokenError(c, "invalid_request", "Invalid form request", http.StatusBadRequest)
			return
		}
	}

	if req.ClientID == "" || req.ClientSecret == "" {
		if clientID, clientSecret, ok := c.Request.BasicAuth(); ok {
			req.ClientID = clientID
			req.ClientSecret = clientSecret
		}
	}

	if req.GrantType == "" {
		h.sendTokenError(c, "invalid_request", "Missing grant_type", http.StatusBadRequest)
		return
	}

	tokenResponse, err := h.oauth2Service.Token(&req)
	if err != nil {
		status := http.StatusBadRequest
		errorCode := "invalid_request"

		switch {
		case strings.Contains(err.Error(), "invalid client"):
			errorCode = "invalid_client"
			status = http.StatusUnauthorized
		case strings.Contains(err.Error(), "invalid grant"):
			errorCode = "invalid_grant"
		case strings.Contains(err.Error(), "unsupported grant"):
			errorCode = "unsupported_grant_type"
		case strings.Contains(err.Error(), "invalid scope"):
			errorCode = "invalid_scope"
		}

		h.sendTokenError(c, errorCode, err.Error(), status)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// ListOAuth2Tokens godoc
// @Summary      List OAuth2 tokens
// @Description  Retrieve all OAuth2 tokens, optionally filtered by service_id
// @Tags         oauth2
// @Produce      json
// @Param        service_id  query  string  false  "Service ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      405  {object}  map[string]string
// @Router       /oauth2/tokens [get]

// CreateOAuth2Token godoc
// @Summary      Create an OAuth2 token
// @Description  Create a new OAuth2 token with access, refresh token, and expiration details
// @Tags         oauth2
// @Accept       json
// @Produce      json
// @Param        token  body  object  true  "OAuth2 Token creation payload"
// @Success      201  {object}  models.OAuth2Token
// @Failure      400  {object}  map[string]string
// @Failure      405  {object}  map[string]string
// @Router       /oauth2/tokens [post]
func (h *OAuth2Handler) OAuth2Tokens(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		h.listTokens(c)
	case "POST":
		h.createToken(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

// GetOAuth2TokenByID godoc
// @Summary      Get a token by ID
// @Description  Retrieve a specific OAuth2 token by ID
// @Tags         oauth2
// @Produce      json
// @Param        token_id  path  string  true  "Token ID"
// @Success      200  {object}  models.OAuth2Token
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /oauth2/tokens/{token_id} [get]

// UpdateOAuth2TokenByID godoc
// @Summary      Update a token by ID
// @Description  Update fields of a specific OAuth2 token by ID
// @Tags         oauth2
// @Accept       json
// @Produce      json
// @Param        token_id  path  string  true  "Token ID"
// @Param        updates   body  object  true  "Fields to update (partial allowed)"
// @Success      200  {object}  models.OAuth2Token
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /oauth2/tokens/{token_id} [put]

// DeleteOAuth2TokenByID godoc
// @Summary      Delete a token by ID
// @Description  Delete a specific OAuth2 token by ID
// @Tags         oauth2
// @Produce      json
// @Param        token_id  path  string  true  "Token ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /oauth2/tokens/{token_id} [delete]
func (h *OAuth2Handler) OAuth2TokenByID(c *gin.Context) {
	tokenID := c.Param("token_id")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token_id required"})
		return
	}

	switch c.Request.Method {
	case "GET":
		h.getToken(c, tokenID)
	case "PUT":
		h.updateToken(c, tokenID)
	case "DELETE":
		h.deleteToken(c, tokenID)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func (h *OAuth2Handler) listTokens(c *gin.Context) {
	var tokens []models.OAuth2Token

	query := h.db.Preload("Credential")
	if serviceID := c.Query("service_id"); serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}

	if err := query.Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(tokens),
		"data":  tokens,
	})
}

func (h *OAuth2Handler) createToken(c *gin.Context) {
	var req struct {
		Credential struct {
			ID uuid.UUID `json:"id" binding:"required"`
		} `json:"credential" binding:"required"`
		AccessToken            string `json:"access_token"`
		RefreshToken           string `json:"refresh_token"`
		AccessTokenExpiration  int    `json:"access_token_expiration"`
		RefreshTokenExpiration int    `json:"refresh_token_expiration"`
		Scope                  string `json:"scope"`
		AuthenticatedUserID    string `json:"authenticated_userid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token := &models.OAuth2Token{
		AccessToken:            req.AccessToken,
		RefreshToken:           req.RefreshToken,
		AccessTokenExpiration:  utils.GetCurrentTS().Add(time.Duration(req.AccessTokenExpiration) * time.Second),
		RefreshTokenExpiration: utils.GetCurrentTS().Add(time.Duration(req.RefreshTokenExpiration) * time.Second),
		Scope:                  req.Scope,
		AuthenticatedUserID:    req.AuthenticatedUserID,
		CredentialID:           req.Credential.ID,
	}

	if err := h.db.Create(token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusCreated, token)
}

func (h *OAuth2Handler) getToken(c *gin.Context, tokenID string) {
	var token models.OAuth2Token
	if err := h.db.Preload("Credential").Where("id = ?", tokenID).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch token"})
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *OAuth2Handler) updateToken(c *gin.Context, tokenID string) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var token models.OAuth2Token
	if err := h.db.Where("id = ?", tokenID).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch token"})
		return
	}

	if err := h.db.Model(&token).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update token"})
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *OAuth2Handler) deleteToken(c *gin.Context, tokenID string) {
	result := h.db.Where("id = ?", tokenID).Delete(&models.OAuth2Token{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete token"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// sendErrorRedirect actualizado para hacer redirects HTTP reales
func (h *OAuth2Handler) sendErrorRedirect(c *gin.Context, errorCode, description, redirectURI, state string) {
	if redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             errorCode,
			"error_description": description,
		})
		return
	}

	errorURL := fmt.Sprintf("%s?error=%s&error_description=%s", redirectURI, errorCode, description)
	if state != "" {
		errorURL += "&state=" + state
	}

	// Hacer redirect HTTP real en lugar de devolver JSON
	c.Redirect(http.StatusFound, errorURL)
}

func (h *OAuth2Handler) sendTokenError(c *gin.Context, errorCode, description string, status int) {
	c.JSON(status, gin.H{
		"error":             errorCode,
		"error_description": description,
	})
}

// ValidateToken godoc
// @Summary      Middleware to validate OAuth2 tokens
// @Description  Checks for a valid Bearer token in the Authorization header. Aborts with 401 if the token is missing, invalid, or expired. Sets user and token info in Gin context on success.
// @Tags         oauth2
// @Security     ApiKeyAuth
// @Produce      json
// @Failure      401  {object}  map[string]string
func (h *OAuth2Handler) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenValue := parts[1]

		var token models.OAuth2Token
		if err := h.db.Preload("Credential").Where("access_token = ?", tokenValue).First(&token).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if token.IsExpired() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("client_id", token.Credential.ClientID)
		c.Set("authenticated_userid", token.AuthenticatedUserID)
		c.Set("scope", token.Scope)

		c.Next()
	}
}
