package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// DeleteShipment godoc
// @Summary      Eliminar envío
// @Description  Elimina (soft delete) un envío del sistema
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del envío"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/{id} [delete]
func (h *Handlers) DeleteShipment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de envío inválido",
			"error":   "El ID debe ser un número válido",
		})
		return
	}

	// Llamar al caso de uso
	err = h.uc.DeleteShipment(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrShipmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Envío no encontrado",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al eliminar envío",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Envío eliminado exitosamente",
	})
}

