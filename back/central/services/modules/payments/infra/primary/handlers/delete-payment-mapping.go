package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeletePaymentMapping godoc
// @Summary      Eliminar mapeo de método de pago
// @Description  Elimina un mapeo de método de pago
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "ID del mapeo"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Router       /payments/mappings/{id} [delete]
func (h *PaymentHandlers) DeletePaymentMapping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.uc.DeletePaymentMapping(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
