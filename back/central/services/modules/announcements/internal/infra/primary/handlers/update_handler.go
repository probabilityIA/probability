package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/request"
)

func (h *handler) Update(c *gin.Context) {
	id, ok := h.parseIDParam(c)
	if !ok {
		return
	}

	var req request.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	dto := mappers.UpdateRequestToDTO(id, req)

	result, err := h.uc.UpdateAnnouncement(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	resp := mappers.EntityToResponse(result)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}
