package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListInvoices(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	params := dtos.ListInvoicesParams{
		Status:   c.Query("status"),
		Page:     page,
		PageSize: pageSize,
	}
	if v := c.Query("is_test"); v == "true" || v == "false" {
		isTest := v == "true"
		params.IsTest = &isTest
	}
	if v := c.Query("business_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil && id > 0 {
			businessID := uint(id)
			params.BusinessID = &businessID
		}
	}
	items, total, err := h.uc.ListInvoices(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	data := make([]response.InvoiceResponse, len(items))
	for i := range items {
		data[i] = response.FromInvoice(&items[i])
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, response.InvoicesListResponse{
		Success:    true,
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
