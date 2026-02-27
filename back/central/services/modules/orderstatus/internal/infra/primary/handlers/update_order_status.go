package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

func (h *handler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID inválido",
		})
		return
	}

	var req request.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos inválidos",
			"error":   err.Error(),
		})
		return
	}

	status := &entities.OrderStatusInfo{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Color:       req.Color,
		Priority:    req.Priority,
		IsActive:    req.IsActive,
	}

	updated, err := h.uc.UpdateOrderStatus(c.Request.Context(), uint(id), status)
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
			"message": "Error al actualizar estado de orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Estado de orden actualizado exitosamente",
		"data": response.OrderStatusCatalogResponse{
			ID:          updated.ID,
			Code:        updated.Code,
			Name:        updated.Name,
			Description: updated.Description,
			Category:    updated.Category,
			Color:       updated.Color,
			Priority:    updated.Priority,
			IsActive:    updated.IsActive,
		},
	})
}
