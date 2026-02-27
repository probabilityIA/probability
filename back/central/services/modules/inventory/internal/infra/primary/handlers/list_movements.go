package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListMovements(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := dtos.ListMovementsParams{
		BusinessID: businessID,
		ProductID:  c.Query("product_id"),
		Type:       c.Query("type"),
		Page:       page,
		PageSize:   pageSize,
	}

	if v := c.Query("warehouse_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			whID := uint(id)
			params.WarehouseID = &whID
		}
	}

	movements, total, err := h.uc.ListMovements(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.StockMovementResponse, len(movements))
	for i, m := range movements {
		data[i] = response.StockMovementFromEntity(&m)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.MovementListResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	})
}
