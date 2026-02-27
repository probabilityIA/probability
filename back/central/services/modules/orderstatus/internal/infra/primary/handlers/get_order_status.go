package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

func (h *handler) GetOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID inválido",
		})
		return
	}

	status, err := h.uc.GetOrderStatus(c.Request.Context(), uint(id))
	if err != nil {
		if err == domainerrors.ErrOrderStatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Estado de orden no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener estado de orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Estado de orden obtenido exitosamente",
		"data": response.OrderStatusCatalogResponse{
			ID:          status.ID,
			Code:        status.Code,
			Name:        status.Name,
			Description: status.Description,
			Category:    status.Category,
			Color:       status.Color,
			Priority:    status.Priority,
			IsActive:    status.IsActive,
		},
	})
}
