package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

const meliTypeID = uint(3)

const meliTokenURL = "https://api.mercadolibre.com/oauth/token"

const defaultMeliAuthDomain = "auth.mercadolibre.com.co"

const defaultMeliScopes = "offline_access read write"

var meliConfigEnvFallback = map[string]string{
	"client_id":     "MELI_CLIENT_ID",
	"client_secret": "MELI_CLIENT_SECRET",
	"redirect_uri":  "MELI_REDIRECT_URI",
	"frontend_url":  "FRONTEND_URL",
	"auth_domain":   "MELI_AUTH_DOMAIN",
	"scopes":        "MELI_SCOPES",
}

func (h *meliHandler) platformCred(ctx context.Context, field string) string {
	if h.coreIntegration != nil {
		creds, err := h.coreIntegration.GetPlatformCredential(ctx, fmt.Sprintf("%d", meliTypeID), field)
		if err == nil && creds != "" {
			return creds
		}
	}
	return ""
}

func (h *meliHandler) envCred(field string) string {
	envKey, ok := meliConfigEnvFallback[field]
	if !ok {
		envKey = "MELI_" + strings.ToUpper(field)
	}
	return h.config.Get(envKey)
}

func (h *meliHandler) getMeliConfig(ctx context.Context, field string, testMode bool) string {
	if testMode {
		if v := h.platformCred(ctx, "test_"+field); v != "" {
			return v
		}
		if v := h.envCred("test_" + field); v != "" {
			return v
		}
	}
	if v := h.platformCred(ctx, field); v != "" {
		return v
	}
	return h.envCred(field)
}

type InitiateOAuthRequest struct {
	IntegrationName string `json:"integration_name" binding:"required"`
	BusinessID      uint   `json:"business_id"`
	IsTesting       bool   `json:"is_testing"`
}

type InitiateOAuthResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	AuthorizationURL string `json:"authorization_url,omitempty"`
	State            string `json:"state,omitempty"`
	Error            string `json:"error,omitempty"`
}

func (h *meliHandler) InitiateOAuthHandler(c *gin.Context) {
	var req InitiateOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, InitiateOAuthResponse{
			Success: false,
			Message: "Datos de entrada invalidos",
			Error:   err.Error(),
		})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, InitiateOAuthResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "token de autenticacion invalido o ausente",
		})
		return
	}
	businessID := c.GetUint("business_id")
	if businessID == 0 && req.BusinessID > 0 {
		businessID = req.BusinessID
	}

	testMode := req.IsTesting

	clientID := h.getMeliConfig(c.Request.Context(), "client_id", testMode)
	if clientID == "" {
		c.JSON(http.StatusInternalServerError, InitiateOAuthResponse{
			Success: false,
			Message: "Error de configuracion: falta el Client ID de MercadoLibre",
			Error:   "client_id is missing (configurar credenciales de plataforma o MELI_CLIENT_ID)",
		})
		return
	}

	redirectURI := h.getMeliConfig(c.Request.Context(), "redirect_uri", testMode)
	if redirectURI == "" {
		scheme := "https"
		if h.config.Get("APP_ENV") == "development" {
			scheme = "http"
		}
		redirectURI = fmt.Sprintf("%s://%s/api/v1/meli/callback", scheme, c.Request.Host)
	}

	authDomain := h.getMeliConfig(c.Request.Context(), "auth_domain", testMode)
	if authDomain == "" {
		authDomain = defaultMeliAuthDomain
	}

	scopes := h.getMeliConfig(c.Request.Context(), "scopes", testMode)
	if scopes == "" {
		scopes = defaultMeliScopes
	}

	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, InitiateOAuthResponse{
			Success: false,
			Message: "Error al generar token de seguridad",
			Error:   err.Error(),
		})
		return
	}

	verifier, challenge, err := generatePKCE()
	if err != nil {
		c.JSON(http.StatusInternalServerError, InitiateOAuthResponse{
			Success: false,
			Message: "Error al generar PKCE",
			Error:   err.Error(),
		})
		return
	}

	storeOAuthState(state, &OAuthStateData{
		IntegrationName: req.IntegrationName,
		UserID:          userID,
		BusinessID:      businessID,
		CodeVerifier:    verifier,
		IsTesting:       testMode,
		Expiry:          time.Now().Add(10 * time.Minute),
	})

	authURL := fmt.Sprintf(
		"https://%s/authorization?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&code_challenge=%s&code_challenge_method=S256&state=%s",
		authDomain,
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(scopes),
		url.QueryEscape(challenge),
		url.QueryEscape(state),
	)

	h.logger.Info(c.Request.Context()).
		Uint("user_id", userID).
		Uint("business_id", businessID).
		Str("state", state).
		Msg("MeLi OAuth iniciado")

	c.JSON(http.StatusOK, InitiateOAuthResponse{
		Success:          true,
		Message:          "URL de autorizacion generada",
		AuthorizationURL: authURL,
		State:            state,
	})
}

