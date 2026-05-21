package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListClientGroups(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	page, pageSize := parsePagination(c)
	groups, total, err := h.uc.ListClientGroups(c.Request.Context(), dtos.ListClientGroupsParams{
		BusinessID: businessID,
		Search:     c.Query("search"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	data := make([]response.ClientGroupResponse, len(groups))
	for i := range groups {
		data[i] = response.FromGroupEntity(&groups[i])
	}

	c.JSON(http.StatusOK, response.ClientGroupsListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: response.TotalPages(total, pageSize),
	})
}
