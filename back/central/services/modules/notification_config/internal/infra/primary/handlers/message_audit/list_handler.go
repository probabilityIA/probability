package message_audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/request"
)

// List godoc
// @Summary Listar logs de auditoría de mensajes
// @Description Obtiene logs de mensajes con filtros y paginación
// @Tags message-audit
// @Produce json
// @Param business_id query uint true "ID del negocio"
// @Param status query string false "Estado del mensaje (sent, delivered, read, failed)"
// @Param direction query string false "Dirección (outbound, inbound)"
// @Param template_name query string false "Nombre del template (búsqueda parcial)"
// @Param date_from query string false "Fecha inicio (YYYY-MM-DD)"
// @Param date_to query string false "Fecha fin (YYYY-MM-DD)"
// @Param page query int false "Página" default(1)
// @Param page_size query int false "Tamaño de página" default(20)
// @Success 200 {object} response.PaginatedMessageAuditResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/integrations/notification-configs/message-audit [get]
func (h *handler) List(c *gin.Context) {
	h.logger.Info().Msg("[GET /notification-configs/message-audit] Request received")

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var query request.ListMessageAudit
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error().Err(err).Msg("Invalid query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := mappers.ListRequestToDomain(&query, businessID)

	result, err := h.useCase.ListMessageAudit(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error listing message audit logs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info().
		Int64("total", result.Total).
		Int("page", result.Page).
		Msg("Message audit logs fetched successfully")

	c.JSON(http.StatusOK, mappers.DomainToListResponse(result))
}
