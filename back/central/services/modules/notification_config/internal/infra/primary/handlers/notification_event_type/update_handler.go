package notification_event_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/request"
)

// Update godoc
// @Summary Actualizar tipo de evento
// @Description Actualiza un tipo de evento de notificación existente
// @Tags notification-event-types
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de evento"
// @Param body body request.UpdateNotificationEventType true "Datos a actualizar"
// @Success 200 {object} response.NotificationEventType
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-event-types/{id} [put]
func (h *handler) Update(c *gin.Context) {
	h.logger.Info().Msg("🌐 [PUT /notification-event-types/:id] Request received")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("📋 Path parameter parsed")

	var req request.UpdateNotificationEventType
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Interface("request_body", req).
		Msg("📋 Request body parsed")

	// Obtener la entidad existente
	h.logger.Info().Uint64("id", id).Msg("🔍 Fetching existing notification event type")

	existing, err := h.useCase.GetNotificationEventTypeByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Warn().Err(err).Uint64("id", id).Msg("⚠️ Notification event type not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification event type not found"})
		return
	}

	// Aplicar los cambios del request a la entidad existente
	updated := mappers.UpdateRequestToDomain(&req, existing)

	h.logger.Info().Uint64("id", id).Msg("🔄 Updating notification event type via use case")

	// Actualizar en el dominio
	err = h.useCase.UpdateNotificationEventType(c.Request.Context(), updated)
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("❌ Error updating notification event type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("✅ Notification event type updated successfully")

	// Convertir entidad actualizada a response HTTP
	response := mappers.DomainToResponse(*updated)
	c.JSON(http.StatusOK, response)
}
