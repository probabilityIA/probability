package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de orden inválido",
			"error":   "El ID de la orden es requerido",
		})
		return
	}

	if !h.ensureOrderOwnership(c, id) {
		return
	}

	err := h.orderCRUD.DeleteOrder(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Orden no encontrada",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al eliminar orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Orden eliminada exitosamente",
	})
}
