package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
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

// OAuthStateData almacena información temporal durante el flujo OAuth
type OAuthStateData struct {
	IntegrationName string
	ShopDomain      string
	UserID          uint
	BusinessID      uint
	Expiry          time.Time
}

// OAuthStateStore almacena temporalmente los estados CSRF para validación
// En producción, esto debería usar Redis o similar
var oauthStateStore = make(map[string]*OAuthStateData)

// InitiateOAuthRequest representa la solicitud para iniciar OAuth
type InitiateOAuthRequest struct {
	IntegrationName string `json:"integration_name" binding:"required"`
	ShopDomain      string `json:"shop_domain" binding:"required"`
}

// InitiateOAuthResponse representa la respuesta con la URL de autorización
type InitiateOAuthResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	AuthorizationURL string `json:"authorization_url,omitempty"`
	State            string `json:"state,omitempty"`
	Error            string `json:"error,omitempty"`
}

// InitiateOAuthHandler inicia el flujo OAuth de Shopify
//
//	@Summary		Iniciar OAuth con Shopify
//	@Description	Genera la URL de autorización para que el usuario autorice la app en Shopify
//	@Tags			Shopify OAuth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		InitiateOAuthRequest	true	"Datos para iniciar OAuth"
//	@Success		200		{object}	InitiateOAuthResponse
//	@Failure		400		{object}	InitiateOAuthResponse
//	@Failure		401		{object}	InitiateOAuthResponse
//	@Failure		500		{object}	InitiateOAuthResponse
//	@Router			/integrations/shopify/connect [post]
func (h *ShopifyHandler) InitiateOAuthHandler(c *gin.Context) {
	var req InitiateOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos de entrada para OAuth")
		c.JSON(http.StatusBadRequest, InitiateOAuthResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Validar que el usuario esté autenticado
	userID, exists := middleware.GetUserID(c)
	if !exists {
		h.logger.Error().Msg("Intento de iniciar OAuth sin usuario autenticado")
		c.JSON(http.StatusUnauthorized, InitiateOAuthResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "token de autenticación inválido o ausente",
		})
		return
	}

	// Obtener business_id del contexto
	businessID := c.GetUint("business_id")

	// Normalizar el dominio de la tienda
	shopDomain := normalizeShopDomain(req.ShopDomain)

	// Generar state CSRF
	state, err := generateState()
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al generar state CSRF")
		c.JSON(http.StatusInternalServerError, InitiateOAuthResponse{
			Success: false,
			Message: "Error al generar token de seguridad",
			Error:   err.Error(),
		})
		return
	}

	// Almacenar state con metadata (expira en 10 minutos)
	oauthStateStore[state] = &OAuthStateData{
		IntegrationName: req.IntegrationName,
		ShopDomain:      shopDomain,
		UserID:          userID,
		BusinessID:      businessID,
		Expiry:          time.Now().Add(10 * time.Minute),
	}

	// Obtener configuración de OAuth desde variables de entorno
	clientID := h.config.Get("SHOPIFY_CLIENT_ID")
	redirectURI := h.config.Get("SHOPIFY_REDIRECT_URI")
	scopes := h.config.Get("SHOPIFY_SCOPES")

	if clientID == "" || redirectURI == "" || scopes == "" {
		h.logger.Error().Msg("Configuración OAuth incompleta")
		c.JSON(http.StatusInternalServerError, InitiateOAuthResponse{
			Success: false,
			Message: "Configuración OAuth incompleta",
			Error:   "SHOPIFY_CLIENT_ID, SHOPIFY_REDIRECT_URI o SHOPIFY_SCOPES no configurados",
		})
		return
	}

	// Construir URL de autorización de Shopify
	authURL := fmt.Sprintf(
		"https://%s/admin/oauth/authorize?client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		shopDomain,
		url.QueryEscape(clientID),
		url.QueryEscape(scopes),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)

	h.logger.Info().
		Uint("user_id", userID).
		Str("shop_domain", shopDomain).
		Str("state", state).
		Msg("OAuth iniciado exitosamente")

	c.JSON(http.StatusOK, InitiateOAuthResponse{
		Success:          true,
		Message:          "URL de autorización generada",
		AuthorizationURL: authURL,
		State:            state,
	})
}

