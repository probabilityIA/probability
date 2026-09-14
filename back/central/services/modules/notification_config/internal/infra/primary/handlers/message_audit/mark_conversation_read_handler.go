package message_audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) MarkConversationRead(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id is required"})
		return
	}

	conversationID := c.Param("id")
	if _, err := uuid.Parse(conversationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid conversation ID"})
		return
	}

	var userID *uint
	if id := c.GetUint("user_id"); id > 0 {
		userID = &id
	}

	if err := h.useCase.MarkConversationRead(c.Request.Context(), conversationID, businessID, userID); err != nil {
		h.logger.Warn().Err(err).Str("conversation_id", conversationID).Msg("No se pudo marcar la conversacion como leida")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Conversation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
