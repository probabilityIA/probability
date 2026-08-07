package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetStats(c *gin.Context) {
	stats, err := h.uc.GetStats(c.Request.Context())
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Msg("Error al calcular estadisticas de prospectos comerciales")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response.FromStats(stats),
	})
}
