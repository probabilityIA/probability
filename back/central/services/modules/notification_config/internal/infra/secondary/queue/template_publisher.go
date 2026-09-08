package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/templates"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type templatePublisher struct {
	rabbit rabbitmq.IQueue
	logger log.ILogger
}

func NewTemplatePublisher(rabbit rabbitmq.IQueue, logger log.ILogger) templates.ISubmissionPublisher {
	return &templatePublisher{
		rabbit: rabbit,
		logger: logger.WithModule("whatsapp_template_publisher"),
	}
}

func (p *templatePublisher) PublishTemplateSubmission(ctx context.Context, message dtos.TemplateSubmissionMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error serializando la plantilla: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppTemplateSubmitRequests, payload); err != nil {
		p.logger.Error(ctx).Err(err).
			Uint("template_id", message.TemplateID).
			Msg("Error publicando la plantilla a la cola de envio a Meta")
		return err
	}

	p.logger.Info(ctx).
		Uint("template_id", message.TemplateID).
		Str("name", message.Name).
		Msg("Plantilla publicada a whatsapp.templates.submit.requests")

	return nil
}
