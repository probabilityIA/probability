package events

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IntegrationEventManager implementa IIntegrationEventPublisher para manejar eventos de integraciones
type IntegrationEventManager struct {
	connectionManager *domain.ConnectionManager
	mutex             sync.RWMutex
	eventChan         chan domain.IntegrationEvent
	stopChan          chan struct{}

	// Estadísticas por business_id
	eventCount     map[uint]int
	eventTypeCount map[uint]map[domain.IntegrationEventType]int

	// Caché de eventos recientes por business_id para re-hidratación
	recentEvents      map[uint][]domain.IntegrationEvent
	maxRecent         int
	logger            log.ILogger
	connectionCounter uint64 // Contador para generar IDs únicos
}

// NewIntegrationEventManager crea un nuevo manager de eventos de integraciones
func New(logger log.ILogger) domain.IIntegrationEventPublisher {
	manager := &IntegrationEventManager{
		connectionManager: domain.NewConnectionManager(),
		eventChan:         make(chan domain.IntegrationEvent, 1000),
		stopChan:          make(chan struct{}),
		eventCount:        make(map[uint]int),
		eventTypeCount:    make(map[uint]map[domain.IntegrationEventType]int),
		recentEvents:      make(map[uint][]domain.IntegrationEvent),
		maxRecent:         2000,
		logger:            logger,
		connectionCounter: 0,
	}

	// Iniciar worker para procesar eventos
	go manager.startEventWorker()

	return manager
}

// startEventWorker procesa eventos en background
func (m *IntegrationEventManager) startEventWorker() {
	for {
		select {
		case event := <-m.eventChan:
			// #region agent log
			if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				logData, _ := json.Marshal(map[string]interface{}{
					"location": "constructor.go:startEventWorker",
					"message":  "Event received from channel, broadcasting",
					"data": map[string]interface{}{
						"event_id":       event.ID,
						"event_type":     string(event.Type),
						"integration_id": event.IntegrationID,
						"business_id":    event.BusinessID,
					},
					"timestamp":    time.Now().UnixMilli(),
					"sessionId":    "debug-session",
					"runId":        "run1",
					"hypothesisId": "D",
				})
				f.WriteString(string(logData) + "\n")
				f.Close()
			}
			// #endregion
			m.broadcastEvent(event)
		case <-m.stopChan:
			return
		}
	}
}
