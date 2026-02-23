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

// syncOrdersParams representa los parámetros para sincronizar órdenes (tipo local para evitar ciclo de importación)
type syncOrdersParams struct {
	CreatedAtMin      *time.Time
	CreatedAtMax      *time.Time
	Status            string
	FinancialStatus   string
	FulfillmentStatus string
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

	// Parsear el body para obtener parámetros opcionales
	var req SyncOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Si no hay body o está vacío, usar comportamiento por defecto
		req = SyncOrdersRequest{}
	}

	// Parsear fechas si están presentes
	var syncParams *syncOrdersParams
	if req.CreatedAtMin != nil || req.CreatedAtMax != nil || req.Status != "" || req.FinancialStatus != "" || req.FulfillmentStatus != "" {
		syncParams = &syncOrdersParams{}

		if req.CreatedAtMin != nil && *req.CreatedAtMin != "" {
			parsedMin, err := parseFlexibleDate(*req.CreatedAtMin)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid created_at_min format. Use RFC3339 or YYYY-MM-DD",
					"error":   err.Error(),
				})
				return
			}
			// Ajustar a inicio del día en UTC
			parsedMin = time.Date(parsedMin.Year(), parsedMin.Month(), parsedMin.Day(), 0, 0, 0, 0, time.UTC)
			syncParams.CreatedAtMin = &parsedMin
		}

		if req.CreatedAtMax != nil && *req.CreatedAtMax != "" {
			parsedMax, err := parseFlexibleDate(*req.CreatedAtMax)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid created_at_max format. Use RFC3339 or YYYY-MM-DD",
					"error":   err.Error(),
				})
				return
			}
			// Ajustar a fin del día en UTC (23:59:59.999)
			parsedMax = time.Date(parsedMax.Year(), parsedMax.Month(), parsedMax.Day(), 23, 59, 59, 999999999, time.UTC)
			syncParams.CreatedAtMax = &parsedMax
		}

		if req.Status != "" {
			syncParams.Status = req.Status
		}
		if req.FinancialStatus != "" {
			syncParams.FinancialStatus = req.FinancialStatus
		}
		if req.FulfillmentStatus != "" {
			syncParams.FulfillmentStatus = req.FulfillmentStatus
		}
	}

	// Usar método con parámetros si están presentes, sino usar el método por defecto
	var err error
	if syncParams != nil {
		err = h.usecase.SyncOrdersByIntegrationIDWithParams(c.Request.Context(), integrationID, syncParams)
	} else {
		err = h.usecase.SyncOrdersByIntegrationID(c.Request.Context(), integrationID)
	}

	if err != nil {
		h.logger.Error().Err(err).Str("integration_id", integrationID).Msg("Failed to sync orders")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to start sync",
			"error":   err.Error(),
		})
		return
	}

	message := "Order synchronization started in background"
	if syncParams != nil {
		message += " with custom filters"
	} else {
		message += " (last 30 days)"
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": message,
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

	if err := h.usecase.SyncOrdersByBusiness(c.Request.Context(), uint(businessID)); err != nil {
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
