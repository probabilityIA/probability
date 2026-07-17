package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type getOAuthTokenResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IsTesting    bool   `json:"is_testing"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	Error        string `json:"error,omitempty"`
}

func (h *jumpsellerHandler) GetOAuthToken(c *gin.Context) {
	exchangeToken := c.Query("exchange_token")
	if exchangeToken == "" {
		c.JSON(http.StatusBadRequest, getOAuthTokenResponse{Success: false, Error: "Parametros requeridos faltantes"})
		return
	}

	data, ok := retrieveExchangeToken(exchangeToken)
	if !ok {
		c.JSON(http.StatusGone, getOAuthTokenResponse{
			Success: false,
			Error:   "El token de autorizacion expiro o ya fue consumido. Inicia la conexion con Jumpseller nuevamente.",
		})
		return
	}

	c.JSON(http.StatusOK, getOAuthTokenResponse{
		Success:      true,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		IsTesting:    data.IsTesting,
		ExpiresAt:    data.ExpiresAt.Format(time.RFC3339),
	})
}
