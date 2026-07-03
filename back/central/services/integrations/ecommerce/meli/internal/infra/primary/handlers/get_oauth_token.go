package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetOAuthTokenResponse struct {
	Success         bool   `json:"success"`
	AccessToken     string `json:"access_token,omitempty"`
	RefreshToken    string `json:"refresh_token,omitempty"`
	ClientID        string `json:"client_id,omitempty"`
	ClientSecret    string `json:"client_secret,omitempty"`
	SellerID        int64  `json:"seller_id,omitempty"`
	IsTesting       bool   `json:"is_testing"`
	ExpiresAt       string `json:"expires_at,omitempty"`
	IntegrationName string `json:"integration_name,omitempty"`
	IntegrationCode string `json:"integration_code,omitempty"`
	Error           string `json:"error,omitempty"`
}

func (h *meliHandler) GetOAuthTokenHandler(c *gin.Context) {
	state := c.Query("state")
	exchangeToken := c.Query("exchange_token")
	integrationName := c.Query("integration_name")
	integrationCode := c.Query("integration_code")

	if state == "" || exchangeToken == "" {
		c.JSON(http.StatusBadRequest, GetOAuthTokenResponse{
			Success: false,
			Error:   "Parametros requeridos faltantes",
		})
		return
	}

	data, ok := retrieveExchangeToken(exchangeToken)
	if !ok {
		c.JSON(http.StatusGone, GetOAuthTokenResponse{
			Success: false,
			Error:   "El token de autorizacion expiro o ya fue consumido. Inicia la conexion con MercadoLibre nuevamente.",
		})
		return
	}

	c.JSON(http.StatusOK, GetOAuthTokenResponse{
		Success:         true,
		AccessToken:     data.AccessToken,
		RefreshToken:    data.RefreshToken,
		ClientID:        data.ClientID,
		ClientSecret:    data.ClientSecret,
		SellerID:        data.SellerID,
		IsTesting:       data.IsTesting,
		ExpiresAt:       data.ExpiresAt.Format(time.RFC3339),
		IntegrationName: integrationName,
		IntegrationCode: integrationCode,
	})
}
