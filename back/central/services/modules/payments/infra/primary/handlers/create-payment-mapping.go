package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/domain"
)

// CreatePaymentMapping godoc
// @Summary      Crear mapeo de método de pago
// @Description  Crea un nuevo mapeo entre un método de pago original y uno unificado
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CreatePaymentMappingRequest  true  "Datos del mapeo"
// @Success      201      {object}  domain.PaymentMappingResponse
// @Failure      400      {object}  map[string]string
// @Router       /payments/mappings [post]
func (h *PaymentHandlers) CreatePaymentMapping(c *gin.Context) {
	var req domain.CreatePaymentMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.uc.CreatePaymentMapping(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
