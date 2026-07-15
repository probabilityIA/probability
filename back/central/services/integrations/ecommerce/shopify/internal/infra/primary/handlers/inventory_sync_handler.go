package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

type syncInventoryRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id"`
}

func (h *ShopifyHandler) SyncInventoryHandler(c *gin.Context) {
	var req syncInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return
	}
	if businessID == 0 {
		if req.BusinessID == nil || *req.BusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
			return
		}
		businessID = *req.BusinessID
	}

	integrationID := strconv.FormatUint(uint64(req.IntegrationID), 10)
	correlationID := uuid.New().String()

	go func() {
		ctx := context.Background()
		if err := h.useCase.SyncInventory(ctx, integrationID, businessID, correlationID); err != nil {
			h.logger.Error(ctx).Err(err).Msg("Error sincronizando inventario a Shopify")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"correlation_id": correlationID,
		"message":        "Sincronizacion de inventario iniciada",
	})
}
