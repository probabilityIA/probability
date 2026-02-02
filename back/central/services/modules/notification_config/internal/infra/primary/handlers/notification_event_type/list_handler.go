package notification_event_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/mappers"
)

// List godoc
// @Summary Listar tipos de evento de notificación
// @Description Obtiene todos los tipos de evento, opcionalmente filtrados por tipo de notificación
// @Tags notification-event-types
// @Accept json
// @Produce json
// @Param notification_type_id query int false "Filtrar por ID del tipo de notificación"
// @Success 200 {array} response.NotificationEventType
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-event-types [get]
func (h *handler) List(c *gin.Context) {
	h.logger.Info().Msg("🌐 [GET /notification-event-types] Request received")

	// Verificar si hay filtro por notification_type_id
	notificationTypeIDStr := c.Query("notification_type_id")
	h.logger.Info().Str("notification_type_id_param", notificationTypeIDStr).Msg("📋 Query params")

	if notificationTypeIDStr != "" {
		// Filtrar por tipo de notificación
		notificationTypeID, err := strconv.ParseUint(notificationTypeIDStr, 10, 32)
		if err != nil {
			h.logger.Error().Err(err).Msg("❌ Invalid notification_type_id parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification_type_id"})
			return
		}

		h.logger.Info().Uint64("notification_type_id", notificationTypeID).Msg("🔍 Fetching event types by notification type")

		events, err := h.useCase.GetEventTypesByNotificationType(c.Request.Context(), uint(notificationTypeID))
		if err != nil {
			h.logger.Error().Err(err).Uint64("notification_type_id", notificationTypeID).Msg("❌ Error listing event types by notification type")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		h.logger.Info().Int("count", len(events)).Msg("✅ Event types fetched successfully")

		response := mappers.DomainListToResponse(events)
		c.JSON(http.StatusOK, response)
		return
	}

	// Sin filtro - retornar todos los event types
	h.logger.Info().Msg("🔍 Fetching all event types (no filter)")

	// Llamar al use case para obtener todos
	events, err := h.useCase.ListAllEventTypes(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("❌ Error listing all event types")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Int("count", len(events)).Msg("✅ All event types fetched successfully")

	response := mappers.DomainListToResponse(events)
	c.JSON(http.StatusOK, response)
}
