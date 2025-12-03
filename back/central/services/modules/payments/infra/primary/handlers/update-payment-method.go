package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/domain"
)

// UpdatePaymentMethod godoc
// @Summary      Actualizar método de pago
// @Description  Actualiza un método de pago existente
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        id       path      int                                 true  "ID del método de pago"
// @Param        request  body      domain.UpdatePaymentMethodRequest  true  "Datos a actualizar"
// @Success      200      {object}  domain.PaymentMethodResponse
// @Failure      400      {object}  map[string]string
// @Router       /payments/methods/{id} [put]
func (h *PaymentHandlers) UpdatePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.uc.UpdatePaymentMethod(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
