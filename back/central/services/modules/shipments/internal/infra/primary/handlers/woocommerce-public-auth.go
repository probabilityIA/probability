package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) authWooPublic(c *gin.Context) (uint, bool) {
	integrationID64, err := strconv.ParseUint(c.Param("integration_id"), 10, 64)
	if err != nil || integrationID64 == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_integration"})
		return 0, false
	}

	resolved, err := h.resolveWoo(c.Request.Context(), uint(integrationID64))
	if err != nil || resolved == nil || !resolved.Found || resolved.Revoked ||
		!h.wooTokenMatches(uint(integrationID64), resolved.Salt, c.GetHeader("X-Probability-Token")) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return 0, false
	}

	return uint(integrationID64), true
}
