package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/request"
)

func (h *handler) RegisterView(c *gin.Context) {
	id, ok := h.parseIDParam(c)
	if !ok {
		return
	}

	var req request.RegisterViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	dto := dtos.RegisterViewDTO{
		AnnouncementID: id,
		UserID:         c.GetUint("user_id"),
		BusinessID:     c.GetUint("business_id"),
		Action:         entities.ViewAction(req.Action),
		LinkID:         req.LinkID,
	}

	if err := h.uc.RegisterView(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "view registered"})
}
