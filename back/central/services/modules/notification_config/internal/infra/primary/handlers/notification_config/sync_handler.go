package notification_config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/request"
)

// SyncByIntegration godoc
// @Summary Sincronizar reglas de notificación por integración
// @Description Crea, actualiza y elimina reglas de notificación para una integración en batch
// @Tags notification-config
// @Accept json
// @Produce json
// @Param body body request.SyncNotificationConfigs true "Datos de sincronización"
// @Param business_id query int false "Business ID (required for super admin)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/sync [put]
func (h *handler) SyncByIntegration(c *gin.Context) {
	h.logger.Info().Msg("🌐 [PUT /notification-configs/sync] Request received")

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.SyncNotificationConfigs
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("❌ Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Uint("business_id", businessID).
		Uint("integration_id", req.IntegrationID).
		Int("rules_count", len(req.Rules)).
		Msg("📋 Sync request parsed")

	// Convertir request a DTO de dominio
	syncDTO := dtos.SyncNotificationConfigsDTO{
		BusinessID:    businessID,
		IntegrationID: req.IntegrationID,
		Rules:         make([]dtos.SyncRuleDTO, len(req.Rules)),
	}

	for i, rule := range req.Rules {
		syncDTO.Rules[i] = dtos.SyncRuleDTO{
			ID:                      rule.ID,
			NotificationTypeID:      rule.NotificationTypeID,
			NotificationEventTypeID: rule.NotificationEventTypeID,
			Enabled:                 rule.Enabled,
			Description:             rule.Description,
			OrderStatusIDs:          rule.OrderStatusIDs,
		}
	}

	result, err := h.useCase.SyncByIntegration(c.Request.Context(), syncDTO)
	if err != nil {
		if err == errors.ErrDuplicateConfig {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate rules detected in sync request"})
			return
		}
		if err == errors.ErrNotificationConfigNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "One or more rule IDs not found"})
			return
		}
		h.logger.Error().Err(err).Msg("❌ Error syncing notification configs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convertir configs del resultado a response HTTP
	configResponses := mappers.DomainListToResponse(result.Configs)

	h.logger.Info().
		Int("created", result.Created).
		Int("updated", result.Updated).
		Int("deleted", result.Deleted).
		Msg("✅ Sync completed")

	c.JSON(http.StatusOK, gin.H{
		"created": result.Created,
		"updated": result.Updated,
		"deleted": result.Deleted,
		"configs": configResponses,
	})
}
