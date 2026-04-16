package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/request"
)

// CreatePaymentMethod godoc
// @Summary      Crear método de pago
// @Description  Crea un nuevo método de pago en el sistema
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreatePaymentMethod  true  "Datos del método de pago"
// @Success      201      {object}  response.PaymentMethod
// @Failure      400      {object}  response.Error
// @Router       /payments/methods [post]
func (h *PaymentHandlers) CreatePaymentMethod(c *gin.Context) {
	var req request.CreatePaymentMethod
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a DTO de dominio
	dto := mappers.RequestToCreatePaymentMethodDTO(&req)

	// Llamar caso de uso
	domainResponse, err := h.uc.CreatePaymentMethod(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMethodResponse(domainResponse)

	c.JSON(http.StatusCreated, httpResponse)
}
