package notification_config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/mappers"
)

// GetByID godoc
// @Summary Obtener configuración por ID
// @Description Obtiene una configuración de notificación por su ID
// @Tags notification-config
// @Accept json
// @Produce json
// @Param id path uint true "ID de la configuración"
// @Success 200 {object} response.NotificationConfig
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/{id} [get]
func (h *handler) GetByID(c *gin.Context) {
	h.logger.Info().Msg("🌐 [GET /notification-configs/:id] Request received")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id_param", idStr).Msg("❌ Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("🔍 Fetching notification config by ID")

	result, err := h.useCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrNotificationConfigNotFound {
			h.logger.Warn().Uint64("id", id).Msg("⚠️ Notification config not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification config not found"})
			return
		}
		h.logger.Error().Err(err).Uint64("id", id).Msg("❌ Error getting notification config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("✅ Notification config fetched successfully")

	// Convertir DTO de dominio a response HTTP usando mapper
	response := mappers.DomainToResponse(*result)
	c.JSON(http.StatusOK, response)
}
