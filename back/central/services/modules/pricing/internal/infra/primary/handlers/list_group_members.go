package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListGroupMembers(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	groupID, ok := parseUintParam(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	page, pageSize := parsePagination(c)
	members, total, err := h.uc.ListGroupMembers(c.Request.Context(), dtos.ListGroupMembersParams{
		BusinessID:    businessID,
		ClientGroupID: groupID,
		Search:        c.Query("search"),
		Page:          page,
		PageSize:      pageSize,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	data := make([]response.ClientSummaryResponse, len(members))
	for i := range members {
		data[i] = response.FromClientSummary(&members[i])
	}

	c.JSON(http.StatusOK, response.ClientsListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: response.TotalPages(total, pageSize),
	})
}
