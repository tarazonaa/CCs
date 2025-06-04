package services

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"log"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OAuth2Service struct {
	db     *gorm.DB
	config *config.Config
}

func NewOAuth2Service(db *gorm.DB, cfg *config.Config) *OAuth2Service {
	return &OAuth2Service{
		db:     db,
		config: cfg,
	}
}

type AuthorizeRequest struct {
	ResponseType        string `json:"response_type" form:"response_type"`
	ClientID            string `json:"client_id" form:"client_id"`
	RedirectURI         string `json:"redirect_uri" form:"redirect_uri"`
	Scope               string `json:"scope" form:"scope"`
	State               string `json:"state" form:"state"`
	CodeChallenge       string `json:"code_challenge" form:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method" form:"code_challenge_method"`
	// Kong-specific fields
	ProvisionKey        string `json:"provision_key" form:"provision_key"`
	AuthenticatedUserID string `json:"authenticated_userid" form:"authenticated_userid"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type"`
	ClientID     string `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri"`
	Code         string `json:"code" form:"code"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
	Scope        string `json:"scope" form:"scope"`
	// Kong-specific fields
	ProvisionKey        string `json:"provision_key" form:"provision_key"`
	AuthenticatedUserID string `json:"authenticated_userid" form:"authenticated_userid"`
	Username            string `json:"username" form:"username"`
	Password            string `json:"password" form:"password"`
}

type AuthorizeResponse struct {
	RedirectURI string `json:"redirect_uri"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func (s *OAuth2Service) Authorize(req *AuthorizeRequest) (*AuthorizeResponse, error) {
	log.Printf("DEBUG: Received provision_key: '%s'", req.ProvisionKey)
	log.Printf("DEBUG: Expected provision_key: '%s'", s.config.ProvisionKey)
	
	// Validate provision key
	if req.ProvisionKey != s.config.ProvisionKey {
		return nil, errors.New("invalid provision key")
	}

	var app models.OAuth2Application
	if err := s.db.Where("client_id = ?", req.ClientID).First(&app).Error; err != nil {
		return nil, errors.New("invalid client")
	}

	if !s.isValidRedirectURI(req.RedirectURI, app.RedirectURIs) {
		return nil, errors.New("invalid redirect URI")
	}

	switch req.ResponseType {
	case "code":
		return s.handleAuthorizationCodeFlow(req, &app)
	case "token":
		return s.handleImplicitFlow(req, &app)
	default:
		return nil, errors.New("unsupported response type")
	}
}

func (s *OAuth2Service) handleAuthorizationCodeFlow(req *AuthorizeRequest, app *models.OAuth2Application) (*AuthorizeResponse, error) {
	if !s.config.OAuth2.EnableAuthorizationCode {
		return nil, errors.New("authorization code flow is disabled")
	}

	
	if s.config.OAuth2.PKCERequired && req.CodeChallenge == "" {
		return nil, errors.New("PKCE is required")
	}

	authCode := &models.AuthorizationCode{
		ClientID:            req.ClientID,
		UserID:              s.parseUserID(req.AuthenticatedUserID),
		RedirectURI:         req.RedirectURI,
		Scopes:              strings.Fields(req.Scope),
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		ExpiresAt:           time.Now().Add(time.Duration(s.config.OAuth2.AuthCodeExpiration) * time.Second),
	}

	if err := s.db.Create(authCode).Error; err != nil {
		return nil, fmt.Errorf("failed to create authorization code: %w", err)
	}


	redirectURL, _ := url.Parse(req.RedirectURI)
	query := redirectURL.Query()
	query.Set("code", authCode.Code)
	if req.State != "" {
		query.Set("state", req.State)
	}
	redirectURL.RawQuery = query.Encode()

	return &AuthorizeResponse{
		RedirectURI: redirectURL.String(),
	}, nil
}

func (s *OAuth2Service) handleImplicitFlow(req *AuthorizeRequest, app *models.OAuth2Application) (*AuthorizeResponse, error) {
	if !s.config.OAuth2.EnableImplicitGrant {
		return nil, errors.New("implicit grant flow is disabled")
	}

	// creamos el token directo
	token := &models.OAuth2Token{
		TokenType:           "bearer",
		ExpiresIn:           s.config.OAuth2.AccessTokenExpiration,
		Scope:               req.Scope,
		AuthenticatedUserID: req.AuthenticatedUserID,
		CredentialID:        app.ID,
	}

	if err := s.db.Create(token).Error; err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	redirectURL, _ := url.Parse(req.RedirectURI)
	fragment := fmt.Sprintf("access_token=%s&token_type=bearer&expires_in=%d",
		token.AccessToken, token.ExpiresIn)
	if req.Scope != "" {
		fragment += "&scope=" + url.QueryEscape(req.Scope)
	}
	if req.State != "" {
		fragment += "&state=" + url.QueryEscape(req.State)
	}
	redirectURL.Fragment = fragment

	return &AuthorizeResponse{
		RedirectURI: redirectURL.String(),
	}, nil
}

// endpoint
func (s *OAuth2Service) Token(req *TokenRequest) (*TokenResponse, error) {
	switch req.GrantType {
	case "authorization_code":
		return s.handleAuthorizationCodeGrant(req)
	case "client_credentials":
		return s.handleClientCredentialsGrant(req)
	case "refresh_token":
		return s.handleRefreshTokenGrant(req)
	case "password":
		return s.handlePasswordGrant(req)
	default:
		return nil, errors.New("unsupported grant type")
	}
}

func (s *OAuth2Service) handleAuthorizationCodeGrant(req *TokenRequest) (*TokenResponse, error) {
	var app models.OAuth2Application
	if err := s.validateClient(req.ClientID, req.ClientSecret, &app); err != nil {
		return nil, err
	}

	var authCode models.AuthorizationCode
	if err := s.db.Where("code = ? AND client_id = ?", req.Code, req.ClientID).First(&authCode).Error; err != nil {
		return nil, errors.New("invalid authorization code")
	}

	if !authCode.IsValid() {
		return nil, errors.New("authorization code expired or already used")
	}

	if authCode.RedirectURI != req.RedirectURI {
		return nil, errors.New("redirect URI mismatch")
	}

	if authCode.CodeChallenge != "" {
		// PKCE validation would go here
		// For now, we'll skip the implementation details
	}

	// Mark code as used
	authCode.IsUsed = true
	s.db.Save(&authCode)

	// Create access token
	return s.createTokenResponse(&app, authCode.UserID, strings.Join(authCode.Scopes, " "))
}

func (s *OAuth2Service) handleClientCredentialsGrant(req *TokenRequest) (*TokenResponse, error) {
	if !s.config.OAuth2.EnableClientCredentials {
		return nil, errors.New("client credentials flow is disabled")
	}

	// Validate client credentials
	var app models.OAuth2Application
	if err := s.validateClient(req.ClientID, req.ClientSecret, &app); err != nil {
		return nil, err
	}
	return s.createTokenResponse(&app, 0, req.Scope)
}

func (s *OAuth2Service) validateClient(clientID, clientSecret string, app *models.OAuth2Application) error {
	if err := s.db.Where("client_id = ?", clientID).First(app).Error; err != nil {
		return errors.New("invalid client")
	}

	if app.ClientSecret != clientSecret {
		return errors.New("invalid client credentials")
	}

	return nil
}

func (s *OAuth2Service) createTokenResponse(app *models.OAuth2Application, userID uint, scope string) (*TokenResponse, error) {
	// Create access token
	accessToken := &models.OAuth2Token{
		TokenType:           "bearer",
		ExpiresIn:           s.config.OAuth2.AccessTokenExpiration,
		Scope:               scope,
		CredentialID:        app.ID,
	}

	if userID > 0 {
		accessToken.AuthenticatedUserID = fmt.Sprintf("%d", userID)
	}

	if err := s.db.Create(accessToken).Error; err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	response := &TokenResponse{
		AccessToken: accessToken.AccessToken,
		TokenType:   accessToken.TokenType,
		ExpiresIn:   accessToken.ExpiresIn,
		Scope:       accessToken.Scope,
	}

	// Create refresh token if enabled
	if s.config.OAuth2.RefreshTokenExpiration > 0 {
		refreshToken := &models.OAuth2Token{
			RefreshToken:        uuid.New().String(),
			TokenType:           "refresh",
			ExpiresIn:           s.config.OAuth2.RefreshTokenExpiration,
			CredentialID:        app.ID,
			AuthenticatedUserID: accessToken.AuthenticatedUserID,
		}

		if err := s.db.Create(refreshToken).Error; err == nil {
			response.RefreshToken = refreshToken.RefreshToken
		}
	}

	return response, nil
}

// Helper functions
func (s *OAuth2Service) isValidRedirectURI(uri string, validURIs []string) bool {
	for _, validURI := range validURIs {
		if uri == validURI {
			return true
		}
	}
	return false
}

func (s *OAuth2Service) parseUserID(userIDStr string) uint {
	
	if userIDStr == "" {
		return 0
	}
	
	return 1 // Placeholder
}


func (s *OAuth2Service) handleRefreshTokenGrant(req *TokenRequest) (*TokenResponse, error) {
	// Implementation for refresh token grant
	return nil, errors.New("refresh token grant not implemented yet")
}

func (s *OAuth2Service) handlePasswordGrant(req *TokenRequest) (*TokenResponse, error) {
	return nil, errors.New("password grant not implemented")
}