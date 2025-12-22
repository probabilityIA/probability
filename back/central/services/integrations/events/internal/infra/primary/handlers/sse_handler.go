package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IntegrationSSEHandler maneja las conexiones Server-Sent Events para eventos de integraciones
type IntegrationSSEHandler struct {
	eventManager domain.IIntegrationEventPublisher
	logger       log.ILogger
}

// IntegrationSSEHandlerInterface define la interfaz del handler SSE
type IntegrationSSEHandlerInterface interface {
	HandleSSE(c *gin.Context)
	GetManager() domain.IIntegrationEventPublisher
	GetSyncStatus(c *gin.Context)
}

// NewIntegrationSSEHandler crea un nuevo handler de SSE para eventos de integraciones
func New(eventManager domain.IIntegrationEventPublisher, logger log.ILogger) IntegrationSSEHandlerInterface {
	return &IntegrationSSEHandler{
		eventManager: eventManager,
		logger:       logger,
	}
}

// GetManager retorna el manager interno para acceso externo
func (h *IntegrationSSEHandler) GetManager() domain.IIntegrationEventPublisher {
	return h.eventManager
}

// HandleSSE maneja la conexión SSE por business_id con filtros opcionales
func (h *IntegrationSSEHandler) HandleSSE(c *gin.Context) {
	h.logger.Info(c.Request.Context()).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("remote_addr", c.ClientIP()).
		Msg("Integration SSE endpoint called")

	// Manejar preflight OPTIONS
	if c.Request.Method == "OPTIONS" {
		h.setupSSEHeaders(c.Writer)
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	// Obtener business_id
	var businessID uint
	businessIDFromParam := c.Param("businessID")
	businessIDFromQuery := c.Query("business_id")

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"location": "sse_handler.go:HandleSSE",
			"message":  "Parsing business_id from request",
			"data": map[string]interface{}{
				"path":                  c.Request.URL.Path,
				"businessID_from_param": businessIDFromParam,
				"businessID_from_query": businessIDFromQuery,
				"query_string":          c.Request.URL.RawQuery,
			},
			"timestamp":    time.Now().UnixMilli(),
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "E",
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	if businessIDFromParam != "" {
		if id, err := strconv.ParseUint(businessIDFromParam, 10, 32); err == nil {
			businessID = uint(id)
		} else {
			// #region agent log
			if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				logData, _ := json.Marshal(map[string]interface{}{
					"location": "sse_handler.go:HandleSSE",
					"message":  "Failed to parse businessID from param",
					"data": map[string]interface{}{
						"businessID_from_param": businessIDFromParam,
						"error":                 err.Error(),
					},
					"timestamp":    time.Now().UnixMilli(),
					"sessionId":    "debug-session",
					"runId":        "run1",
					"hypothesisId": "E",
				})
				f.WriteString(string(logData) + "\n")
				f.Close()
			}
			// #endregion
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "invalid business_id",
			})
			return
		}
	} else if businessIDFromQuery != "" {
		if id, err := strconv.ParseUint(businessIDFromQuery, 10, 32); err == nil {
			businessID = uint(id)
		} else {
			// #region agent log
			if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				logData, _ := json.Marshal(map[string]interface{}{
					"location": "sse_handler.go:HandleSSE",
					"message":  "Failed to parse businessID from query",
					"data": map[string]interface{}{
						"businessID_from_query": businessIDFromQuery,
						"error":                 err.Error(),
					},
					"timestamp":    time.Now().UnixMilli(),
					"sessionId":    "debug-session",
					"runId":        "run1",
					"hypothesisId": "E",
				})
				f.WriteString(string(logData) + "\n")
				f.Close()
			}
			// #endregion
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "invalid business_id",
			})
			return
		}
	}

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"location": "sse_handler.go:HandleSSE",
			"message":  "BusinessID parsed successfully",
			"data": map[string]interface{}{
				"final_business_id": businessID,
			},
			"timestamp":    time.Now().UnixMilli(),
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "E",
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	// Construir filtros
	filter := h.buildFilterFromQuery(c)

	h.logger.Info(c.Request.Context()).
		Uint("business_id", businessID).
		Interface("filter", filter).
		Msg("New integration SSE connection requested")

	// Configurar headers SSE
	h.setupSSEHeaders(c.Writer)

	// Agregar conexión
	connectionID := h.eventManager.AddConnection(businessID, filter, c.Writer)

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"location": "sse_handler.go:HandleSSE",
			"message":  "SSE connection registered",
			"data": map[string]interface{}{
				"connection_id": connectionID,
				"business_id":   businessID,
				"filter":        filter,
			},
			"timestamp":    time.Now().UnixMilli(),
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A",
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	// Precargar eventos del caché
	if businessID > 0 {
		h.preloadCacheEvents(c.Writer, businessID)
	}

	// Enviar mensaje de conexión establecida
	message := fmt.Sprintf("Conexión SSE establecida para eventos de integraciones (business %d)", businessID)
	if businessID == 0 {
		message = "Conexión SSE establecida para eventos de integraciones (super usuario)"
	}
	connectionEvent := fmt.Sprintf("event: connection_established\ndata: {\"message\":\"%s\",\"connection_id\":\"%s\",\"timestamp\":\"%s\"}\n\n",
		message, connectionID, time.Now().Format(time.RFC3339))

	if _, err := c.Writer.Write([]byte(connectionEvent)); err != nil {
		h.logger.Error(c.Request.Context()).
			Err(err).
			Str("connection_id", connectionID).
			Msg("Error writing connection message")
		return
	}

	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	// Mantener conexión viva
	h.keepConnectionAlive(c.Writer, connectionID, c.Request.Context())
}

