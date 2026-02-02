package notification_event_type

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary Eliminar tipo de evento
// @Description Elimina un tipo de evento de notificación por su ID
// @Tags notification-event-types
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de evento"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-event-types/{id} [delete]
func (h *handler) Delete(c *gin.Context) {
	h.logger.Info().Msg("🌐 [DELETE /notification-event-types/:id] Request received")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("🗑️ Deleting notification event type via use case")

	err = h.useCase.DeleteNotificationEventType(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("❌ Error deleting notification event type")

		// Si el error contiene un mensaje específico sobre configuraciones activas, usar 409 Conflict
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "configuración") || strings.Contains(errorMsg, "siendo usado") {
			c.JSON(http.StatusConflict, gin.H{"error": errorMsg})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("✅ Notification event type deleted successfully")

	c.Status(http.StatusNoContent)
}
