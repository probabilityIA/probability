package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListAvailableClients(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	page, pageSize := parsePagination(c)
	clients, total, err := h.uc.ListAvailableClients(c.Request.Context(), dtos.ListAvailableClientsParams{
		BusinessID:    businessID,
		Search:        c.Query("search"),
		OnlyUngrouped: c.Query("only_ungrouped") == "true",
		Page:          page,
		PageSize:      pageSize,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	data := make([]response.ClientSummaryResponse, len(clients))
	for i := range clients {
		data[i] = response.FromClientSummary(&clients[i])
	}

	c.JSON(http.StatusOK, response.ClientsListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: response.TotalPages(total, pageSize),
	})
}
