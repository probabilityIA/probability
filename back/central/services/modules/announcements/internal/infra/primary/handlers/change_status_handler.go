package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/request"
)

func (h *handler) ChangeStatus(c *gin.Context) {
	id, ok := h.parseIDParam(c)
	if !ok {
		return
	}

	var req request.ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	dto := dtos.ChangeStatusDTO{
		ID:     id,
		Status: entities.AnnouncementStatus(req.Status),
	}

	if err := h.uc.ChangeStatus(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "status updated"})
}
