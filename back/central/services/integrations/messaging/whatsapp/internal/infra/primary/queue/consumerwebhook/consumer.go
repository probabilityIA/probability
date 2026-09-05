package consumerwebhook

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasetemplates"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Consumer struct {
	queue     rabbitmq.IQueue
	useCase   usecasemessaging.IUseCase
	templates usecasetemplates.IUseCase
	log       log.ILogger
}

func New(queue rabbitmq.IQueue, useCase usecasemessaging.IUseCase, templates usecasetemplates.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		queue:     queue,
		useCase:   useCase,
		templates: templates,
		log:       logger.WithModule("whatsapp.webhook_consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueWebhooksWhatsappReceived, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al declarar la cola de webhooks WhatsApp")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueWebhooksWhatsappReceived, func(body []byte) error {
			var webhook request.WebhookPayload
			if err := json.Unmarshal(body, &webhook); err != nil {
				c.log.Error(ctx).Err(err).Msg("Mensaje de webhook WhatsApp invalido, descartado")
				return nil
			}
			Dispatch(context.Background(), c.useCase, c.templates, c.log, webhook)
			return nil
		})
		if err != nil {
			c.log.Error(ctx).Err(err).Msg("Error al consumir la cola de webhooks WhatsApp")
		}
	}()

	c.log.Info(ctx).Msg("Consumer de webhooks WhatsApp iniciado")
}

func Dispatch(ctx context.Context, useCase usecasemessaging.IUseCase, templates usecasetemplates.IUseCase, logger log.ILogger, webhook request.WebhookPayload) {
	webhookDTO := mappers.WebhookPayloadToDomain(webhook)

	for _, entry := range webhookDTO.Entry {
		for _, change := range entry.Changes {
			switch change.Field {
			case "messages":
				if len(change.Value.Messages) > 0 {
					if err := useCase.HandleIncomingMessage(ctx, webhookDTO); err != nil {
						logger.Error(ctx).Err(err).Msg("Error procesando mensajes entrantes de WhatsApp")
					}
				}
				if len(change.Value.Statuses) > 0 {
					if err := useCase.HandleMessageStatus(ctx, webhookDTO); err != nil {
						logger.Error(ctx).Err(err).Msg("Error procesando estados de mensajes de WhatsApp")
					}
				}
			case "message_template_status_update":
				logger.Info(ctx).
					Str("waba_id", entry.ID).
					Str("template", change.Value.TemplateName).
					Str("event", change.Value.TemplateEvent).
					Msg("Actualizacion de estado de plantilla recibida")
				if templates == nil {
					continue
				}
				if err := templates.HandleStatusUpdate(
					ctx,
					entry.ID,
					change.Value.TemplateName,
					change.Value.TemplateLanguage,
					change.Value.TemplateEvent,
					change.Value.TemplateReason,
				); err != nil {
					logger.Error(ctx).Err(err).Msg("Error actualizando estado de plantilla")
				}
			default:
				logger.Warn(ctx).Str("field", change.Field).Msg("Campo de webhook WhatsApp no reconocido")
			}
		}
	}
}
