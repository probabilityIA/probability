package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
)

// ListPaymentMappings godoc
// @Summary      Listar mapeos de métodos de pago
// @Description  Obtiene una lista de todos los mapeos
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        integration_type  query     string  false  "Filtrar por tipo de integración"
// @Param        is_active         query     bool    false  "Filtrar por estado activo"
// @Success      200               {object}  response.PaymentMappingsList
// @Failure      400               {object}  response.Error
// @Router       /payments/mappings [get]
func (h *PaymentHandlers) ListPaymentMappings(c *gin.Context) {
	filters := make(map[string]interface{})
	if integrationType := c.Query("integration_type"); integrationType != "" {
		filters["integration_type"] = integrationType
	}
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"], _ = strconv.ParseBool(isActive)
	}

	domainResponse, err := h.uc.ListPaymentMappings(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMappingsListResponse(domainResponse)

	c.JSON(http.StatusOK, httpResponse)
}
