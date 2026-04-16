package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
)

// TogglePaymentMethodActive godoc
// @Summary      Activar/desactivar método de pago
// @Description  Cambia el estado activo de un método de pago
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del método de pago"
// @Success      200  {object}  response.PaymentMethod
// @Failure      400  {object}  response.Error
// @Router       /payments/methods/{id}/toggle [patch]
func (h *PaymentHandlers) TogglePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method ID"})
		return
	}

	domainResponse, err := h.uc.TogglePaymentMethodActive(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMethodResponse(domainResponse)

	c.JSON(http.StatusOK, httpResponse)
}
