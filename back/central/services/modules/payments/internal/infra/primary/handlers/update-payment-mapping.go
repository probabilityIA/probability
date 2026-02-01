package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/request"
)

// UpdatePaymentMapping godoc
// @Summary      Actualizar mapeo de método de pago
// @Description  Actualiza un mapeo existente
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "ID del mapeo"
// @Param        request  body      request.UpdatePaymentMapping  true  "Datos actualizados"
// @Success      200      {object}  response.PaymentMapping
// @Failure      400      {object}  response.Error
// @Router       /payments/mappings/{id} [put]
func (h *PaymentHandlers) UpdatePaymentMapping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment mapping ID"})
		return
	}

	var req request.UpdatePaymentMapping
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a DTO de dominio
	dto := mappers.RequestToUpdatePaymentMappingDTO(&req)

	// Llamar caso de uso
	domainResponse, err := h.uc.UpdatePaymentMapping(c.Request.Context(), uint(id), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMappingResponse(domainResponse)

	c.JSON(http.StatusOK, httpResponse)
}
