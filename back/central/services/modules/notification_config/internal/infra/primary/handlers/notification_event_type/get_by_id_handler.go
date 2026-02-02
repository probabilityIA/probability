package notification_event_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/mappers"
)

// GetByID godoc
// @Summary Obtener tipo de evento por ID
// @Description Obtiene un tipo de evento de notificación específico por su ID
// @Tags notification-event-types
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de evento"
// @Success 200 {object} response.NotificationEventType
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-event-types/{id} [get]
func (h *handler) GetByID(c *gin.Context) {
	h.logger.Info().Msg("🌐 [GET /notification-event-types/:id] Request received")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("🔍 Fetching notification event type by ID")

	entity, err := h.useCase.GetNotificationEventTypeByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Warn().Err(err).Uint64("id", id).Msg("⚠️ Notification event type not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification event type not found"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("✅ Notification event type fetched successfully")

	// Convertir entidad de dominio a response HTTP
	response := mappers.DomainToResponse(*entity)
	c.JSON(http.StatusOK, response)
}
