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
	// Verificar si hay filtro por notification_type_id
	notificationTypeIDStr := c.Query("notification_type_id")

	if notificationTypeIDStr != "" {
		// Filtrar por tipo de notificación
		notificationTypeID, err := strconv.ParseUint(notificationTypeIDStr, 10, 32)
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid notification_type_id parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification_type_id"})
			return
		}

		events, err := h.useCase.GetEventTypesByNotificationType(c.Request.Context(), uint(notificationTypeID))
		if err != nil {
			h.logger.Error().Err(err).Uint64("notification_type_id", notificationTypeID).Msg("Error listing event types by notification type")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := mappers.DomainListToResponse(events)
		c.JSON(http.StatusOK, response)
		return
	}

	// Sin filtro - retornar lista vacía o implementar GetAll si se necesita
	c.JSON(http.StatusOK, []interface{}{})
}
