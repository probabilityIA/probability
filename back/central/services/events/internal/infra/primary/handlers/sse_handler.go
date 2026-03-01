package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// HandleSSE maneja la conexión SSE por business_id con filtros opcionales
func (h *SSEHandler) HandleSSE(c *gin.Context) {
	h.logger.Info(c.Request.Context()).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("remote_addr", c.ClientIP()).
		Msg("SSE endpoint llamado")

	if c.Request.Method == "OPTIONS" {
		h.setupSSEHeaders(c.Writer)
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	var businessID uint

	if businessIDStr := c.Param("businessID"); businessIDStr != "" {
		if id, parseErr := strconv.ParseUint(businessIDStr, 10, 32); parseErr == nil {
			businessID = uint(id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID de negocio inválido",
				"error":   parseErr.Error(),
			})
			return
		}
	} else if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		if id, parseErr := strconv.ParseUint(businessIDStr, 10, 32); parseErr == nil {
			businessID = uint(id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID de negocio inválido",
				"error":   parseErr.Error(),
			})
			return
		}
	}

	filter := h.buildFilterFromQuery(c)

	h.setupSSEHeaders(c.Writer)

	connectionID := h.eventManager.AddConnection(businessID, filter, c.Writer)

	// Precargar caché
	if businessID > 0 {
		h.preloadCacheEvents(c.Writer, businessID)
	}

	// Enviar mensaje de conexión
	message := fmt.Sprintf("Conexión SSE establecida para business %d", businessID)
	if businessID == 0 {
		message = "Conexión SSE establecida (super usuario - todos los businesses)"
	}
	connectionEvent := fmt.Sprintf("event: connection_established\ndata: {\"message\":\"%s\",\"connection_id\":\"%s\",\"timestamp\":\"%s\"}\n\n",
		message, connectionID, time.Now().Format(time.RFC3339))

	if _, err := c.Writer.Write([]byte(connectionEvent)); err != nil {
		h.logger.Error(c.Request.Context()).
			Err(err).
			Str("connection_id", connectionID).
			Msg("Error escribiendo mensaje de conexión SSE")
		return
	}

	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	h.keepConnectionAlive(c.Writer, connectionID, c.Request.Context())
}

// buildFilterFromQuery construye filtros desde query parameters
func (h *SSEHandler) buildFilterFromQuery(c *gin.Context) *entities.SSEConnectionFilter {
	filter := &entities.SSEConnectionFilter{}

	if integrationIDStr := c.Query("integration_id"); integrationIDStr != "" {
		if id, err := strconv.ParseUint(integrationIDStr, 10, 32); err == nil {
			integrationID := uint(id)
			filter.IntegrationID = &integrationID
		}
	}

	if eventTypesStr := c.Query("event_types"); eventTypesStr != "" {
		eventTypes := strings.Split(eventTypesStr, ",")
		filter.EventTypes = make([]string, 0, len(eventTypes))
		for _, et := range eventTypes {
			et = strings.TrimSpace(et)
			if et != "" {
				filter.EventTypes = append(filter.EventTypes, et)
			}
		}
	}

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
func (h *SSEHandler) setupSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control, Last-Event-ID")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Cache-Control")
	w.Header().Del("Access-Control-Allow-Credentials")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// keepConnectionAlive mantiene la conexión viva y detecta desconexiones
func (h *SSEHandler) keepConnectionAlive(w http.ResponseWriter, connectionID string, ctx context.Context) {
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
				Msg("Cliente SSE desconectado")
			return
		}
	}
}

// preloadCacheEvents precarga eventos del caché para nueva conexión
func (h *SSEHandler) preloadCacheEvents(w http.ResponseWriter, businessID uint) {
	events := h.eventManager.GetRecentEventsByBusiness(businessID, 0)

	if len(events) > 0 {
		h.logger.Info(context.Background()).
			Uint("business_id", businessID).
			Int("cache_events_count", len(events)).
			Msg("Precargando eventos del caché")

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
func (h *SSEHandler) eventToSSEJSON(event entities.Event) string {
	eventData := map[string]interface{}{
		"id":          event.ID,
		"type":        event.Type,
		"business_id": event.BusinessID,
		"timestamp":   event.Timestamp,
		"metadata":    event.Metadata,
	}

	if event.Data != nil {
		eventData["data"] = event.Data
	}

	jsonBytes, err := json.Marshal(eventData)
	if err != nil {
		return "{}"
	}

	return string(jsonBytes)
}

// sendSSEMessage envía un mensaje SSE formateado
func (h *SSEHandler) sendSSEMessage(w http.ResponseWriter, eventType, data string) {
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, data)
	w.Write([]byte(message))

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}
