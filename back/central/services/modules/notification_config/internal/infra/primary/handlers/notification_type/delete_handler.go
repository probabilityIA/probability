package notification_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary Eliminar tipo de notificación
// @Description Elimina un tipo de notificación por su ID
// @Tags notification-types
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de notificación"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-types/{id} [delete]
func (h *handler) Delete(c *gin.Context) {
	h.logger.Info().Msg("🌐 [DELETE /notification-types/:id] Request received")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("🗑️ Deleting notification type via use case")

	err = h.useCase.DeleteNotificationType(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("❌ Error deleting notification type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("✅ Notification type deleted successfully")

	c.Status(http.StatusNoContent)
}
