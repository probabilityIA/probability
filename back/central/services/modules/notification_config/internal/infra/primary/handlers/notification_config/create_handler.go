package notification_config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/request"
)

// Create godoc
// @Summary Crear configuración de notificación
// @Description Crea una nueva configuración de notificación para una integración
// @Tags notification-config
// @Accept json
// @Produce json
// @Param body body request.CreateNotificationConfig true "Datos de la configuración"
// @Success 201 {object} response.NotificationConfig
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs [post]
func (h *handler) Create(c *gin.Context) {
	h.logger.Info().Msg("🌐 [POST /notification-configs] Request received")

	var req request.CreateNotificationConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Interface("request_body", req).
		Msg("📋 Request body parsed")

	// Convertir request HTTP a DTO de dominio usando mapper
	dto := mappers.CreateRequestToDomain(&req)

	h.logger.Info().Msg("➕ Creating notification config via use case")

	result, err := h.useCase.Create(c.Request.Context(), dto)
	if err != nil {
		if err == errors.ErrDuplicateConfig {
			h.logger.Warn().Msg("⚠️ Duplicate notification config")
			c.JSON(http.StatusConflict, gin.H{"error": "A similar notification config already exists"})
			return
		}
		h.logger.Error().Err(err).Msg("❌ Error creating notification config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Uint("id", result.ID).Msg("✅ Notification config created successfully")

	// Convertir DTO de dominio a response HTTP usando mapper
	response := mappers.DomainToResponse(*result)
	c.JSON(http.StatusCreated, response)
}
