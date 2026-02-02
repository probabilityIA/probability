package notification_type

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type/request"
)

// Create godoc
// @Summary Crear tipo de notificación
// @Description Crea un nuevo tipo de notificación (WhatsApp, Email, SMS, etc.)
// @Tags notification-types
// @Accept json
// @Produce json
// @Param body body request.CreateNotificationType true "Datos del tipo de notificación"
// @Success 201 {object} response.NotificationType
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/notification-types [post]
func (h *handler) Create(c *gin.Context) {
	h.logger.Info().Msg("🌐 [POST /notification-types] Request received")

	var req request.CreateNotificationType
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Interface("request_body", req).
		Msg("📋 Request body parsed")

	// Convertir request HTTP a entidad de dominio
	entity := mappers.CreateRequestToDomain(&req)

	h.logger.Info().Msg("➕ Creating notification type via use case")

	err := h.useCase.CreateNotificationType(c.Request.Context(), entity)
	if err != nil {
		// TODO: Manejar error de duplicado si es necesario
		h.logger.Error().Err(err).Msg("❌ Error creating notification type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Uint("id", entity.ID).Msg("✅ Notification type created successfully")

	// Convertir entidad de dominio a response HTTP
	response := mappers.DomainToResponse(*entity)
	c.JSON(http.StatusCreated, response)
}
