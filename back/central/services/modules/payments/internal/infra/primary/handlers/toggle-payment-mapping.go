package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
)

// TogglePaymentMappingActive godoc
// @Summary      Activar/desactivar mapeo
// @Description  Cambia el estado activo de un mapeo
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del mapeo"
// @Success      200  {object}  response.PaymentMapping
// @Failure      400  {object}  response.Error
// @Router       /payments/mappings/{id}/toggle [patch]
func (h *PaymentHandlers) TogglePaymentMapping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment mapping ID"})
		return
	}

	domainResponse, err := h.uc.TogglePaymentMappingActive(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMappingResponse(domainResponse)

	c.JSON(http.StatusOK, httpResponse)
}
