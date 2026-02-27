package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

func (h *handler) CreateOrderStatus(c *gin.Context) {
	var req request.CreateOrderStatusRequest
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

	created, err := h.uc.CreateOrderStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear estado de orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Estado de orden creado exitosamente",
		"data": response.OrderStatusCatalogResponse{
			ID:          created.ID,
			Code:        created.Code,
			Name:        created.Name,
			Description: created.Description,
			Category:    created.Category,
			Color:       created.Color,
			Priority:    created.Priority,
			IsActive:    created.IsActive,
		},
	})
}
