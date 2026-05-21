package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/response"
)

func parseOptionalUint(c *gin.Context, key string) *uint {
	raw := c.Query(key)
	if raw == "" {
		return nil
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		return nil
	}
	id := uint(value)
	return &id
}

func (h *Handlers) GetCatalogPrices(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	page, pageSize := parsePagination(c)
	rows, total, err := h.uc.ListCatalogPrices(c.Request.Context(), dtos.ListCatalogPricesParams{
		Target: dtos.CatalogPriceTarget{
			BusinessID:    businessID,
			ClientGroupID: parseOptionalUint(c, "client_group_id"),
			ClientID:      parseOptionalUint(c, "client_id"),
		},
		Search:   c.Query("search"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	data := make([]response.CatalogPriceRowResponse, len(rows))
	for i := range rows {
		data[i] = response.FromCatalogPriceRow(&rows[i])
	}

	c.JSON(http.StatusOK, response.CatalogPricesListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: response.TotalPages(total, pageSize),
	})
}
