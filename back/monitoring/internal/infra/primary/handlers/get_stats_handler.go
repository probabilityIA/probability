package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers/mappers"
)

func (h *handler) GetStats(c *gin.Context) {
	id := c.Param("id")

	stats, err := h.useCase.GetContainerStats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.StatsToResponse(stats))
}
