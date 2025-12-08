package events

import (
	"sync"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// EventManager implementa EventManagerPort para manejar eventos en tiempo real
type EventManager struct {
	// Conexiones por connectionID (permite filtros por business_id)
	connections map[string]*domain.SSEConnection // connectionID -> connection
	mutex       sync.RWMutex
	eventChan   chan domain.Event
	stopChan    chan struct{}

	// Estadísticas por business_id
	eventCount     map[uint]int
	eventTypeCount map[uint]map[domain.EventType]int

	// Caché de eventos recientes por business_id para re-hidratación
	recentEvents      map[uint][]domain.Event
	maxRecent         int
	logger            log.ILogger
	connectionCounter uint64 // Contador para generar IDs únicos
}

// NewEventManager crea un nuevo manager de eventos
func New(logger log.ILogger) domain.IEventPublisher {
	manager := &EventManager{
		connections:       make(map[string]*domain.SSEConnection),
		eventChan:         make(chan domain.Event, 1000),
		stopChan:          make(chan struct{}),
		eventCount:        make(map[uint]int),
		eventTypeCount:    make(map[uint]map[domain.EventType]int),
		recentEvents:      make(map[uint][]domain.Event),
		maxRecent:         2000,
		logger:            logger,
		connectionCounter: 0,
	}

	// Iniciar worker para procesar eventos
	go manager.startEventWorker()

	return manager
}
