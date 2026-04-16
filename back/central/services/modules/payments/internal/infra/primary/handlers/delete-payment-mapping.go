package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeletePaymentMapping godoc
// @Summary      Eliminar mapeo
// @Description  Elimina un mapeo del sistema
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "ID del mapeo"
// @Success      204  "No Content"
// @Failure      400  {object}  response.Error
// @Router       /payments/mappings/{id} [delete]
func (h *PaymentHandlers) DeletePaymentMapping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment mapping ID"})
		return
	}

	if err := h.uc.DeletePaymentMapping(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
