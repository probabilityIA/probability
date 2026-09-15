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
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

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

	businessID, ok := h.resolveSSEBusinessID(c)
	if !ok {
		return
	}

	filter := h.buildFilterFromQuery(c)

	h.setupSSEHeaders(c.Writer)

	connectionID := h.eventManager.AddConnection(businessID, filter, c.Writer)

	if businessID > 0 {
		if lastEventID := c.GetHeader("Last-Event-ID"); lastEventID != "" {
			sinceSeq, _ := strconv.ParseInt(lastEventID, 10, 64)
			h.preloadCacheEventsSince(c.Writer, businessID, sinceSeq, c.Request.Context())
		}
	}

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

func (h *SSEHandler) resolveSSEBusinessID(c *gin.Context) (uint, bool) {
	tokenBusinessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "UNAUTHORIZED",
			"message": "contexto de negocio no encontrado",
		})
		return 0, false
	}
	if tokenBusinessID > 0 {
		return tokenBusinessID, true
	}

	raw := c.Param("businessID")
	if raw == "" {
		raw = c.Query("business_id")
	}
	if raw == "" {
		return 0, true
	}

	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		h.logger.Warn(c.Request.Context()).
			Err(err).
			Str("business_id_raw", raw).
			Msg("ID de negocio invalido en SSE")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_BUSINESS_ID",
			"message": "El ID de negocio proporcionado no es valido",
		})
		return 0, false
	}
	return uint(id), true
}

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

func (h *SSEHandler) setupSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Cache-Control")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (h *SSEHandler) keepConnectionAlive(w http.ResponseWriter, connectionID string, ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	done := ctx.Done()
	flusher, hasFlusher := w.(http.Flusher)

	for {
		select {
		case <-ticker.C:
			if err := h.sendSSEMessage(w, "keep-alive", "ping"); err != nil {
				h.eventManager.RemoveConnection(connectionID)
				h.logger.Info(ctx).
					Str("connection_id", connectionID).
					Msg("Cliente SSE desconectado (error en keep-alive)")
				return
			}
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

func (h *SSEHandler) preloadCacheEventsSince(w http.ResponseWriter, businessID uint, sinceSeq int64, ctx context.Context) {
	events := h.eventManager.GetRecentEventsByBusiness(businessID, sinceSeq)

	if len(events) > 0 {
		h.logger.Info(ctx).
			Uint("business_id", businessID).
			Int64("since_seq", sinceSeq).
			Int("cache_events_count", len(events)).
			Msg("Precargando eventos del caché (reconexión)")

		for _, event := range events {
			eventJSON := h.eventToSSEJSON(event)
			idLine := ""
			if seqVal, ok := event.Metadata["sse_seq"]; ok {
				idLine = fmt.Sprintf("id: %v\n", seqVal)
			}
			message := fmt.Sprintf("%sevent: %s\ndata: %s\n\n", idLine, event.Type, eventJSON)
			if _, err := w.Write([]byte(message)); err != nil {
				h.logger.Warn(ctx).
					Err(err).
					Uint("business_id", businessID).
					Msg("Error escribiendo eventos de cache, cliente desconectado")
				return
			}
		}

		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

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
		h.logger.Error(context.Background()).
			Err(err).
			Str("event_id", event.ID).
			Str("event_type", event.Type).
			Msg("Error serializando evento SSE a JSON")
		return "{}"
	}

	return string(jsonBytes)
}

func (h *SSEHandler) sendSSEMessage(w http.ResponseWriter, eventType, data string) error {
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, data)
	if _, err := w.Write([]byte(message)); err != nil {
		return err
	}

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
