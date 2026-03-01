package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
)

func (h *handlers) TransferStock(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.TransferStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	var createdByID *uint
	if userID > 0 {
		createdByID = &userID
	}

	dto := dtos.TransferStockDTO{
		ProductID:       req.ProductID,
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		FromLocationID:  req.FromLocationID,
		ToLocationID:    req.ToLocationID,
		BusinessID:      businessID,
		Quantity:        req.Quantity,
		Reason:          req.Reason,
		Notes:           req.Notes,
		CreatedByID:     createdByID,
	}

	if err := h.uc.TransferStock(c.Request.Context(), dto); err != nil {
		if errors.Is(err, domainerrors.ErrProductNotFound) || errors.Is(err, domainerrors.ErrWarehouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrInsufficientStock) || errors.Is(err, domainerrors.ErrSameWarehouse) ||
			errors.Is(err, domainerrors.ErrTransferQtyNeg) || errors.Is(err, domainerrors.ErrProductNoTracking) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock transferred successfully"})
}
