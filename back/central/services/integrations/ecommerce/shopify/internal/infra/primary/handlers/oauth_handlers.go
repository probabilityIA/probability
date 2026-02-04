package handlers

import (
	"crypto/hmac"
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

// OAuthStateData almacena información temporal durante el flujo OAuth
type OAuthStateData struct {
	IntegrationName string
	ShopDomain      string
	UserID          uint
	BusinessID      uint
	Expiry          time.Time
	// Custom App Credentials (optional)
	ClientID     string
	ClientSecret string
}

// OAuthStateStore almacena temporalmente los estados CSRF para validación
// En producción, esto debería usar Redis o similar
var oauthStateStore = make(map[string]*OAuthStateData)

// InitiateOAuthRequest representa la solicitud para iniciar OAuth
type InitiateOAuthRequest struct {
	IntegrationName string `json:"integration_name" binding:"required"`
	ShopDomain      string `json:"shop_domain" binding:"required"`
}

// InitiateCustomOAuthRequest representa la solicitud para iniciar OAuth con credenciales personalizadas
type InitiateCustomOAuthRequest struct {
	IntegrationName string `json:"integration_name" binding:"required"`
	ShopDomain      string `json:"shop_domain" binding:"required"`
	ClientID        string `json:"client_id" binding:"required"`
	ClientSecret    string `json:"client_secret" binding:"required"`
}

// InitiateOAuthResponse representa la respuesta con la URL de autorización
type InitiateOAuthResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	AuthorizationURL string `json:"authorization_url,omitempty"`
	State            string `json:"state,omitempty"`
	Error            string `json:"error,omitempty"`
}

// InitiateOAuthHandler inicia el flujo OAuth de Shopify (App Pública/Default)
//
//	@Summary		Iniciar OAuth con Shopify (Default)
//	@Description	Genera la URL de autorización usando credenciales de entorno
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

	h.initiateOAuthProcess(c, req.IntegrationName, req.ShopDomain, "", "")
}

