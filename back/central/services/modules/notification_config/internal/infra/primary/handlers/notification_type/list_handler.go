package notification_type

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type/mappers"
)

// List godoc
// @Summary Listar tipos de notificación
// @Description Obtiene todos los tipos de notificación disponibles
// @Tags notification-types
// @Accept json
// @Produce json
// @Success 200 {array} response.NotificationType
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-types [get]
func (h *handler) List(c *gin.Context) {
	types, err := h.useCase.GetNotificationTypes(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Error listing notification types")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convertir entidades de dominio a respuestas HTTP
	response := mappers.DomainListToResponse(types)
	c.JSON(http.StatusOK, response)
}