// OAuthCallbackHandler maneja el callback de OAuth desde Shopify
//
//	@Summary		Callback de OAuth
//	@Description	Recibe el código de autorización de Shopify y lo intercambia por un access token
//	@Tags			Shopify OAuth
//	@Accept			json
//	@Produce		json
//	@Param			code	query	string	true	"Código de autorización"
//	@Param			shop	query	string	true	"Dominio de la tienda"
//	@Param			state	query	string	true	"State CSRF"
//	@Param			hmac	query	string	true	"HMAC signature"
//	@Success		302		{string}	string	"Redirección al frontend con datos de integración"
//	@Failure		400		{object}	map[string]interface{}	"Error en validación"
//	@Router			/shopify/callback [get]
func (h *ShopifyHandler) OAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	shop := c.Query("shop")
	state := c.Query("state")
	hmacParam := c.Query("hmac")

	// Validar parámetros requeridos
	if code == "" || shop == "" || state == "" || hmacParam == "" {
		h.logger.Error().Msg("Parámetros faltantes en callback OAuth")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Parámetros faltantes en la solicitud",
		})
		return
	}

	// Validar state CSRF
	stateData, exists := oauthStateStore[state]
	if !exists || time.Now().After(stateData.Expiry) {
		h.logger.Error().Str("state", state).Msg("State CSRF inválido o expirado")
		delete(oauthStateStore, state)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Token de seguridad inválido o expirado",
		})
		return
	}

	// Eliminar state usado
	delete(oauthStateStore, state)

	// Validar HMAC
	if !h.validateHMAC(c.Request.URL.Query(), hmacParam) {
		h.logger.Error().Msg("HMAC inválido en callback OAuth")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Firma HMAC inválida",
		})
		return
	}

	// Intercambiar código por access token
	clientID := h.config.Get("SHOPIFY_CLIENT_ID")
	clientSecret := h.config.Get("SHOPIFY_CLIENT_SECRET")

	accessToken, scope, err := h.exchangeCodeForToken(shop, code, clientID, clientSecret)
	if err != nil {
		h.logger.Error().Err(err).Str("shop", shop).Msg("Error al intercambiar código por token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error al obtener token de acceso",
		})
		return
	}

	h.logger.Info().
		Str("shop", shop).
		Str("scope", scope).
		Str("integration_name", stateData.IntegrationName).
		Msg("Token de acceso obtenido exitosamente")

	// Codificar datos para el frontend
	integrationCode := generateIntegrationCode(stateData.IntegrationName)

	// Redirigir al frontend con los datos necesarios para crear la integración
	frontendURL := h.config.Get("WEBHOOK_BASE_URL")
	redirectURL := fmt.Sprintf(
		"%s/integrations?shopify_oauth=success&shop=%s&integration_name=%s&integration_code=%s&access_token=%s&user_id=%d&business_id=%d",
		frontendURL,
		url.QueryEscape(shop),
		url.QueryEscape(stateData.IntegrationName),
		url.QueryEscape(integrationCode),
		url.QueryEscape(accessToken),
		stateData.UserID,
		stateData.BusinessID,
	)

	h.logger.Info().
		Str("redirect_url", redirectURL).
		Uint("user_id", stateData.UserID).
		Uint("business_id", stateData.BusinessID).
		Msg("Redirigiendo al frontend después de OAuth exitoso")

	c.Redirect(http.StatusFound, redirectURL)
}

// exchangeCodeForToken intercambia el código de autorización por un access token
func (h *ShopifyHandler) exchangeCodeForToken(shop, code, clientID, clientSecret string) (string, string, error) {
	tokenURL := fmt.Sprintf("https://%s/admin/oauth/access_token", shop)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", "", fmt.Errorf("error al hacer request de token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("respuesta no exitosa: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	return result.AccessToken, result.Scope, nil
}

// generateState genera un state CSRF aleatorio
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// normalizeShopDomain normaliza el dominio de la tienda
func normalizeShopDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)

	// Remover protocolo si existe
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	// Asegurar que termine en .myshopify.com
	if !strings.HasSuffix(domain, ".myshopify.com") {
		if !strings.Contains(domain, ".") {
			domain = domain + ".myshopify.com"
		}
	}

	return domain
}

// validateHMAC valida la firma HMAC de Shopify
func (h *ShopifyHandler) validateHMAC(queryParams url.Values, receivedHMAC string) bool {
	clientSecret := h.config.Get("SHOPIFY_CLIENT_SECRET")

	// Crear una copia de los parámetros sin el HMAC
	params := url.Values{}
	for k, v := range queryParams {
		if k != "hmac" && k != "signature" {
			params[k] = v
		}
	}

	// Construir query string ordenado
	message := params.Encode()

	// Calcular HMAC
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(message))
	expectedHMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedHMAC), []byte(receivedHMAC))
}

// generateIntegrationCode genera un código único para la integración basado en el nombre
func generateIntegrationCode(name string) string {
	code := strings.ToLower(name)
	code = strings.TrimSpace(code)
	code = strings.ReplaceAll(code, " ", "_")
	// Remover caracteres especiales
	code = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, code)
	// Agregar timestamp para unicidad
	code = fmt.Sprintf("%s_%d", code, time.Now().Unix())
	return code
}
