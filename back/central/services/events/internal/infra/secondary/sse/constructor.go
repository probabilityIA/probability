package sse

import (
	"sync"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// EventManager implementa ISSEPublisher para manejar eventos SSE en tiempo real
type EventManager struct {
	connections map[string]*entities.SSEConnection
	mutex       sync.RWMutex
	eventChan   chan entities.Event
	stopChan    chan struct{}

	eventCount     map[uint]int
	eventTypeCount map[uint]map[string]int

	recentEvents      map[uint][]entities.Event
	maxRecent         int
	logger            log.ILogger
	connectionCounter uint64
}

// New crea un nuevo EventManager
func New(logger log.ILogger) ports.ISSEPublisher {
	manager := &EventManager{
		connections:       make(map[string]*entities.SSEConnection),
		eventChan:         make(chan entities.Event, 1000),
		stopChan:          make(chan struct{}),
		eventCount:        make(map[uint]int),
		eventTypeCount:    make(map[uint]map[string]int),
		recentEvents:      make(map[uint][]entities.Event),
		maxRecent:         2000,
		logger:            logger,
		connectionCounter: 0,
	}

	go manager.startEventWorker()

	return manager
}
