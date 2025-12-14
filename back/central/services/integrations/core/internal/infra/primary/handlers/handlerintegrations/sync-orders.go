package handlerintegrations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *IntegrationHandler) SyncOrdersByIntegrationIDHandler(c *gin.Context) {
	integrationID := c.Param("id")
	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "integration_id is required",
		})
		return
	}

	if h.orderSyncSvc == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Order sync service not configured",
		})
		return
	}

	if err := h.orderSyncSvc.SyncOrdersByIntegrationID(c.Request.Context(), integrationID); err != nil {
		h.logger.Error().Err(err).Str("integration_id", integrationID).Msg("Failed to sync orders")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to start sync",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Order synchronization started in background",
	})
}

func (h *IntegrationHandler) SyncOrdersByBusinessHandler(c *gin.Context) {
	businessIDStr := c.Param("business_id")
	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid business_id",
		})
		return
	}

	if h.orderSyncSvc == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Order sync service not configured",
		})
		return
	}

	if err := h.orderSyncSvc.SyncOrdersByBusiness(c.Request.Context(), uint(businessID)); err != nil {
		h.logger.Error().Err(err).Uint("business_id", uint(businessID)).Msg("Failed to sync orders")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to start sync",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Order synchronization started in background for all integrations",
	})
}
