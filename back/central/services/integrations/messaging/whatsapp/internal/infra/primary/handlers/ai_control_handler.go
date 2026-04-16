package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/request"
)

// PauseAI pausa el bot AI para una conversación. El humano toma el control.
// POST /whatsapp/conversations/:id/pause-ai
func (h *handler) PauseAI(c *gin.Context) {
	ctx := c.Request.Context()
	conversationID := c.Param("id")

	var req request.AIControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}

	if err := h.useCase.PauseAI(ctx, conversationID, req.PhoneNumber, req.BusinessID); err != nil {
		h.log.Error(ctx).Err(err).Str("conversation_id", conversationID).Msg("[PauseAI Handler] - error pausando AI")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "pause_failed", "message": err.Error()})
		return
	}

	h.log.Info(ctx).Str("conversation_id", conversationID).Msg("[PauseAI Handler] - AI pausado")
	c.JSON(http.StatusOK, gin.H{"status": "paused"})
}

// ResumeAI reactiva el bot AI para una conversación.
// POST /whatsapp/conversations/:id/resume-ai
func (h *handler) ResumeAI(c *gin.Context) {
	ctx := c.Request.Context()
	conversationID := c.Param("id")

	var req request.AIControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}

	if err := h.useCase.ResumeAI(ctx, conversationID, req.PhoneNumber, req.BusinessID); err != nil {
		h.log.Error(ctx).Err(err).Str("conversation_id", conversationID).Msg("[ResumeAI Handler] - error reactivando AI")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "resume_failed", "message": err.Error()})
		return
	}

	h.log.Info(ctx).Str("conversation_id", conversationID).Msg("[ResumeAI Handler] - AI reactivado")
	c.JSON(http.StatusOK, gin.H{"status": "active"})
}