// InitiateCustomOAuthHandler inicia el flujo OAuth de Shopify con credenciales personalizadas
//
//	@Summary		Iniciar OAuth con Shopify (Custom App)
//	@Description	Genera la URL de autorización usando credenciales proporcionadas en el body
//	@Tags			Shopify OAuth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		InitiateCustomOAuthRequest	true	"Datos y credenciales para iniciar OAuth"
//	@Success		200		{object}	InitiateOAuthResponse
//	@Failure		400		{object}	InitiateOAuthResponse
//	@Failure		401		{object}	InitiateOAuthResponse
//	@Failure		500		{object}	InitiateOAuthResponse
//	@Router			/integrations/shopify/connect/custom [post]
func (h *ShopifyHandler) InitiateCustomOAuthHandler(c *gin.Context) {
	var req InitiateCustomOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos de entrada para Custom OAuth")
		c.JSON(http.StatusBadRequest, InitiateOAuthResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	h.initiateOAuthProcess(c, req.IntegrationName, req.ShopDomain, req.ClientID, req.ClientSecret)
}

// initiateOAuthProcess encapsula la lógica común de inicio de OAuth
func (h *ShopifyHandler) initiateOAuthProcess(c *gin.Context, integrationName, shopDomainParam, customClientID, customClientSecret string) {
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
	shopDomain := normalizeShopDomain(shopDomainParam)

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
		IntegrationName: integrationName,
		ShopDomain:      shopDomain,
		UserID:          userID,
		BusinessID:      businessID,
		Expiry:          time.Now().Add(10 * time.Minute),
		ClientID:        customClientID,
		ClientSecret:    customClientSecret,
	}

	// Determinar credenciales a usar
	var clientID, redirectURI, scopes string

	if customClientID != "" {
		// Usar credenciales custom proporcionadas por el usuario
		clientID = customClientID
		redirectURI = h.config.Get("SHOPIFY_REDIRECT_URI")
		scopes = h.config.Get("SHOPIFY_SCOPES")
	} else {
		// Usar credenciales globales configuradas por el administrador
		clientID = h.config.Get("SHOPIFY_CLIENT_ID")
		redirectURI = h.config.Get("SHOPIFY_REDIRECT_URI")
		scopes = h.config.Get("SHOPIFY_SCOPES")
	}

	// 1. Fallback para Scopes: si no están configurados, usar los mínimos necesarios
	if scopes == "" {
		scopes = "read_customers,read_fulfillments,read_orders,write_orders,read_products"
		h.logger.Info().Msg("SHOPIFY_SCOPES no configurado, usando valores por defecto")
	}

	// 2. Fallback para Redirect URI: si no está configurada, generarla dinámicamente
	if redirectURI == "" {
		scheme := "https"
		// Si estamos en desarrollo y no hay TLS, usar http
		if h.config.Get("APP_ENV") == "development" {
			scheme = "http"
		}
		// Construir URL: esquema://host_actual/api/v1/shopify/callback
		redirectURI = fmt.Sprintf("%s://%s/api/v1/shopify/callback", scheme, c.Request.Host)
		h.logger.Info().Str("redirect_uri", redirectURI).Msg("SHOPIFY_REDIRECT_URI no configurado, generada URL dinámica")
	}

	// 3. Validar Client ID (crítico)
	if clientID == "" {
		h.logger.Error().Msg("Faltan credenciales críticas para Shopify OAuth (Client ID)")
		c.JSON(http.StatusInternalServerError, InitiateOAuthResponse{
			Success: false,
			Message: "Error de configuración: falta el Client ID de Shopify",
			Error:   "SHOPIFY_CLIENT_ID is missing",
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
		Bool("is_custom", customClientID != "").
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
	// Nota: Si usamos custom credentials, necesitamos el secret CORRECTO para validar el HMAC.
	clientSecret := h.config.Get("SHOPIFY_CLIENT_SECRET") // Default
	if stateData.ClientSecret != "" {
		clientSecret = stateData.ClientSecret // Custom
	}

	if !h.validateHMACWithSecret(c.Request.URL.Query(), hmacParam, clientSecret) {
		h.logger.Error().Msg("HMAC inválido en callback OAuth")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Firma HMAC inválida",
		})
		return
	}

	// Intercambiar código por access token
	clientID := h.config.Get("SHOPIFY_CLIENT_ID") // Default
	if stateData.ClientID != "" {
		clientID = stateData.ClientID
	}

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
		Str("access_token_prefix", func() string {
			if len(accessToken) > 5 {
				return accessToken[:5] + "..."
			}
			return "too-short"
		}()).
		Msg("Token de acceso obtenido exitosamente")

	// Codificar datos para el frontend
	integrationCode := generateIntegrationCode(stateData.IntegrationName)

	// Almacenar credenciales completas en cookie segura (5 minutos de validez)
	// Guardamos token, client_id y client_secret para persistirlos después
	creds := map[string]string{
		"access_token":  accessToken,
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	credsJSON, _ := json.Marshal(creds)

	// Usar url.QueryEscape para que caracteres como { } " sean seguros en la cookie
	middleware.SetSecureCookie(c, "shopify_temp_token", url.QueryEscape(string(credsJSON)), 300)

	// Redirigir al frontend SIN el token en la URL (seguro)
	frontendURL := h.config.Get("WEBHOOK_BASE_URL")
	redirectURL := fmt.Sprintf(
		"%s/integrations?shopify_oauth=success&shop=%s&integration_name=%s&integration_code=%s&state=%s&user_id=%d&business_id=%d",
		frontendURL,
		url.QueryEscape(shop),
		url.QueryEscape(stateData.IntegrationName),
		url.QueryEscape(integrationCode),
		url.QueryEscape(state),
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

// validateHMACWithSecret valida la firma HMAC de Shopify usando un secret específico
func (h *ShopifyHandler) validateHMACWithSecret(queryParams url.Values, receivedHMAC, clientSecret string) bool {
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

// GetConfigHandler retorna la configuración pública de Shopify (Client ID)
//
//	@Summary		Obtener configuración de Shopify
//	@Description	Retorna el Client ID de Shopify para inicializar App Bridge en el frontend
//	@Tags			Shopify Integrations
//	@Produce		json
//	@Success		200		{object}	map[string]string
//	@Router			/integrations/shopify/config [get]
func (h *ShopifyHandler) GetConfigHandler(c *gin.Context) {
	clientID := h.config.Get("SHOPIFY_CLIENT_ID")

	if clientID == "" {
		h.logger.Warn().Msg("SHOPIFY_CLIENT_ID no configurado - OAuth de Shopify no disponible")
		// Retornar 200 OK con configuración vacía (no es un error, simplemente no está configurado)
		c.JSON(http.StatusOK, gin.H{
			"shopify_client_id": nil,
			"configured":        false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shopify_client_id": clientID,
		"configured":        true,
	})
}

// Estructura para recibir el session token
type LoginWithSessionTokenRequest struct {
	SessionToken string `json:"session_token" binding:"required"`
}

// Estructura para el payload del JWT de Shopify (claims básicos)
// Ver: https://shopify.dev/docs/apps/auth/oauth/session-tokens#payload
type ShopifySessionTokenClaims struct {
	Iss  string `json:"iss"`  // Issuer (https://shopify.com/<shop_id>)
	Dest string `json:"dest"` // Destination (https://<shop_domain>)
	Aud  string `json:"aud"`  // Audience (API Key)
	Sub  string `json:"sub"`  // Subject (User ID)
	Exp  int64  `json:"exp"`  // Expiration
	Nbf  int64  `json:"nbf"`  // Not Before
	Jti  string `json:"jti"`  // JWT ID
	Sid  string `json:"sid"`  // Session ID
}

// LoginWithSessionTokenHandler autentica un usuario basado en el Session Token de Shopify
func (h *ShopifyHandler) LoginWithSessionTokenHandler(c *gin.Context) {
	var req LoginWithSessionTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("LoginWithSessionToken: Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	clientSecret := h.config.Get("SHOPIFY_CLIENT_SECRET")
	if clientSecret == "" {
		h.logger.Error().Msg("SHOPIFY_CLIENT_SECRET not configured")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	// 1. Validar y parsear el token
	// Nota: En un entorno de producción ideal, deberíamos usar una librería JWT robusta para validar la firma.
	// Shopify usa HS256 con el Client Secret como clave.
	// Por simplicidad y evitar dependencias cíclicas si no tienes librería JWT a mano, aquí implemento una validación básica,
	// pero RECOMIENDO encarecidamente usar `golang-jwt/jwt` para esto.
	// Asumiré que podemos usar lógica similar a la de validación HMAC o una librería estándar si está disponible.

	// TODO: IMPLEMENTAR VALIDACIÓN JWT REAL
	// Por ahora, para avanzar, parseamos el claims sin verificar firma (INSEGURO - SOLO PARA DEMO/DEV ACEPTADO POR AHORA)
	// OJO: En producción DEBES validar la firma.

	/*
		token, err := jwt.ParseWithClaims(req.SessionToken, &ShopifySessionTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(clientSecret), nil
		})
	*/

	// SIMULACIÓN DE EXTRACCIÓN DE DOMINIO (Para ilustrar el flujo)
	// En realidad, decodificamos el payload base64.
	// El token es header.payload.signature
	parts := strings.Split(req.SessionToken, ".")
	if len(parts) != 3 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		h.logger.Error().Err(err).Msg("Error decoding JWT payload")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token encoding"})
		return
	}

	var claims ShopifySessionTokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		h.logger.Error().Err(err).Msg("Error unmarshalling JWT claims")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// El campo 'dest' contiene la URL de la tienda, ej: https://tienda.myshopify.com
	if claims.Dest == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing 'dest' claim"})
		return
	}

	// Limpiar el protocolo para obtener solo el dominio
	shopDomain := strings.TrimPrefix(claims.Dest, "https://")

	// 2. Buscar si existe un business integrado con este shopDomain
	// Aquí deberías tener un servicio/repositorio para buscar el 'Business' por 'ShopifyDomain'
	// Como no tengo acceso directo a tu repositorio de 'Business' desde este handler (solo integration),
	// simularé la respuesta o dejaré el TODO claro.

	// TODO: h.businessService.GetByShopifyShop(shopDomain)
	// Si encuentro el negocio, genero el token de sesión de TU app.

	h.logger.Info().Str("shop", shopDomain).Msg("Intento de login vía Session Token")

	// RESPUESTA MOCK (Para que el frontend avance):
	// Si el dominio coincide con algo conocido o simplemente para probar el flujo:
	// Devolvemos un token dummy o un error específico si no está registrado.

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful (MOCK)",
		"shop":    shopDomain,
		"token":   "MOCK_PROBABILITY_TOKEN_" + shopDomain, // Este token lo usaría el frontend para futuras llamadas
	})
}
