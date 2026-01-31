package notification_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type/mappers"
)

// GetByID godoc
// @Summary Obtener tipo de notificación por ID
// @Description Obtiene un tipo de notificación específico por su ID
// @Tags notification-types
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de notificación"
// @Success 200 {object} response.NotificationType
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-types/{id} [get]
func (h *handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	entity, err := h.useCase.GetNotificationTypeByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("Error getting notification type by ID")
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification type not found"})
		return
	}

	// Convertir entidad de dominio a response HTTP
	response := mappers.DomainToResponse(*entity)
	c.JSON(http.StatusOK, response)
}
