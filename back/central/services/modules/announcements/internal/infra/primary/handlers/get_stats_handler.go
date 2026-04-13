package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/mappers"
)

func (h *handler) GetStats(c *gin.Context) {
	id, ok := h.parseIDParam(c)
	if !ok {
		return
	}

	stats, err := h.uc.GetAnnouncementStats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	resp := mappers.StatsToResponse(stats)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}
