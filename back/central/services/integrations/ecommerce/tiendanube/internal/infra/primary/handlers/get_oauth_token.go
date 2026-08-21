package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getOAuthTokenResponse struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token,omitempty"`
	StoreID     string `json:"store_id,omitempty"`
	Scope       string `json:"scope,omitempty"`
	IsTesting   bool   `json:"is_testing"`
	Error       string `json:"error,omitempty"`
}

func (h *tiendanubeHandler) GetOAuthToken(c *gin.Context) {
	exchangeToken := c.Query("exchange_token")
	if exchangeToken == "" {
		c.JSON(http.StatusBadRequest, getOAuthTokenResponse{Success: false, Error: "Parametros requeridos faltantes"})
		return
	}

	data, ok := retrieveExchangeToken(exchangeToken)
	if !ok {
		c.JSON(http.StatusGone, getOAuthTokenResponse{
			Success: false,
			Error:   "El token de autorizacion expiro o ya fue consumido. Inicia la conexion con Tiendanube nuevamente.",
		})
		return
	}

	c.JSON(http.StatusOK, getOAuthTokenResponse{
		Success:     true,
		AccessToken: data.AccessToken,
		StoreID:     data.StoreID,
		Scope:       data.Scope,
		IsTesting:   data.IsTesting,
	})
}
