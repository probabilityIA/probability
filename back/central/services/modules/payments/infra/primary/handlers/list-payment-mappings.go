package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListPaymentMappings godoc
// @Summary      Listar mapeos de métodos de pago
// @Description  Obtiene una lista de mapeos de métodos de pago con filtros opcionales
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        integration_type  query    string  false  "Filtrar por tipo de integración (shopify, whatsapp, mercadolibre)"
// @Success      200  {object}  domain.PaymentMappingsListResponse
// @Failure      500  {object}  map[string]string
// @Router       /payments/mappings [get]
func (h *PaymentHandlers) ListPaymentMappings(c *gin.Context) {
	filters := make(map[string]interface{})
	if integrationType := c.Query("integration_type"); integrationType != "" {
		filters["integration_type"] = integrationType
	}

	response, err := h.uc.ListPaymentMappings(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
