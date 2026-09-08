package template_result_consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/templates"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Consumer struct {
	rabbit  rabbitmq.IQueue
	useCase templates.IUseCase
	logger  log.ILogger
}

func New(rabbit rabbitmq.IQueue, useCase templates.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		rabbit:  rabbit,
		useCase: useCase,
		logger:  logger.WithModule("template_result_consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.rabbit.DeclareQueue(rabbitmq.QueueWhatsAppTemplateSubmitResults, true); err != nil {
		c.logger.Error(ctx).Err(err).Str("queue", rabbitmq.QueueWhatsAppTemplateSubmitResults).
			Msg("Error declarando la cola")
		return err
	}

	return c.rabbit.Consume(ctx, rabbitmq.QueueWhatsAppTemplateSubmitResults, func(body []byte) error {
		var result dtos.TemplateSubmissionResult

		if err := json.Unmarshal(body, &result); err != nil {
			c.logger.Warn(ctx).Err(err).
				Msg("Mensaje de resultado de plantilla ilegible - se descarta (ACK)")
			return nil
		}

		if result.TemplateID == 0 {
			if result.WABAID == "" || result.Name == "" {
				c.logger.Warn(ctx).
					Msg("Resultado de plantilla sin template_id ni waba/nombre - se descarta (ACK)")
				return nil
			}

			if err := c.useCase.ApplyMetaStatus(ctx, result.WABAID, result.Name, result.Language, result.Status, result.Reason); err != nil {
				c.logger.Error(ctx).Err(err).
					Str("waba_id", result.WABAID).
					Str("name", result.Name).
					Msg("Error aplicando el estado de plantilla del webhook - se reintenta")
				return err
			}

			return nil
		}

		if err := c.useCase.ApplySubmissionResult(ctx, result); err != nil {
			c.logger.Error(ctx).Err(err).
				Uint("template_id", result.TemplateID).
				Msg("Error aplicando el resultado de envio de plantilla - se reintenta")
			return err
		}

		return nil
	})
}
