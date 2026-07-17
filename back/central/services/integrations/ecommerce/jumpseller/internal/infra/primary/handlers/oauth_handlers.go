package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

const jumpsellerTypeID = uint(33)

func (h *jumpsellerHandler) platformConfig(ctx context.Context, field string, testMode bool) string {
	if h.coreIntegration == nil {
		return ""
	}
	creds, err := h.coreIntegration.GetCachedPlatformCredentials(ctx, jumpsellerTypeID)
	if err != nil || creds == nil {
		return ""
	}
	if testMode {
		if v, ok := creds["test_"+field].(string); ok && v != "" {
			return v
		}
	}
	if v, ok := creds[field].(string); ok && v != "" {
		return v
	}
	return ""
}

func (h *jumpsellerHandler) resolveRedirectURI(c *gin.Context, testMode bool) string {
	if v := h.platformConfig(c.Request.Context(), "redirect_uri", testMode); v != "" {
		return v
	}
	scheme := "https"
	if h.config != nil && h.config.Get("APP_ENV") == "development" {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/api/v1/jumpseller/callback", scheme, c.Request.Host)
}

func (h *jumpsellerHandler) resolveFrontendURL(c *gin.Context) string {
	if h.config != nil {
		if v := h.config.Get("FRONTEND_BASE_URL"); v != "" {
			return strings.TrimRight(v, "/")
		}
	}
	scheme := "https"
	if h.config != nil && h.config.Get("APP_ENV") == "development" {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

type initiateOAuthRequest struct {
	IntegrationName string `json:"integration_name" binding:"required"`
	BusinessID      uint   `json:"business_id"`
	IsTesting       bool   `json:"is_testing"`
}

type initiateOAuthResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	AuthorizationURL string `json:"authorization_url,omitempty"`
	State            string `json:"state,omitempty"`
	Error            string `json:"error,omitempty"`
}

func (h *jumpsellerHandler) InitiateOAuth(c *gin.Context) {
	var req initiateOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, initiateOAuthResponse{Success: false, Message: "Datos de entrada invalidos", Error: err.Error()})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, initiateOAuthResponse{Success: false, Message: "Usuario no autenticado", Error: "token invalido o ausente"})
		return
	}

	businessID := c.GetUint("business_id")
	if businessID == 0 && req.BusinessID > 0 {
		businessID = req.BusinessID
	}

	testMode := req.IsTesting
	clientID := h.platformConfig(c.Request.Context(), "client_id", testMode)
	if clientID == "" {
		c.JSON(http.StatusInternalServerError, initiateOAuthResponse{
			Success: false,
			Message: "Falta configurar el APP ID de la aplicacion de Jumpseller",
			Error:   "client_id ausente en las credenciales de plataforma del tipo Jumpseller",
		})
		return
	}

	scopes := h.platformConfig(c.Request.Context(), "scopes", testMode)
	if scopes == "" {
		scopes = domain.DefaultScopes
	}

	redirectURI := h.resolveRedirectURI(c, testMode)

	state, err := generateRandomToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, initiateOAuthResponse{Success: false, Message: "Error al generar el token de seguridad", Error: err.Error()})
		return
	}

	storeOAuthState(state, &oAuthStateData{
		IntegrationName: req.IntegrationName,
		UserID:          userID,
		BusinessID:      businessID,
		IsTesting:       testMode,
		Expiry:          time.Now().Add(10 * time.Minute),
	})

	authURL := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		domain.OAuthAuthorizeURL,
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(scopes),
		url.QueryEscape(state),
	)

	h.logger.Info(c.Request.Context()).
		Uint("user_id", userID).
		Uint("business_id", businessID).
		Str("state", state).
		Msg("Jumpseller OAuth iniciado")

	c.JSON(http.StatusOK, initiateOAuthResponse{
		Success:          true,
		Message:          "URL de autorizacion generada",
		AuthorizationURL: authURL,
		State:            state,
	})
}

