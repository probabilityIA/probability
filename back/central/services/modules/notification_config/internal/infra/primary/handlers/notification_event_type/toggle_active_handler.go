package notification_event_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/mappers"
)

// ToggleActive godoc
// @Summary Activar/Desactivar tipo de evento
// @Description Cambia el estado activo/inactivo de un tipo de evento de notificación
// @Tags notification-event-types
// @Produce json
// @Param id path int true "ID del tipo de evento"
// @Success 200 {object} response.NotificationEventType
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-event-types/{id}/toggle-active [patch]
func (h *handler) ToggleActive(c *gin.Context) {
	h.logger.Info().Msg("🌐 [PATCH /notification-event-types/:id/toggle-active] Request received")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("📋 Path parameter parsed")

	// Obtener la entidad existente
	h.logger.Info().Uint64("id", id).Msg("🔍 Fetching existing notification event type")

	existing, err := h.useCase.GetNotificationEventTypeByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Warn().Err(err).Uint64("id", id).Msg("⚠️ Notification event type not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification event type not found"})
		return
	}

	// Invertir el estado activo
	existing.IsActive = !existing.IsActive

	h.logger.Info().Uint64("id", id).Bool("new_is_active", existing.IsActive).Msg("🔄 Toggling active state via use case")

	// Actualizar en el dominio
	err = h.useCase.UpdateNotificationEventType(c.Request.Context(), existing)
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("❌ Error toggling notification event type state")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Uint64("id", id).Bool("is_active", existing.IsActive).Msg("✅ Notification event type state toggled successfully")

	// Convertir entidad actualizada a response HTTP
	response := mappers.DomainToResponse(*existing)
	c.JSON(http.StatusOK, response)
}
