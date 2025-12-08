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
	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// SSEHandler maneja las conexiones Server-Sent Events adaptado a Gin
type SSEHandler struct {
	eventManager domain.IEventPublisher
	logger       log.ILogger
}

// SSEHandlerInterface define la interfaz del handler SSE
type SSEHandlerInterface interface {
	HandleSSE(c *gin.Context)
	GetManager() domain.IEventPublisher
}

// NewSSEHandler crea un nuevo handler de SSE
func New(eventManager domain.IEventPublisher, logger log.ILogger) SSEHandlerInterface {
	return &SSEHandler{
		eventManager: eventManager,
		logger:       logger,
	}
}

// GetManager retorna el manager interno para acceso externo
func (h *SSEHandler) GetManager() domain.IEventPublisher {
	return h.eventManager
}

// HandleSSE maneja la conexión SSE por business_id con filtros opcionales (adaptado a Gin)
func (h *SSEHandler) HandleSSE(c *gin.Context) {
	// Obtener business_id de los parámetros de la URL o query params
	var businessID uint

	// Intentar obtener de parámetro de ruta primero
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
		// Si no está en la ruta, intentar desde query params
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
	} else {
		// Si no se proporciona business_id, asumir super usuario (business_id = 0)
		businessID = 0
	}

	// Construir filtros desde query params
	filter := h.buildFilterFromQuery(c)

	h.logger.Info(c.Request.Context()).
		Uint("business_id", businessID).
		Interface("filter", filter).
		Bool("is_super_user", businessID == 0).
		Msg("Nueva conexión SSE solicitada")

	// Configurar headers para SSE
	h.setupSSEHeaders(c.Writer)

	// Agregar la conexión al manager (retorna connectionID)
	connectionID := h.eventManager.AddConnection(businessID, filter, c.Writer)

	// PRECARGAR CACHÉ: Enviar eventos históricos si existen
	if businessID > 0 {
		h.preloadCacheEvents(c.Writer, businessID)
	} else {
		// Para super usuario, precargar eventos de todos los businesses
		h.preloadCacheEventsForSuperUser(c.Writer)
	}

	// Enviar mensaje de conexión establecida
	message := fmt.Sprintf("Conexión SSE establecida para business %d", businessID)
	if businessID == 0 {
		message = "Conexión SSE establecida (super usuario - todos los businesses)"
	}
	connectionEvent := fmt.Sprintf("event: connection_established\ndata: {\"message\":\"%s\",\"connection_id\":\"%s\",\"timestamp\":\"%s\"}\n\n",
		message, connectionID, time.Now().Format(time.RFC3339))

	c.Writer.Write([]byte(connectionEvent))
	c.Writer.Flush()

	h.logger.Info(c.Request.Context()).
		Uint("business_id", businessID).
		Str("connection_id", connectionID).
		Msg("Conexión SSE establecida y mensaje de confirmación enviado")

	// Mantener la conexión viva y detectar desconexión
	h.keepConnectionAlive(c.Writer, connectionID, c.Request.Context())
}

