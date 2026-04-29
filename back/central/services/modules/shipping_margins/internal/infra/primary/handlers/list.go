package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/infra/primary/handlers/response"
)

func (h *Handlers) List(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := dtos.ListShippingMarginsParams{
		BusinessID:  businessID,
		CarrierCode: c.Query("carrier_code"),
		Page:        page,
		PageSize:    pageSize,
	}

	items, total, err := h.uc.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.ShippingMarginResponse, len(items))
	for i, m := range items {
		data[i] = response.FromEntity(&m)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.ShippingMarginsListResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	})
}
