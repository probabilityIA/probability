package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListWarehouses(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := dtos.ListWarehousesParams{
		BusinessID: businessID,
		Search:     c.Query("search"),
		Page:       page,
		PageSize:   pageSize,
	}

	if v := c.Query("is_active"); v != "" {
		b, _ := strconv.ParseBool(v)
		params.IsActive = &b
	}
	if v := c.Query("is_fulfillment"); v != "" {
		b, _ := strconv.ParseBool(v)
		params.IsFulfillment = &b
	}

	warehouses, total, err := h.uc.ListWarehouses(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.WarehouseResponse, len(warehouses))
	for i, w := range warehouses {
		data[i] = response.FromEntity(&w)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.WarehouseListResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	})
}