// buildFilterFromQuery construye filtros desde query parameters
func (h *IntegrationSSEHandler) buildFilterFromQuery(c *gin.Context) *domain.IntegrationSSEFilter {
	filter := &domain.IntegrationSSEFilter{}

	// Filtro por integration_id
	if integrationIDStr := c.Query("integration_id"); integrationIDStr != "" {
		if id, err := strconv.ParseUint(integrationIDStr, 10, 32); err == nil {
			integrationID := uint(id)
			filter.IntegrationID = &integrationID
		}
	}

	// Filtro por event_types
	if eventTypesStr := c.Query("event_types"); eventTypesStr != "" {
		eventTypes := strings.Split(eventTypesStr, ",")
		filter.EventTypes = make([]domain.IntegrationEventType, 0, len(eventTypes))
		for _, et := range eventTypes {
			et = strings.TrimSpace(et)
			if et != "" {
				filter.EventTypes = append(filter.EventTypes, domain.IntegrationEventType(et))
			}
		}
	}

	// Filtro por order_ids
	if orderIDsStr := c.Query("order_ids"); orderIDsStr != "" {
		orderIDs := strings.Split(orderIDsStr, ",")
		filter.OrderIDs = make([]string, 0, len(orderIDs))
		for _, oid := range orderIDs {
			oid = strings.TrimSpace(oid)
			if oid != "" {
				filter.OrderIDs = append(filter.OrderIDs, oid)
			}
		}
	}

	return filter
}

// setupSSEHeaders configura los headers HTTP para SSE
func (h *IntegrationSSEHandler) setupSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control, Last-Event-ID")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// keepConnectionAlive mantiene la conexión viva
func (h *IntegrationSSEHandler) keepConnectionAlive(w http.ResponseWriter, connectionID string, ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	done := ctx.Done()
	flusher, hasFlusher := w.(http.Flusher)

	for {
		select {
		case <-ticker.C:
			h.sendSSEMessage(w, "keep-alive", "ping")
			if hasFlusher {
				flusher.Flush()
			}
		case <-done:
			h.eventManager.RemoveConnection(connectionID)
			h.logger.Info(ctx).
				Str("connection_id", connectionID).
				Msg("Integration SSE client disconnected")
			return
		}
	}
}

