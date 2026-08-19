package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListCatalogs(c *gin.Context) {
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

	catalogs, err := h.useCase.ListCatalogs(ctx, uint(integrationID))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint64("integration_id", integrationID).Msg("No se pudieron consultar los catalogos de Siigo")
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": catalogs})
}
