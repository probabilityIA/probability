package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/domain"
)

// CreatePaymentMethod godoc
// @Summary      Crear método de pago
// @Description  Crea un nuevo método de pago en el sistema
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CreatePaymentMethodRequest  true  "Datos del método de pago"
// @Success      201      {object}  domain.PaymentMethodResponse
// @Failure      400      {object}  map[string]string
// @Router       /payments/methods [post]
func (h *PaymentHandlers) CreatePaymentMethod(c *gin.Context) {
	var req domain.CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.uc.CreatePaymentMethod(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
