package handlerintegrations

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SyncOrdersRequest representa los parámetros para sincronizar órdenes
type SyncOrdersRequest struct {
	CreatedAtMin      *string `json:"created_at_min"`     // Formato: RFC3339 o YYYY-MM-DD
	CreatedAtMax      *string `json:"created_at_max"`     // Formato: RFC3339 o YYYY-MM-DD
	Status            string  `json:"status"`             // any, open, closed, cancelled
	FinancialStatus   string  `json:"financial_status"`   // any, paid, pending, refunded, etc.
	FulfillmentStatus string  `json:"fulfillment_status"` // any, shipped, partial, unshipped, etc.
}

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

	// Intentar parsear el body para obtener parámetros opcionales
	// Por ahora, los parámetros se ignoran y se usa el comportamiento por defecto
	// El soporte completo de parámetros se implementará cuando se agregue a la interfaz IOrderSyncService
	var req SyncOrdersRequest
	c.ShouldBindJSON(&req)

	// Sin parámetros, usar comportamiento por defecto (últimos 30 días)
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
		"message": "Order synchronization started in background (last 30 days)",
	})
}

// parseFlexibleDate parsea una fecha en varios formatos posibles
func parseFlexibleDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, &time.ParseError{Layout: "RFC3339 or YYYY-MM-DD", Value: dateStr}
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
