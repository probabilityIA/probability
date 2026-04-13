package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/request"
)

func (h *handler) Create(c *gin.Context) {
	var req request.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	dto := mappers.CreateRequestToDTO(req, userID)

	result, err := h.uc.CreateAnnouncement(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	resp := mappers.EntityToResponse(result)
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}
