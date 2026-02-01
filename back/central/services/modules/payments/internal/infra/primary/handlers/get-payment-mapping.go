package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
)

// GetPaymentMapping godoc
// @Summary      Obtener mapeo por ID
// @Description  Obtiene un mapeo específico por su ID
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del mapeo"
// @Success      200  {object}  response.PaymentMapping
// @Failure      404  {object}  response.Error
// @Router       /payments/mappings/{id} [get]
func (h *PaymentHandlers) GetPaymentMapping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment mapping ID"})
		return
	}

	domainResponse, err := h.uc.GetPaymentMappingByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMappingResponse(domainResponse)

	c.JSON(http.StatusOK, httpResponse)
}
