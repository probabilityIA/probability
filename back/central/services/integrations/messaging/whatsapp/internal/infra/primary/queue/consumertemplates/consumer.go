package consumertemplates

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasetemplates"
	whaErrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Consumer struct {
	rabbit  rabbitmq.IQueue
	useCase usecasetemplates.IUseCase
	logger  log.ILogger
}

func New(rabbit rabbitmq.IQueue, useCase usecasetemplates.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		rabbit:  rabbit,
		useCase: useCase,
		logger:  logger.WithModule("whatsapp-template-submit-consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueWhatsAppTemplateSubmitRequests, true); err != nil {
		c.logger.Error(ctx).Err(err).Str("queue", rabbitmq.QueueWhatsAppTemplateSubmitRequests).
			Msg("Error declarando la cola")
		return err
	}

	return c.rabbit.Consume(ctx, rabbitmq.QueueWhatsAppTemplateSubmitRequests, func(body []byte) error {
		var submission usecasetemplates.CustomTemplateSubmission

		if err := json.Unmarshal(body, &submission); err != nil {
			c.logger.Warn(ctx).Err(err).
				Msg("Solicitud de plantilla ilegible - se descarta (ACK)")
			return nil
		}

		if submission.Name == "" || (submission.Action != "delete" && submission.TemplateID == 0) {
			c.logger.Warn(ctx).
				Uint("template_id", submission.TemplateID).
				Msg("Solicitud de plantilla incompleta - se descarta (ACK)")
			return nil
		}

		if submission.Action == "delete" {
			if err := c.useCase.DeleteCustom(ctx, submission); err != nil {
				if whaErrors.IsNonRetryable(err) {
					c.logger.Warn(ctx).Err(err).
						Str("name", submission.Name).
						Msg("No se pudo borrar la plantilla en Meta - error permanente (ACK)")
					return nil
				}
				return err
			}
			return nil
		}

		result := c.useCase.SubmitCustom(ctx, submission)

		payload, err := json.Marshal(result)
		if err != nil {
			c.logger.Error(ctx).Err(err).
				Uint("template_id", submission.TemplateID).
				Msg("Error serializando el resultado de la plantilla")
			return nil
		}

		if err := c.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppTemplateSubmitResults, payload); err != nil {
			c.logger.Error(ctx).Err(err).
				Uint("template_id", submission.TemplateID).
				Msg("Error publicando el resultado de la plantilla - se reintenta")
			return err
		}

		return nil
	})
}
