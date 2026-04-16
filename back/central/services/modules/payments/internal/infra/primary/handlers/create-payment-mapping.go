package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/request"
)

// CreatePaymentMapping godoc
// @Summary      Crear mapeo de método de pago
// @Description  Crea un nuevo mapeo entre un método de pago externo y uno interno
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreatePaymentMapping  true  "Datos del mapeo"
// @Success      201      {object}  response.PaymentMapping
// @Failure      400      {object}  response.Error
// @Router       /payments/mappings [post]
func (h *PaymentHandlers) CreatePaymentMapping(c *gin.Context) {
	var req request.CreatePaymentMapping
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a DTO de dominio
	dto := mappers.RequestToCreatePaymentMappingDTO(&req)

	// Llamar caso de uso
	domainResponse, err := h.uc.CreatePaymentMapping(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMappingResponse(domainResponse)

	c.JSON(http.StatusCreated, httpResponse)
}
