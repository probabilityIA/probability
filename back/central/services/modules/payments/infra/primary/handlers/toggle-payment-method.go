package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TogglePaymentMethod godoc
// @Summary      Activar/Desactivar método de pago
// @Description  Cambia el estado activo/inactivo de un método de pago
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del método de pago"
// @Success      200  {object}  domain.PaymentMethodResponse
// @Failure      400  {object}  map[string]string
// @Router       /payments/methods/{id}/toggle [patch]
func (h *PaymentHandlers) TogglePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	response, err := h.uc.TogglePaymentMethodActive(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
