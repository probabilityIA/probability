package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetOAuthTokenResponse representa la respuesta del endpoint que retorna el token OAuth
type GetOAuthTokenResponse struct {
	Success         bool   `json:"success"`
	AccessToken     string `json:"access_token,omitempty"`
	ClientID        string `json:"client_id,omitempty"`
	ClientSecret    string `json:"client_secret,omitempty"`
	Shop            string `json:"shop,omitempty"`
	IntegrationName string `json:"integration_name,omitempty"`
	IntegrationCode string `json:"integration_code,omitempty"`
	Error           string `json:"error,omitempty"`
}

// GetOAuthTokenHandler recupera el access token almacenado temporalmente en cookie
// después del flujo OAuth. Este endpoint debe ser llamado inmediatamente después
// del redirect para obtener el token de forma segura.
//
//	@Summary		Obtener token OAuth temporal
//	@Description	Recupera el access token de Shopify almacenado en cookie después del flujo OAuth
//	@Tags			Shopify OAuth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			state	query		string	true	"State CSRF del flujo OAuth"
//	@Success		200		{object}	GetOAuthTokenResponse
//	@Failure		401		{object}	GetOAuthTokenResponse
//	@Failure		404		{object}	GetOAuthTokenResponse
//	@Router			/integrations/shopify/oauth/token [get]
func (h *ShopifyHandler) GetOAuthTokenHandler(c *gin.Context) {
	state := c.Query("state")
	shop := c.Query("shop")
	integrationName := c.Query("integration_name")
	integrationCode := c.Query("integration_code")
	exchangeToken := c.Query("exchange_token")

	if state == "" || shop == "" {
		h.logger.Error().Msg("State o shop faltantes en petición de token")
		c.JSON(http.StatusBadRequest, GetOAuthTokenResponse{
			Success: false,
			Error:   "Parámetros requeridos faltantes",
		})
		return
	}

	// 1. Intentar recuperar por exchange_token (prioridad)
	if exchangeToken != "" {
		if data, ok := RetrieveExchangeToken(exchangeToken); ok {
			h.logger.Info().Str("shop", shop).Msg("Credenciales recuperadas vía token de intercambio")

			c.JSON(http.StatusOK, GetOAuthTokenResponse{
				Success:         true,
				AccessToken:     data.AccessToken,
				ClientID:        data.ClientID,
				ClientSecret:    data.ClientSecret,
				Shop:            shop,
				IntegrationName: integrationName,
				IntegrationCode: integrationCode,
			})
			return
		}
		// Si falla, logueamos y seguimos con cookie como fallback
		h.logger.Warn().Str("exchange_token", exchangeToken).Msg("Token de intercambio inválido o expirado")
	}

	// 2. Intentar recuperar por cookie (fallback)
	cookieData, err := c.Cookie("shopify_temp_token")
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("state", state).
			Msg("Cookie de token no encontrada o expirada")

		// Usamos 410 Gone (no 404) para distinguir "token expirado" de "ruta no encontrada"
		c.JSON(http.StatusGone, GetOAuthTokenResponse{
			Success: false,
			Error:   "El token de autorización expiró o ya fue consumido. Por favor inicia el proceso de conexión con Shopify nuevamente.",
		})
		return
	}

	// URL-decode la cookie
	decodedData, err := url.QueryUnescape(cookieData)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Error al de-escapar cookie de token, intentando usar raw")
		decodedData = cookieData
	}

	// Decodificar JSON de la cookie
	var creds struct {
		AccessToken  string `json:"access_token"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	if strings.HasPrefix(decodedData, "{") {
		if err := json.Unmarshal([]byte(decodedData), &creds); err != nil {
			h.logger.Error().Err(err).Msg("Error al decodificar JSON de credenciales")
			creds.AccessToken = decodedData
		}
	} else {
		creds.AccessToken = decodedData
	}

	// Borrar la cookie
	host := c.Request.Host
	domainName := ".probabilityia.com.co"

	if strings.Contains(host, "localhost") ||
		strings.Contains(host, "127.0.0.1") ||
		strings.Contains(host, "ngrok") ||
		strings.Contains(host, ".dev") {
		domainName = ""
	}

	c.SetCookie(
		"shopify_temp_token",
		"",
		-1,
		"/",
		domainName,
		true,
		true,
	)

	h.logger.Info().
		Str("shop", shop).
		Str("integration_name", integrationName).
		Msg("Credenciales OAuth recuperadas exitosamente vía cookie")

	c.JSON(http.StatusOK, GetOAuthTokenResponse{
		Success:         true,
		AccessToken:     creds.AccessToken,
		ClientID:        creds.ClientID,
		ClientSecret:    creds.ClientSecret,
		Shop:            shop,
		IntegrationName: integrationName,
		IntegrationCode: integrationCode,
	})
}
