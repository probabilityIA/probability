package handlerintegrations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

func (h *IntegrationHandler) GetIntegrationStatsHandler(c *gin.Context) {
	var businessID uint
	if ctxBusinessID, exists := c.Get("business_id"); exists {
		if bID, ok := ctxBusinessID.(uint); ok {
			businessID = bID
		}
	}

	if businessID == 0 {
		param := c.Query("business_id")
		if param == "" {
			c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
				Success: false,
				Message: "business_id es requerido para super admin",
			})
			return
		}
		parsed, err := strconv.ParseUint(param, 10, 64)
		if err != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
				Success: false,
				Message: "business_id invalido",
			})
			return
		}
		businessID = uint(parsed)
	}

	stats, err := h.usecase.GetIntegrationStats(c.Request.Context(), businessID)
	if err != nil {
		h.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error al obtener stats de integraciones")
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error interno al obtener estadisticas de integraciones",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.ToIntegrationStatsResponse(stats))
}