// preloadCacheEvents precarga eventos del caché
func (h *IntegrationSSEHandler) preloadCacheEvents(w http.ResponseWriter, businessID uint) {
	events := h.eventManager.GetRecentEventsByBusiness(businessID, 0)

	if len(events) > 0 {
		h.logger.Info(context.Background()).
			Uint("business_id", businessID).
			Int("cache_events_count", len(events)).
			Msg("Preloading integration events from cache")

		for _, event := range events {
			eventJSON := h.eventToSSEJSON(event)
			message := fmt.Sprintf("event: %s\ndata: %s\n\n", event.Type, eventJSON)
			w.Write([]byte(message))
		}

		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

// eventToSSEJSON convierte un evento a JSON para SSE
func (h *IntegrationSSEHandler) eventToSSEJSON(event domain.IntegrationEvent) string {
	eventData := map[string]interface{}{
		"id":             event.ID,
		"type":           event.Type,
		"integration_id": event.IntegrationID,
		"business_id":    event.BusinessID,
		"timestamp":      event.Timestamp,
		"metadata":       event.Metadata,
	}

	if event.Data != nil {
		if dataMap, ok := event.Data.(map[string]interface{}); ok {
			for key, value := range dataMap {
				eventData[key] = value
			}
		} else {
			eventData["data"] = event.Data
		}
	}

	jsonBytes, err := json.Marshal(eventData)
	if err != nil {
		h.logger.Error(context.Background()).
			Err(err).
			Msg("Error serializing integration event for SSE")
		return "{}"
	}

	return string(jsonBytes)
}

// sendSSEMessage envía un mensaje SSE formateado
func (h *IntegrationSSEHandler) sendSSEMessage(w http.ResponseWriter, eventType, data string) {
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, data)
	w.Write([]byte(message))

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// GetSyncStatus consulta el estado de sincronización para una integración
func (h *IntegrationSSEHandler) GetSyncStatus(c *gin.Context) {
	integrationIDStr := c.Param("integrationID")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "integration_id is required",
		})
		return
	}

	integrationID, err := strconv.ParseUint(integrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid integration_id",
		})
		return
	}

	// Obtener business_id del query o del contexto
	var businessID uint
	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		if id, err := strconv.ParseUint(businessIDStr, 10, 32); err == nil {
			businessID = uint(id)
		}
	}

	// Obtener eventos recientes para este business
	events := h.eventManager.GetRecentEventsByBusiness(businessID, 0)

	// Buscar el último evento de inicio de sincronización para esta integración
	var lastStartedEvent *domain.IntegrationEvent
	var lastCompletedEvent *domain.IntegrationEvent
	var lastFailedEvent *domain.IntegrationEvent

	for i := range events {
		event := events[i]
		if event.IntegrationID != uint(integrationID) {
			continue
		}

		switch event.Type {
		case domain.IntegrationEventTypeSyncStarted:
			if lastStartedEvent == nil || event.Timestamp.After(lastStartedEvent.Timestamp) {
				lastStartedEvent = &event
			}
		case domain.IntegrationEventTypeSyncCompleted:
			if lastCompletedEvent == nil || event.Timestamp.After(lastCompletedEvent.Timestamp) {
				lastCompletedEvent = &event
			}
		case domain.IntegrationEventTypeSyncFailed:
			if lastFailedEvent == nil || event.Timestamp.After(lastFailedEvent.Timestamp) {
				lastFailedEvent = &event
			}
		}
	}

	// Determinar si hay una sincronización en curso
	isInProgress := false
	var syncState map[string]interface{}

	if lastStartedEvent != nil {
		// Verificar si hay un evento de completado o fallido posterior
		hasCompletedAfter := lastCompletedEvent != nil && lastCompletedEvent.Timestamp.After(lastStartedEvent.Timestamp)
		hasFailedAfter := lastFailedEvent != nil && lastFailedEvent.Timestamp.After(lastStartedEvent.Timestamp)

		if !hasCompletedAfter && !hasFailedAfter {
			// Hay una sincronización en curso
			isInProgress = true

			// Extraer datos del evento de inicio
			// El Data puede venir como map[string]interface{} desde JSON o como struct
			var startedData domain.SyncStartedEvent
			if dataMap, ok := lastStartedEvent.Data.(map[string]interface{}); ok {
				// Convertir desde map
				if integrationIDVal, ok := dataMap["integration_id"].(float64); ok {
					startedData.IntegrationID = uint(integrationIDVal)
				}
				if integrationTypeVal, ok := dataMap["integration_type"].(string); ok {
					startedData.IntegrationType = integrationTypeVal
				}
				if startedAtVal, ok := dataMap["started_at"].(string); ok {
					if t, err := time.Parse(time.RFC3339, startedAtVal); err == nil {
						startedData.StartedAt = t
					}
				}
				if paramsVal, ok := dataMap["params"].(map[string]interface{}); ok {
					if createdAtMinVal, ok := paramsVal["created_at_min"].(string); ok {
						if t, err := time.Parse(time.RFC3339, createdAtMinVal); err == nil {
							startedData.Params.CreatedAtMin = &t
						}
					}
					if createdAtMaxVal, ok := paramsVal["created_at_max"].(string); ok {
						if t, err := time.Parse(time.RFC3339, createdAtMaxVal); err == nil {
							startedData.Params.CreatedAtMax = &t
						}
					}
					if statusVal, ok := paramsVal["status"].(string); ok {
						startedData.Params.Status = statusVal
					}
					if financialStatusVal, ok := paramsVal["financial_status"].(string); ok {
						startedData.Params.FinancialStatus = financialStatusVal
					}
					if fulfillmentStatusVal, ok := paramsVal["fulfillment_status"].(string); ok {
						startedData.Params.FulfillmentStatus = fulfillmentStatusVal
					}
				}
			} else if dataStruct, ok := lastStartedEvent.Data.(domain.SyncStartedEvent); ok {
				startedData = dataStruct
			}

			syncState = map[string]interface{}{
				"integration_id":   startedData.IntegrationID,
				"integration_type": startedData.IntegrationType,
				"status":           "in_progress",
				"started_at":       startedData.StartedAt,
				"params": map[string]interface{}{
					"created_at_min":     startedData.Params.CreatedAtMin,
					"created_at_max":     startedData.Params.CreatedAtMax,
					"status":             startedData.Params.Status,
					"financial_status":   startedData.Params.FinancialStatus,
					"fulfillment_status": startedData.Params.FulfillmentStatus,
				},
			}

			// Si no se pudo extraer, usar fallback
			if startedData.IntegrationID == 0 {
				syncState = map[string]interface{}{
					"integration_id": uint(integrationID),
					"status":         "in_progress",
					"started_at":     lastStartedEvent.Timestamp,
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"in_progress":    isInProgress,
		"sync_state":     syncState,
		"last_started":   lastStartedEvent != nil,
		"last_completed": lastCompletedEvent != nil,
		"last_failed":    lastFailedEvent != nil,
	})
}
