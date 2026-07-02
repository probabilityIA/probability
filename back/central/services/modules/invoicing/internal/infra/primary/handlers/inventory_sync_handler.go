package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

type inventorySyncRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id,omitempty"`
}

func (h *handler) SyncInventory(c *gin.Context) {
	var req inventorySyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id is required"})
		return
	}

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business context not found"})
		return
	}

	if businessID == 0 {
		if req.BusinessID == nil || *req.BusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required for super admin"})
			return
		}
		businessID = *req.BusinessID
	}

	correlationID, err := h.useCase.RequestInventorySync(c.Request.Context(), businessID, req.IntegrationID)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to start inventory sync")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"correlation_id": correlationID,
		"message":        "Sincronizacion de inventario iniciada. Recibiras el progreso por SSE.",
	})
}

func (h *handler) ListSiigoWarehouses(c *gin.Context) {
	var req inventorySyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id is required"})
		return
	}

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business context not found"})
		return
	}

	if businessID == 0 {
		if req.BusinessID == nil || *req.BusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required for super admin"})
			return
		}
		businessID = *req.BusinessID
	}

	correlationID, err := h.useCase.RequestListSiigoWarehouses(c.Request.Context(), businessID, req.IntegrationID)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to start list siigo warehouses")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"correlation_id": correlationID,
		"message":        "Consulta de bodegas Siigo iniciada. Recibiras el resultado por SSE.",
	})
}
