package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/request"
)

func (h *handler) PauseAI(c *gin.Context) {
	ctx := c.Request.Context()
	conversationID := c.Param("id")

	var req request.AIControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}

	businessID, ok := resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	if err := h.useCase.PauseAI(ctx, conversationID, req.PhoneNumber, businessID); err != nil {
		h.log.Error(ctx).Err(err).Str("conversation_id", conversationID).Msg("[PauseAI Handler] - error pausando AI")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "pause_failed", "message": err.Error()})
		return
	}

	h.log.Info(ctx).Str("conversation_id", conversationID).Msg("[PauseAI Handler] - AI pausado")
	c.JSON(http.StatusOK, gin.H{"status": "paused"})
}

func (h *handler) ResumeAI(c *gin.Context) {
	ctx := c.Request.Context()
	conversationID := c.Param("id")

	var req request.AIControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}

	businessID, ok := resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	if err := h.useCase.ResumeAI(ctx, conversationID, req.PhoneNumber, businessID); err != nil {
		h.log.Error(ctx).Err(err).Str("conversation_id", conversationID).Msg("[ResumeAI Handler] - error reactivando AI")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "resume_failed", "message": err.Error()})
		return
	}

	h.log.Info(ctx).Str("conversation_id", conversationID).Msg("[ResumeAI Handler] - AI reactivado")
	c.JSON(http.StatusOK, gin.H{"status": "active"})
}
