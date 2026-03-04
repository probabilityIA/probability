package handlerintegrations

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
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
			"message": "integration_id es requerido",
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
					"message": "Formato inválido en created_at_min. Use RFC3339 o YYYY-MM-DD",
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
					"message": "Formato inválido en created_at_max. Use RFC3339 o YYYY-MM-DD",
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

	// Decidir entre flujo directo y flujo por lotes
	useBatches := false

	if syncParams != nil && syncParams.CreatedAtMin != nil && syncParams.CreatedAtMax != nil {
		dateRange := syncParams.CreatedAtMax.Sub(*syncParams.CreatedAtMin)
		if dateRange > 14*24*time.Hour {
			useBatches = true
		}
	}

	if useBatches {
		// Flujo por lotes: publica mensajes a la cola y retorna rápido
		batchParams := &domain.SyncBatchParams{
			CreatedAtMin:      syncParams.CreatedAtMin,
			CreatedAtMax:      syncParams.CreatedAtMax,
			Status:            syncParams.Status,
			FinancialStatus:   syncParams.FinancialStatus,
			FulfillmentStatus: syncParams.FulfillmentStatus,
		}
		if err := h.usecase.SyncOrdersByIntegrationIDWithBatches(c.Request.Context(), integrationID, batchParams); err != nil {
			h.logger.Error().Err(err).Str("integration_id", integrationID).Msg("Error al iniciar sincronización por lotes")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error al iniciar sincronización",
				"error":   err.Error(),
			})
			return
		}
	} else if syncParams != nil {
		// Flujo directo con params: ejecutar en goroutine (es sincrónico y puede tardar)
		// Los eventos SSE (started/completed/failed) notificarán al frontend
		go func() {
			if err := h.usecase.SyncOrdersByIntegrationIDWithParams(context.Background(), integrationID, syncParams); err != nil {
				h.logger.Error().Err(err).Str("integration_id", integrationID).Msg("Error en sincronización directa con params")
			}
		}()
	} else {
		// Flujo directo sin params: ejecutar en goroutine
		go func() {
			if err := h.usecase.SyncOrdersByIntegrationID(context.Background(), integrationID); err != nil {
				h.logger.Error().Err(err).Str("integration_id", integrationID).Msg("Error en sincronización directa")
			}
		}()
	}

	message := "Sincronización de órdenes iniciada en segundo plano"
	if useBatches {
		message += " usando procesamiento por lotes (rango de fechas grande)"
	} else if syncParams != nil {
		message += " con filtros personalizados"
	} else {
		message += " (últimos 30 días)"
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
			"message": "business_id inválido",
		})
		return
	}

	if err := h.usecase.SyncOrdersByBusiness(c.Request.Context(), uint(businessID)); err != nil {
		h.logger.Error().Err(err).Uint("business_id", uint(businessID)).Msg("Error al sincronizar órdenes por negocio")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al iniciar sincronización",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Sincronización de órdenes iniciada en segundo plano para todas las integraciones",
	})
}
