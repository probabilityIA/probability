package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
)

// GetPaymentMappingsByIntegrationType godoc
// @Summary      Obtener mapeos por tipo de integración
// @Description  Obtiene todos los mapeos de un tipo de integración específico
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        type  path      string  true  "Tipo de integración (shopify, whatsapp, mercadolibre)"
// @Success      200   {array}   response.PaymentMapping
// @Failure      400   {object}  response.Error
// @Router       /payments/mappings/integration/{type} [get]
func (h *PaymentHandlers) GetPaymentMappingsByIntegration(c *gin.Context) {
	integrationType := c.Param("type")
	if integrationType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration type is required"})
		return
	}

	domainResponses, err := h.uc.GetPaymentMappingsByIntegrationType(c.Request.Context(), integrationType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir cada DTO a response HTTP
	httpResponses := make([]interface{}, len(domainResponses))
	for i, dto := range domainResponses {
		httpResponses[i] = mappers.DomainToPaymentMappingResponse(&dto)
	}

	c.JSON(http.StatusOK, httpResponses)
}
