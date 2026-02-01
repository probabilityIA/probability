package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOAuthTokenResponse representa la respuesta del endpoint que retorna el token OAuth
type GetOAuthTokenResponse struct {
	Success         bool   `json:"success"`
	AccessToken     string `json:"access_token,omitempty"`
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

	// Validar state (debe existir en el store aunque ya fue usado)
	// En el flujo normal, el state ya fue eliminado por el callback
	// pero los datos están en los query params del redirect

	if state == "" || shop == "" {
		h.logger.Error().Msg("State o shop faltantes en petición de token")
		c.JSON(http.StatusBadRequest, GetOAuthTokenResponse{
			Success: false,
			Error:   "Parámetros requeridos faltantes",
		})
		return
	}

	// Leer cookie temporal con el token
	tokenCookie, err := c.Cookie("shopify_temp_token")
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

	// Borrar la cookie inmediatamente después de leerla (one-time use)
	c.SetCookie(
		"shopify_temp_token",
		"",
		-1, // MaxAge negativo elimina la cookie
		"/",
		".probabilityia.com.co", // Con punto inicial para subdominios y iframes
		true,  // Secure
		true,  // HttpOnly
	)

	h.logger.Info().
		Str("shop", shop).
		Str("integration_name", integrationName).
		Msg("Token OAuth recuperado exitosamente desde cookie")

	c.JSON(http.StatusOK, GetOAuthTokenResponse{
		Success:         true,
		AccessToken:     tokenCookie,
		Shop:            shop,
		IntegrationName: integrationName,
		IntegrationCode: integrationCode,
	})
}
