package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
)

func (h *handler) DeleteOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID inválido",
		})
		return
	}

	if err := h.uc.DeleteOrderStatus(c.Request.Context(), uint(id)); err != nil {
		if err == domainerrors.ErrOrderStatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Estado de orden no encontrado",
			})
			return
		}
		if err == domainerrors.ErrOrderStatusHasMappings {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "No se puede eliminar el estado porque tiene mapeos asociados",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al eliminar estado de orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Estado de orden eliminado exitosamente",
	})
}
