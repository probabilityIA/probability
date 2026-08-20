package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SearchProducts(c *gin.Context) {
	ctx := c.Request.Context()

	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	integrationID, err := strconv.ParseUint(integrationIDStr, 10, 32)
	if err != nil || integrationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id invalido"})
		return
	}

	termino := strings.TrimSpace(c.Query("q"))
	limite, _ := strconv.Atoi(c.Query("limit"))

	items, err := h.useCase.SearchProducts(ctx, uint(integrationID), termino, limite)
	if err != nil {
		h.log.Error(ctx).Err(err).
			Uint64("integration_id", integrationID).
			Str("termino", termino).
			Msg("No se pudieron buscar productos en Siigo")
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}