func (h *meliHandler) OAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errParam := c.Query("error")

	if errParam != "" {
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
	clientID := h.getMeliConfig(c.Request.Context(), "client_id", testMode)
	clientSecret := h.getMeliConfig(c.Request.Context(), "client_secret", testMode)
	redirectURI := h.getMeliConfig(c.Request.Context(), "redirect_uri", testMode)
	if redirectURI == "" {
		scheme := "https"
		if h.config.Get("APP_ENV") == "development" {
			scheme = "http"
		}
		redirectURI = fmt.Sprintf("%s://%s/api/v1/meli/callback", scheme, c.Request.Host)
	}

	tokenResp, err := exchangeCodeForToken(clientID, clientSecret, code, redirectURI, stateData.CodeVerifier)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Msg("Error al intercambiar code por token en MeLi")
		h.redirectError(c, "Error al obtener token de acceso")
		return
	}

	exchangeTokenBytes := make([]byte, 16)
	rand.Read(exchangeTokenBytes)
	exchangeToken := hex.EncodeToString(exchangeTokenBytes)

	storeExchangeToken(exchangeToken, TokenExchangeData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		SellerID:     tokenResp.UserID,
		IsTesting:    testMode,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		Expiry:       time.Now().Add(5 * time.Minute),
	})

	integrationCode := generateIntegrationCode(stateData.IntegrationName)

	frontendURL := h.getMeliConfig(c.Request.Context(), "frontend_url", testMode)
	if frontendURL == "" {
		scheme := "https"
		if h.config.Get("APP_ENV") == "development" {
			scheme = "http"
		}
		frontendURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}

	redirectURL := fmt.Sprintf(
		"%s/integrations?meli_oauth=success&integration_name=%s&integration_code=%s&state=%s&user_id=%d&business_id=%d&seller_id=%d&is_testing=%t&exchange_token=%s",
		frontendURL,
		url.QueryEscape(stateData.IntegrationName),
		url.QueryEscape(integrationCode),
		url.QueryEscape(state),
		stateData.UserID,
		stateData.BusinessID,
		tokenResp.UserID,
		testMode,
		exchangeToken,
	)

	c.Redirect(http.StatusFound, redirectURL)
}

func (h *meliHandler) redirectError(c *gin.Context, message string) {
	frontendURL := h.getMeliConfig(c.Request.Context(), "frontend_url", false)
	if frontendURL == "" {
		scheme := "https"
		if h.config.Get("APP_ENV") == "development" {
			scheme = "http"
		}
		frontendURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}
	redirectURL := fmt.Sprintf("%s/integrations?meli_oauth=error&message=%s", frontendURL, url.QueryEscape(message))
	c.Redirect(http.StatusFound, redirectURL)
}

type meliTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserID       int64  `json:"user_id"`
}

func exchangeCodeForToken(clientID, clientSecret, code, redirectURI, codeVerifier string) (*meliTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest(http.MethodPost, meliTokenURL, strings.NewReader(data.Encode()))
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned %d", resp.StatusCode)
	}

	var result meliTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}
	if result.AccessToken == "" {
		return nil, fmt.Errorf("empty access_token in response")
	}
	return &result, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func generatePKCE() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	verifier := base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	return verifier, challenge, nil
}

func generateIntegrationCode(name string) string {
	code := strings.ToLower(strings.TrimSpace(name))
	code = strings.ReplaceAll(code, " ", "_")
	code = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, code)
	return fmt.Sprintf("mercado_libre_%s_%d", code, time.Now().Unix())
}
