package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
)

type previewRequest struct {
	EventCode  string `json:"event_code" binding:"required"`
	BusinessID *uint  `json:"business_id"`
	Days       int    `json:"days"`
	Limit      int    `json:"limit"`
}

func (h *Handlers) Preview(c *gin.Context) {
	var req previewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scope, err := h.resolveScope(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id inválido"})
		return
	}
	if scope == nil && req.BusinessID != nil && *req.BusinessID > 0 {
		scope = req.BusinessID
	}

	filter := dtos.BackfillFilter{
		EventCode:  req.EventCode,
		BusinessID: scope,
		Days:       req.Days,
		Limit:      req.Limit,
	}

	resp, err := h.uc.Preview(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
