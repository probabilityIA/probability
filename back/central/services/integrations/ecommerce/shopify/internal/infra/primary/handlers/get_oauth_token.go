package handlers

import (
	"encoding/json"
	"net/http"
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

	if state == "" || shop == "" {
		h.logger.Error().Msg("State o shop faltantes en petición de token")
		c.JSON(http.StatusBadRequest, GetOAuthTokenResponse{
			Success: false,
			Error:   "Parámetros requeridos faltantes",
		})
		return
	}

	// Leer cookie temporal con el token o JSON de credenciales
	cookieData, err := c.Cookie("shopify_temp_token")
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("state", state).
			Msg("Cookie de token no encontrada o expirada")
		c.JSON(http.StatusNotFound, GetOAuthTokenResponse{
			Success: false,
			Error:   "Token no encontrado o expirado",
		})
		return
	}

	// Decodificar JSON de la cookie (si es JSON)
	var creds struct {
		AccessToken  string `json:"access_token"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	if strings.HasPrefix(cookieData, "{") {
		if err := json.Unmarshal([]byte(cookieData), &creds); err != nil {
			h.logger.Error().Err(err).Msg("Error al decodificar JSON de credenciales")
			// Fallback: tratar como token simple si falla el unmarshal
			creds.AccessToken = cookieData
		}
	} else {
		// Compatible con versiones anteriores que guardaban solo el token
		creds.AccessToken = cookieData
	}

	// Borrar la cookie inmediatamente después de leerla (one-time use)
	host := c.Request.Host
	domainName := ".probabilityia.com.co"
	if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
		domainName = ""
	}

	c.SetCookie(
		"shopify_temp_token",
		"",
		-1, // MaxAge negativo elimina la cookie
		"/",
		domainName,
		true, // Secure
		true, // HttpOnly
	)

	h.logger.Info().
		Str("shop", shop).
		Str("integration_name", integrationName).
		Msg("Credenciales OAuth recuperadas exitosamente")

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
