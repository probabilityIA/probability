package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListClientPricingRules(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := dtos.ListClientPricingRulesParams{
		BusinessID: businessID,
		Page:       page,
		PageSize:   pageSize,
	}

	if clientIDStr := c.Query("client_id"); clientIDStr != "" {
		if id, err := strconv.ParseUint(clientIDStr, 10, 64); err == nil {
			clientID := uint(id)
			params.ClientID = &clientID
		}
	}

	if productID := c.Query("product_id"); productID != "" {
		params.ProductID = &productID
	}

	rules, total, err := h.uc.ListClientPricingRules(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.ClientPricingRuleResponse, len(rules))
	for i, r := range rules {
		data[i] = response.FromRuleEntity(&r)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.PricingRulesListResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	})
}
