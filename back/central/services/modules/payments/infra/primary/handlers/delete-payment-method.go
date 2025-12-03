package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeletePaymentMethod godoc
// @Summary      Eliminar método de pago
// @Description  Elimina un método de pago (no se puede eliminar si tiene mapeos activos)
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "ID del método de pago"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Router       /payments/methods/{id} [delete]
func (h *PaymentHandlers) DeletePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.uc.DeletePaymentMethod(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
