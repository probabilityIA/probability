package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/request"
)

// UpdatePaymentMethod godoc
// @Summary      Actualizar método de pago
// @Description  Actualiza un método de pago existente
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "ID del método de pago"
// @Param        request  body      request.UpdatePaymentMethod  true  "Datos actualizados"
// @Success      200      {object}  response.PaymentMethod
// @Failure      400      {object}  response.Error
// @Router       /payments/methods/{id} [put]
func (h *PaymentHandlers) UpdatePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method ID"})
		return
	}

	var req request.UpdatePaymentMethod
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a DTO de dominio
	dto := mappers.RequestToUpdatePaymentMethodDTO(&req)

	// Llamar caso de uso
	domainResponse, err := h.uc.UpdatePaymentMethod(c.Request.Context(), uint(id), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMethodResponse(domainResponse)

	c.JSON(http.StatusOK, httpResponse)
}
