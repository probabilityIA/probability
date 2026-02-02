package notification_config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/request"
)

// List godoc
// @Summary Listar configuraciones de notificación
// @Description Obtiene una lista de configuraciones con filtros opcionales
// @Tags notification-config
// @Accept json
// @Produce json
// @Param integration_id query uint false "ID de la integración"
// @Param notification_type query string false "Tipo de notificación (whatsapp, email, sms)"
// @Param is_active query bool false "Filtrar por activas/inactivas"
// @Param trigger query string false "Filtrar por trigger"
// @Success 200 {array} response.NotificationConfig
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs [get]
func (h *handler) List(c *gin.Context) {
	h.logger.Info().Msg("🌐 [GET /notification-configs] Request received")

	var query request.FilterNotificationConfig
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Interface("query_params", query).
		Msg("📋 Query params parsed")

	// Convertir query params a DTO de dominio usando mapper
	filters := mappers.FilterRequestToDomain(&query)

	h.logger.Info().Msg("🔍 Fetching notification configs from use case")

	result, err := h.useCase.List(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("❌ Error listing notification configs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().Int("count", len(result)).Msg("✅ Notification configs fetched successfully")

	// Convertir lista de DTOs de dominio a responses HTTP usando mapper
	responses := mappers.DomainListToResponse(result)
	c.JSON(http.StatusOK, responses)
}
