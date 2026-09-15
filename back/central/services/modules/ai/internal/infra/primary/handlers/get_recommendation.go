package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/mappers"
)

func (h *Handlers) GetRecommendation(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "origin and destination required"})
		return
	}

	result, err := h.uc.GetRecommendation(c.Request.Context(), origin, destination)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Error getting recommendation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendation"})
		return
	}

	c.JSON(http.StatusOK, mappers.FromRecommendation(result))
}
