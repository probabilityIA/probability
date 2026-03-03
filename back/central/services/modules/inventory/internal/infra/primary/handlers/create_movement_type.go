package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
)

func (h *handlers) CreateMovementType(c *gin.Context) {
	_, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.CreateMovementTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dtos.CreateStockMovementTypeDTO{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Direction:   req.Direction,
	}

	movType, err := h.uc.CreateMovementType(c.Request.Context(), dto)
	if err != nil {
		if err == domainerrors.ErrMovementTypeCodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.MovementTypeFromEntity(movType))
}
