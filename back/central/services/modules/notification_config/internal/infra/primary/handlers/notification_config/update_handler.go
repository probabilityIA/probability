package notification_config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/request"
)

// Update godoc
// @Summary Actualizar configuración
// @Description Actualiza una configuración de notificación existente
// @Tags notification-config
// @Accept json
// @Produce json
// @Param id path uint true "ID de la configuración"
// @Param body body request.UpdateNotificationConfig true "Datos a actualizar"
// @Success 200 {object} response.NotificationConfig
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/{id} [put]
func (h *handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req request.UpdateNotificationConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir request HTTP a DTO de dominio usando mapper
	dto := mappers.UpdateRequestToDomain(&req)

	result, err := h.useCase.Update(c.Request.Context(), uint(id), dto)
	if err != nil {
		if err == errors.ErrNotificationConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification config not found"})
			return
		}
		h.logger.Error().Err(err).Msg("Error updating notification config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convertir DTO de dominio a response HTTP usando mapper
	response := mappers.DomainToResponse(*result)
	c.JSON(http.StatusOK, response)
}
