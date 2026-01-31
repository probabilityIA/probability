package notification_event_type

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/request"
)

// Create godoc
// @Summary Crear tipo de evento de notificación
// @Description Crea un nuevo tipo de evento para un tipo de notificación específico
// @Tags notification-event-types
// @Accept json
// @Produce json
// @Param body body request.CreateNotificationEventType true "Datos del tipo de evento"
// @Success 201 {object} response.NotificationEventType
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-event-types [post]
func (h *handler) Create(c *gin.Context) {
	var req request.CreateNotificationEventType
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir request HTTP a entidad de dominio
	entity := mappers.CreateRequestToDomain(&req)

	err := h.useCase.CreateNotificationEventType(c.Request.Context(), entity)
	if err != nil {
		// TODO: Manejar error de tipo de notificación no encontrado
		h.logger.Error().Err(err).Msg("Error creating notification event type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convertir entidad de dominio a response HTTP
	response := mappers.DomainToResponse(*entity)
	c.JSON(http.StatusCreated, response)
}
