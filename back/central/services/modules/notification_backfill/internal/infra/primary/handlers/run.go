package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
)

func (h *Handlers) Run(c *gin.Context) {
	var req dtos.RunRequest
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
	req.BusinessID = scope

	createdBy := c.GetUint("user_id")

	job, err := h.uc.Run(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, dtos.RunResponse{JobID: job.ID})
}
