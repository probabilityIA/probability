package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPaymentMappingsByIntegration godoc
// @Summary      Obtener mapeos por tipo de integración
// @Description  Obtiene todos los mapeos de un tipo de integración específico
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        type  path      string  true  "Tipo de integración (shopify, whatsapp, mercadolibre)"
// @Success      200   {array}   domain.PaymentMappingResponse
// @Failure      500   {object}  map[string]string
// @Router       /payments/mappings/integration/{type} [get]
func (h *PaymentHandlers) GetPaymentMappingsByIntegration(c *gin.Context) {
	integrationType := c.Param("type")

	response, err := h.uc.GetPaymentMappingsByIntegrationType(c.Request.Context(), integrationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
