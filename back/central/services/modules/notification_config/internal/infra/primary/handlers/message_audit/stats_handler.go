package message_audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/request"
)

// Stats godoc
// @Summary Obtener estadísticas de mensajes
// @Description Obtiene estadísticas agregadas de mensajes outbound
// @Tags message-audit
// @Produce json
// @Param business_id query uint true "ID del negocio"
// @Param date_from query string false "Fecha inicio (YYYY-MM-DD)"
// @Param date_to query string false "Fecha fin (YYYY-MM-DD)"
// @Success 200 {object} response.MessageAuditStats
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/message-audit/stats [get]
func (h *handler) Stats(c *gin.Context) {
	h.logger.Info().Msg("[GET /notification-configs/message-audit/stats] Request received")

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var query request.StatsMessageAudit
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error().Err(err).Msg("Invalid query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.GetMessageAuditStats(c.Request.Context(), businessID, query.DateFrom, query.DateTo)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error getting message audit stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, mappers.DomainToStatsResponse(result))
}
