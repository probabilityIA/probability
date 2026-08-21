package handlers

import (
	"bytes"
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
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func (h *tiendanubeHandler) platformConfig(ctx context.Context, field string, testMode bool) string {
	if h.coreIntegration == nil {
		return ""
	}
	creds, err := h.coreIntegration.GetCachedPlatformCredentials(ctx, domain.TypeID)
	if err != nil || creds == nil {
		return ""
	}
	if testMode {
		if v, ok := creds["test_"+field].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	if v, ok := creds[field].(string); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return ""
}

func (h *tiendanubeHandler) resolveFrontendURL(c *gin.Context) string {
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

func (h *tiendanubeHandler) InitiateOAuth(c *gin.Context) {
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
	if businessID == 0 {
		c.JSON(http.StatusBadRequest, initiateOAuthResponse{
			Success: false,
			Message: "Debes seleccionar un negocio antes de conectar",
			Error:   "business_id es requerido",
		})
		return
	}

	testMode := req.IsTesting
	clientID := h.platformConfig(c.Request.Context(), "client_id", testMode)
	if clientID == "" {
		c.JSON(http.StatusInternalServerError, initiateOAuthResponse{
			Success: false,
			Message: "Falta configurar el Client ID de la aplicacion de Tiendanube",
			Error:   "client_id ausente en las credenciales de plataforma del tipo Tiendanube",
		})
		return
	}

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
		"%s?state=%s",
		fmt.Sprintf(domain.OAuthAuthorizeURLTemplate, url.PathEscape(clientID)),
		url.QueryEscape(state),
	)

	h.logger.Info(c.Request.Context()).
		Uint("user_id", userID).
		Uint("business_id", businessID).
		Bool("is_testing", testMode).
		Msg("Tiendanube OAuth iniciado")

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

func (h *tiendanubeHandler) VerifyApp(c *gin.Context) {
	testMode := c.Query("is_testing") == "true"

	clientID := h.platformConfig(c.Request.Context(), "client_id", testMode)
	clientSecret := h.platformConfig(c.Request.Context(), "client_secret", testMode)

	var missing []string
	if clientID == "" {
		missing = append(missing, "Client ID")
	}
	if clientSecret == "" {
		missing = append(missing, "Client Secret")
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
		Message:    "La aplicacion de Tiendanube esta configurada. Ya puedes conectar.",
	})
}

func (h *tiendanubeHandler) OAuthCallback(c *gin.Context) {
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

	tokenResp, err := exchangeCodeForToken(c.Request.Context(), clientID, clientSecret, code)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Msg("Error al intercambiar code por token en Tiendanube")
		h.redirectError(c, "Error al obtener el token de acceso")
		return
	}

	exchangeToken, err := generateRandomToken(16)
	if err != nil {
		h.redirectError(c, "Error al generar el token de intercambio")
		return
	}

	storeExchangeToken(exchangeToken, tokenExchangeData{
		AccessToken: tokenResp.AccessToken,
		StoreID:     tokenResp.UserID,
		Scope:       tokenResp.Scope,
		IsTesting:   testMode,
		Expiry:      time.Now().Add(5 * time.Minute),
	})

	h.logger.Info(c.Request.Context()).
		Uint("business_id", stateData.BusinessID).
		Str("store_id", tokenResp.UserID).
		Msg("Tiendanube OAuth completado")

	redirectURL := fmt.Sprintf(
		"%s/integrations?tiendanube_oauth=success&integration_name=%s&state=%s&user_id=%d&business_id=%d&is_testing=%t&exchange_token=%s",
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

func (h *tiendanubeHandler) redirectError(c *gin.Context, message string) {
	redirectURL := fmt.Sprintf("%s/integrations?tiendanube_oauth=error&message=%s", h.resolveFrontendURL(c), url.QueryEscape(message))
	c.Redirect(http.StatusFound, redirectURL)
}

type tokenResponse struct {
	AccessToken string
	UserID      string
	Scope       string
}

func exchangeCodeForToken(ctx context.Context, clientID, clientSecret, code string) (*tokenResponse, error) {
	payload, err := json.Marshal(map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
	})
	if err != nil {
		return nil, fmt.Errorf("encoding token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, domain.OAuthTokenURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("building token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
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

	return parseTokenResponse(raw)
}

func parseTokenResponse(raw []byte) (*tokenResponse, error) {
	var parsed struct {
		AccessToken string      `json:"access_token"`
		TokenType   string      `json:"token_type"`
		Scope       string      `json:"scope"`
		UserID      json.Number `json:"user_id"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}
	if strings.TrimSpace(parsed.AccessToken) == "" {
		return nil, fmt.Errorf("empty access_token in response")
	}
	storeID := strings.Trim(strings.TrimSpace(parsed.UserID.String()), `"`)
	if storeID == "" {
		return nil, fmt.Errorf("empty user_id in response")
	}

	return &tokenResponse{
		AccessToken: parsed.AccessToken,
		UserID:      storeID,
		Scope:       parsed.Scope,
	}, nil
}

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
