package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
)

func (h *handler) DeleteChannelStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID inválido",
		})
		return
	}

	if err := h.uc.DeleteChannelStatus(c.Request.Context(), uint(id)); err != nil {
		if err == domainerrors.ErrChannelStatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Estado del canal no encontrado",
			})
			return
		}
		if err == domainerrors.ErrChannelStatusHasMappings {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "No se puede eliminar el estado porque tiene mapeos asociados",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al eliminar estado del canal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Estado del canal eliminado exitosamente",
	})
}
