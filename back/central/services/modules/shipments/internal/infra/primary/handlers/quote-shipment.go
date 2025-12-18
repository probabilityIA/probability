package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) QuoteShipment(c *gin.Context) {
	var req domain.EnvioClickQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log received request
	// Note: We need a logger in the Handler struct ideally, using fmt for now or assuming h.uc.Logger access if available?
	// Handlers doesn't have Logger directly usually, but let's assume standard stdout for debug or rely on the client log we just added.
	// Actually, Gin logs requests automatically. Let's just trust Client logs for now, or print to stdout if needed.

	resp, err := h.envioClickUC.QuoteShipment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
