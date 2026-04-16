package notification_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type/request"
)

// Update godoc
// @Summary Actualizar tipo de notificación
// @Description Actualiza un tipo de notificación existente
// @Tags notification-types
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de notificación"
// @Param body body request.UpdateNotificationType true "Datos a actualizar"
// @Success 200 {object} response.NotificationType
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-types/{id} [put]
func (h *handler) Update(c *gin.Context) {
	h.logger.Info().Msg("🌐 [PUT /notification-types/:id] Request received")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("📋 Path parameter parsed")

	var req request.UpdateNotificationType
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Interface("request_body", req).
		Msg("📋 Request body parsed")

	// Obtener la entidad existente
	h.logger.Info().Uint64("id", id).Msg("🔍 Fetching existing notification type")

	existing, err := h.useCase.GetNotificationTypeByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Warn().Err(err).Uint64("id", id).Msg("⚠️ Notification type not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification type not found"})
		return
	}

	// Aplicar los cambios del request a la entidad existente
	updated := mappers.UpdateRequestToDomain(&req, existing)

	h.logger.Info().Uint64("id", id).Msg("🔄 Updating notification type via use case")

	// Actualizar en el dominio
	err = h.useCase.UpdateNotificationType(c.Request.Context(), updated)
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("❌ Error updating notification type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Uint64("id", id).Msg("✅ Notification type updated successfully")

	// Convertir entidad actualizada a response HTTP
	response := mappers.DomainToResponse(*updated)
	c.JSON(http.StatusOK, response)
}
