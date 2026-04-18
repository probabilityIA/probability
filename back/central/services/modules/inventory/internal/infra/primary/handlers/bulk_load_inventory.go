package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apprequest "github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (h *handlers) BulkLoadInventory(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.BulkLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": friendlyValidationError(err)})
		return
	}

	userID := c.GetUint("user_id")
	var createdByID *uint
	if userID > 0 {
		createdByID = &userID
	}

	dto := apprequest.BulkLoadDTO{
		WarehouseID: req.WarehouseID,
		BusinessID:  businessID,
		CreatedByID: createdByID,
		Reason:      req.Reason,
		Items:       make([]apprequest.BulkLoadItem, len(req.Items)),
	}

	for i, item := range req.Items {
		dto.Items[i] = apprequest.BulkLoadItem{
			SKU:          item.SKU,
			Quantity:     item.Quantity,
			MinStock:     item.MinStock,
			MaxStock:     item.MaxStock,
			ReorderPoint: item.ReorderPoint,
		}
	}

	// Si hay RabbitMQ disponible, publicar async
	if h.rabbit != nil {
		msgBytes, err := json.Marshal(dto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize request"})
			return
		}

		if err := h.rabbit.Publish(c.Request.Context(), rabbitmq.QueueInventoryBulkLoad, msgBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue bulk load"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message":     "Bulk load enqueued for processing",
			"total_items": len(dto.Items),
		})
		return
	}

	// Fallback: procesamiento síncrono si no hay RabbitMQ
	result, err := h.uc.BulkLoadInventory(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrWarehouseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrInvalidQuantity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.BulkLoadResultFromDTO(result))
}
