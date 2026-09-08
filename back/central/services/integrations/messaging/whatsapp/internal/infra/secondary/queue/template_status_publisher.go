package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasetemplates"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type templateStatusPublisher struct {
	rabbit rabbitmq.IQueue
	logger log.ILogger
}

func NewTemplateStatusPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) usecasetemplates.IStatusPublisher {
	return &templateStatusPublisher{
		rabbit: rabbit,
		logger: logger.WithModule("whatsapp-template-status-publisher"),
	}
}

func (p *templateStatusPublisher) PublishTemplateStatus(ctx context.Context, result usecasetemplates.CustomTemplateResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("error serializando el estado de la plantilla: %w", err)
	}

	return p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppTemplateSubmitResults, payload)
}
