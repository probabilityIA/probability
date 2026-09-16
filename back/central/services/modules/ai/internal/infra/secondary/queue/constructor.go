package queue

import (
	"sync"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Recorder struct {
	queue    rabbitmq.IQueue
	declared sync.Once
}

func New(queue rabbitmq.IQueue) ports.IConversationRecorder {
	if queue == nil {
		return nil
	}
	return &Recorder{queue: queue}
}
