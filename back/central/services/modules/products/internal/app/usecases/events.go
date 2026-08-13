package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	EventChannelDataStarted   = "products.channel_data.started"
	EventChannelDataProgress  = "products.channel_data.progress"
	EventChannelDataCompleted = "products.channel_data.completed"
	EventChannelDataFailed    = "products.channel_data.failed"
)

type publicadorEventos struct {
	rabbit rabbitmq.IQueue
	logger log.ILogger
}

func (uc *UseCases) WithEvents(rabbit rabbitmq.IQueue, logger log.ILogger) *UseCases {
	uc.eventos = &publicadorEventos{rabbit: rabbit, logger: logger}
	return uc
}

func (uc *UseCases) emitir(ctx context.Context, businessID, integrationID uint, tipo string, data map[string]interface{}) {
	if uc.eventos == nil || uc.eventos.rabbit == nil {
		return
	}
	_ = rabbitmq.PublishEvent(ctx, uc.eventos.rabbit, rabbitmq.EventEnvelope{
		Type:          tipo,
		Category:      "products",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Data:          data,
	})
}
