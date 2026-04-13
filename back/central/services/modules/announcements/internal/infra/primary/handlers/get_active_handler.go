package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/mappers"
)

func (h *handler) GetActive(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")

	params := dtos.ActiveAnnouncementsParams{
		BusinessID: businessID,
		UserID:     userID,
	}

	items, err := h.uc.GetActiveAnnouncements(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	data := make([]interface{}, 0, len(items))
	for _, item := range items {
		data = append(data, mappers.EntityToResponse(&item))
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
