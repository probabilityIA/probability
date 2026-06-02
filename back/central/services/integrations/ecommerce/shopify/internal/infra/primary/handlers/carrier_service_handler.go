package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ShopifyHandler) resolvePublicBaseURL() string {
	baseURL := h.config.Get("WEBHOOK_BASE_URL")
	if baseURL == "" {
		baseURL = h.config.Get("URL_BASE_SWAGGER")
	}
	return baseURL
}

func (h *ShopifyHandler) EnableCarrierServiceHandler(c *gin.Context) {
	integrationID := c.Param("integration_id")
	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "integration_id es requerido"})
		return
	}

	baseURL := h.resolvePublicBaseURL()
	carrierServiceID, err := h.useCase.EnableCarrierCalculatedShipping(c.Request.Context(), integrationID, baseURL)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Str("integration_id", integrationID).Msg("Error al activar carrier service")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"enabled": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":            true,
		"message":            "Cotizacion en checkout activada correctamente",
		"enabled":            true,
		"carrier_service_id": carrierServiceID,
	})
}

func (h *ShopifyHandler) DisableCarrierServiceHandler(c *gin.Context) {
	integrationID := c.Param("integration_id")
	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "integration_id es requerido"})
		return
	}

	if err := h.useCase.DisableCarrierCalculatedShipping(c.Request.Context(), integrationID); err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Str("integration_id", integrationID).Msg("Error al desactivar carrier service")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"enabled": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cotizacion en checkout desactivada correctamente",
		"enabled": false,
	})
}