// buildFilterFromQuery construye filtros desde los query parameters
func (h *SSEHandler) buildFilterFromQuery(c *gin.Context) *domain.SSEConnectionFilter {
	filter := &domain.SSEConnectionFilter{}

	// Filtro por integration_id
	if integrationIDStr := c.Query("integration_id"); integrationIDStr != "" {
		if id, err := strconv.ParseUint(integrationIDStr, 10, 32); err == nil {
			integrationID := uint(id)
			filter.IntegrationID = &integrationID
		}
	}

	// Filtro por event_types (separados por comas)
	if eventTypesStr := c.Query("event_types"); eventTypesStr != "" {
		eventTypes := strings.Split(eventTypesStr, ",")
		filter.EventTypes = make([]domain.EventType, 0, len(eventTypes))
		for _, et := range eventTypes {
			et = strings.TrimSpace(et)
			if et != "" {
				filter.EventTypes = append(filter.EventTypes, domain.EventType(et))
			}
		}
	}

	// Filtro por order_ids (separados por comas)
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

// setupSSEHeaders configura los headers HTTP para Server-Sent Events
func (h *SSEHandler) setupSSEHeaders(w http.ResponseWriter) {
	// Headers básicos para SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Headers anti-buffering para nginx proxy
	w.Header().Set("X-Accel-Buffering", "no") // Deshabilita buffering en nginx
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Headers adicionales para compatibilidad con proxies
	w.Header().Set("Transfer-Encoding", "chunked")      // Fuerza streaming
	w.Header().Set("X-Content-Type-Options", "nosniff") // Previene buffering del browser
}

// keepConnectionAlive mantiene la conexión viva y detecta desconexiones
func (h *SSEHandler) keepConnectionAlive(w http.ResponseWriter, connectionID string, ctx context.Context) {
	// Crear un ticker para enviar keep-alive cada 30 segundos
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Detectar si el cliente se desconectó usando Request.Context
	done := ctx.Done()

	for {
		select {
		case <-ticker.C:
			// Enviar comentario de keep-alive
			h.sendSSEMessage(w, "keep-alive", "ping")
		case <-done:
			// Cliente se desconectó
			h.eventManager.RemoveConnection(connectionID)
			h.logger.Info(ctx).
				Str("connection_id", connectionID).
				Msg("Cliente SSE desconectado")
			return
		}
	}
}

// preloadCacheEvents precarga eventos del caché para una nueva conexión SSE por business_id
func (h *SSEHandler) preloadCacheEvents(w http.ResponseWriter, businessID uint) {
	// Obtener eventos del caché usando type assertion
	type recentGetter interface {
		GetRecentEventsByBusiness(uint, int64) []domain.Event
	}

	if rg, ok := h.eventManager.(recentGetter); ok {
		// Obtener todos los eventos desde el inicio (since_seq=0)
		events := rg.GetRecentEventsByBusiness(businessID, 0)

		if len(events) > 0 {
			h.logger.Info(context.Background()).
				Uint("business_id", businessID).
				Int("cache_events_count", len(events)).
				Msg("Precargando eventos del caché")

			// Enviar cada evento del caché como mensaje SSE
			for _, event := range events {
				// Convertir evento a JSON para enviar por SSE
				eventJSON := h.eventToSSEJSON(event)
				message := fmt.Sprintf("event: %s\ndata: %s\n\n", event.Type, eventJSON)
				w.Write([]byte(message))
			}

			// Flush para asegurar que se envíen inmediatamente
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			h.logger.Info(context.Background()).
				Uint("business_id", businessID).
				Int("events_sent", len(events)).
				Msg("Eventos del caché enviados por SSE")
		}
	}
}

// preloadCacheEventsForSuperUser precarga eventos de todos los businesses para super usuario
func (h *SSEHandler) preloadCacheEventsForSuperUser(w http.ResponseWriter) {
	// Para super usuario, no precargamos eventos por defecto
	// O se puede implementar lógica para obtener eventos de todos los businesses
	// Por ahora, solo enviamos mensaje de confirmación
	h.logger.Info(context.Background()).
		Msg("Super usuario conectado - eventos se enviarán en tiempo real")
}

// eventToSSEJSON convierte un evento a JSON para enviar por SSE
func (h *SSEHandler) eventToSSEJSON(event domain.Event) string {
	// Usar el mismo formato que el manager de eventos
	eventData := map[string]interface{}{
		"id":          event.ID,
		"type":        event.Type,
		"business_id": event.BusinessID,
		"timestamp":   event.Timestamp,
		"metadata":    event.Metadata,
	}

	// Agregar campos específicos según el tipo de evento
	if event.Data != nil {
		if dataMap, ok := event.Data.(map[string]interface{}); ok {
			for key, value := range dataMap {
				eventData[key] = value
			}
		}
	}

	// Convertir a JSON
	jsonBytes, err := json.Marshal(eventData)
	if err != nil {
		h.logger.Error(context.Background()).
			Err(err).
			Msg("Error serializando evento para SSE")
		return "{}"
	}

	return string(jsonBytes)
}

// orderEventToSSEJSON convierte un evento de orden a JSON para enviar por SSE
func (h *SSEHandler) orderEventToSSEJSON(event *domain.OrderEvent) string {
	eventData := map[string]interface{}{
		"id":             event.ID,
		"type":           event.Type,
		"order_id":       event.OrderID,
		"business_id":    event.BusinessID,
		"integration_id": event.IntegrationID,
		"timestamp":      event.Timestamp,
		"data":           event.Data,
		"metadata":       event.Metadata,
	}

	jsonBytes, err := json.Marshal(eventData)
	if err != nil {
		h.logger.Error(context.Background()).
			Err(err).
			Msg("Error serializando evento de orden para SSE")
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
