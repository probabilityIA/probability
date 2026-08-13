package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type inventorySyncRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	BusinessID    *uint    `json:"business_id"`
	SKUs          []string `json:"skus"`
}

func (h *meliHandler) SyncInventory(c *gin.Context) {
	var req inventorySyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}
	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()

	go func() {
		ctx := context.Background()
		if err := h.useCase.SyncInventory(ctx, integrationID, businessID, correlationID, req.SKUs...); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando inventario a MercadoLibre")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de inventario iniciada",
	})
}
