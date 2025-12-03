package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/domain"
)

// UpdatePaymentMapping godoc
// @Summary      Actualizar mapeo de método de pago
// @Description  Actualiza un mapeo existente
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        id       path      int                                 true  "ID del mapeo"
// @Param        request  body      domain.UpdatePaymentMappingRequest  true  "Datos a actualizar"
// @Success      200      {object}  domain.PaymentMappingResponse
// @Failure      400      {object}  map[string]string
// @Router       /payments/mappings/{id} [put]
func (h *PaymentHandlers) UpdatePaymentMapping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdatePaymentMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.uc.UpdatePaymentMapping(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