type verifyAppResponse struct {
	Success    bool   `json:"success"`
	Configured bool   `json:"configured"`
	Message    string `json:"message"`
}

func (h *jumpsellerHandler) VerifyApp(c *gin.Context) {
	testMode := c.Query("is_testing") == "true"

	clientID := h.platformConfig(c.Request.Context(), "client_id", testMode)
	clientSecret := h.platformConfig(c.Request.Context(), "client_secret", testMode)

	var missing []string
	if clientID == "" {
		missing = append(missing, "APP ID")
	}
	if clientSecret == "" {
		missing = append(missing, "APP SECRET")
	}

	if len(missing) > 0 {
		c.JSON(http.StatusOK, verifyAppResponse{
			Success:    false,
			Configured: false,
			Message:    "Falta configurar en el tipo de integracion: " + strings.Join(missing, ", "),
		})
		return
	}

	c.JSON(http.StatusOK, verifyAppResponse{
		Success:    true,
		Configured: true,
		Message:    "La aplicacion de Jumpseller esta configurada. Ya puedes conectar.",
	})
}

func (h *jumpsellerHandler) OAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if errParam := c.Query("error"); errParam != "" {
		h.redirectError(c, "Autorizacion denegada por el usuario")
		return
	}
	if code == "" || state == "" {
		h.redirectError(c, "Parametros faltantes en la solicitud")
		return
	}

	stateData, ok := consumeOAuthState(state)
	if !ok {
		h.redirectError(c, "Token de seguridad invalido o expirado")
		return
	}

	testMode := stateData.IsTesting
	clientID := h.platformConfig(c.Request.Context(), "client_id", testMode)
	clientSecret := h.platformConfig(c.Request.Context(), "client_secret", testMode)
	redirectURI := h.resolveRedirectURI(c, testMode)

	tokenResp, err := exchangeCodeForToken(c.Request.Context(), clientID, clientSecret, code, redirectURI)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Msg("Error al intercambiar code por token en Jumpseller")
		h.redirectError(c, "Error al obtener el token de acceso")
		return
	}

	exchangeToken, err := generateRandomToken(16)
	if err != nil {
		h.redirectError(c, "Error al generar el token de intercambio")
		return
	}

	storeExchangeToken(exchangeToken, tokenExchangeData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		IsTesting:    testMode,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		Expiry:       time.Now().Add(5 * time.Minute),
	})

	redirectURL := fmt.Sprintf(
		"%s/integrations?jumpseller_oauth=success&integration_name=%s&state=%s&user_id=%d&business_id=%d&is_testing=%t&exchange_token=%s",
		h.resolveFrontendURL(c),
		url.QueryEscape(stateData.IntegrationName),
		url.QueryEscape(state),
		stateData.UserID,
		stateData.BusinessID,
		testMode,
		exchangeToken,
	)

	c.Redirect(http.StatusFound, redirectURL)
}

func (h *jumpsellerHandler) redirectError(c *gin.Context, message string) {
	redirectURL := fmt.Sprintf("%s/integrations?jumpseller_oauth=error&message=%s", h.resolveFrontendURL(c), url.QueryEscape(message))
	c.Redirect(http.StatusFound, redirectURL)
}

func exchangeCodeForToken(ctx context.Context, clientID, clientSecret, code, redirectURI string) (*domain.TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, domain.OAuthTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("building token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("reading token response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(raw))
	}

	var parsed struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		CreatedAt    int64  `json:"created_at"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}
	if parsed.AccessToken == "" {
		return nil, fmt.Errorf("empty access_token in response")
	}

	return &domain.TokenResponse{
		AccessToken:  parsed.AccessToken,
		RefreshToken: parsed.RefreshToken,
		TokenType:    parsed.TokenType,
		ExpiresIn:    parsed.ExpiresIn,
		CreatedAt:    parsed.CreatedAt,
	}, nil
}

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
