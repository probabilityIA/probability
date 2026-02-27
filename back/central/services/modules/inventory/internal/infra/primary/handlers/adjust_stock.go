package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
)

func (h *Handlers) AdjustStock(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener user_id del JWT
	userID := c.GetUint("user_id")
	var createdByID *uint
	if userID > 0 {
		createdByID = &userID
	}

	dto := dtos.AdjustStockDTO{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		LocationID:  req.LocationID,
		BusinessID:  businessID,
		Quantity:    req.Quantity,
		Reason:      req.Reason,
		Notes:       req.Notes,
		CreatedByID: createdByID,
	}

	movement, err := h.uc.AdjustStock(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrProductNotFound) || errors.Is(err, domainerrors.ErrWarehouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrInsufficientStock) || errors.Is(err, domainerrors.ErrInvalidQuantity) || errors.Is(err, domainerrors.ErrProductNoTracking) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.StockMovementFromEntity(movement))
}
