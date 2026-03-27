package message_audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/request"
)

// ListConversations godoc
// @Summary Listar conversaciones de WhatsApp
// @Description Obtiene conversaciones agrupadas con resumen para la vista de chat
// @Tags message-audit
// @Produce json
// @Param business_id query uint true "ID del negocio"
// @Param state query string false "Estado de la conversación"
// @Param phone query string false "Búsqueda por teléfono"
// @Param date_from query string false "Fecha inicio (YYYY-MM-DD)"
// @Param date_to query string false "Fecha fin (YYYY-MM-DD)"
// @Param page query int false "Página" default(1)
// @Param page_size query int false "Tamaño de página" default(20)
// @Success 200 {object} response.PaginatedConversationListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/message-audit/conversations [get]
func (h *handler) ListConversations(c *gin.Context) {
	h.logger.Info().Msg("[GET /notification-configs/message-audit/conversations] Request received")

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var query request.ListConversations
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error().Err(err).Msg("Invalid query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := mappers.ConversationListRequestToDomain(&query, businessID)

	result, err := h.useCase.ListConversations(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error listing conversations")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().
		Int64("total", result.Total).
		Int("page", result.Page).
		Msg("Conversations fetched successfully")

	c.JSON(http.StatusOK, mappers.DomainToConversationListResponse(result))
}
