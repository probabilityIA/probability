package server

import (
	"sync"
)

// QueueRegistry mantiene un registro de todas las colas declaradas
type QueueRegistry struct {
	queues []string
	mu     sync.Mutex
}

// NewQueueRegistry crea un nuevo registro de colas
func NewQueueRegistry() *QueueRegistry {
	return &QueueRegistry{
		queues: make([]string, 0),
	}
}

// Register registra una cola como declarada exitosamente
func (r *QueueRegistry) Register(queueName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.queues = append(r.queues, queueName)
}

// GetQueues retorna todas las colas registradas
func (r *QueueRegistry) GetQueues() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]string, len(r.queues))
	copy(result, r.queues)
	return result
}
