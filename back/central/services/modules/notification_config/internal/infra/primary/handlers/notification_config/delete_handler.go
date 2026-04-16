package notification_config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// Delete godoc
// @Summary Eliminar configuración
// @Description Elimina una configuración de notificación
// @Tags notification-config
// @Accept json
// @Produce json
// @Param id path uint true "ID de la configuración"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/{id} [delete]
func (h *handler) Delete(c *gin.Context) {
	h.logger.Info().Msg("🌐 [DELETE /notification-configs/:id] Request received")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id_param", idStr).Msg("❌ Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("🗑️ Deleting notification config via use case")

	if err := h.useCase.Delete(c.Request.Context(), uint(id)); err != nil {
		if err == errors.ErrNotificationConfigNotFound {
			h.logger.Warn().Uint64("id", id).Msg("⚠️ Notification config not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification config not found"})
			return
		}
		h.logger.Error().Err(err).Uint64("id", id).Msg("❌ Error deleting notification config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("✅ Notification config deleted successfully")

	c.Status(http.StatusNoContent)
}
