package message_audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/mappers"
)

// GetConversationMessages godoc
// @Summary Obtener mensajes de una conversación
// @Description Obtiene los mensajes de una conversación específica para la vista de chat
// @Tags message-audit
// @Produce json
// @Param id path string true "ID de la conversación (UUID)"
// @Param business_id query uint true "ID del negocio"
// @Success 200 {object} response.ConversationDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/message-audit/conversations/{id}/messages [get]
func (h *handler) GetConversationMessages(c *gin.Context) {
	h.logger.Info().Msg("[GET /notification-configs/message-audit/conversations/:id/messages] Request received")

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	conversationID := c.Param("id")
	if _, err := uuid.Parse(conversationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation ID"})
		return
	}

	result, err := h.useCase.GetConversationMessages(c.Request.Context(), conversationID, businessID)
	if err != nil {
		h.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Error getting conversation messages")
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	h.logger.Info().
		Str("conversation_id", conversationID).
		Int("message_count", len(result.Messages)).
		Msg("Conversation messages fetched successfully")

	c.JSON(http.StatusOK, mappers.DomainToConversationDetailResponse(result))
}
