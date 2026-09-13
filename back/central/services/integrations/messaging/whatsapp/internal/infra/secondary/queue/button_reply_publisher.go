package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type buttonReplyMessage struct {
	BusinessID       uint   `json:"business_id"`
	PhoneNumber      string `json:"phone_number"`
	ButtonText       string `json:"button_text"`
	ContextMessageID string `json:"context_message_id"`
	MessageID        string `json:"message_id"`
}

type buttonReplyPublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

func NewButtonReplyPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) ports.IButtonReplyPublisher {
	return &buttonReplyPublisher{
		rabbit: rabbit,
		log:    logger.WithModule("whatsapp-button-reply-publisher"),
	}
}

func (p *buttonReplyPublisher) PublishButtonReply(ctx context.Context, event ports.ButtonReplyEvent) error {
	payload, err := json.Marshal(buttonReplyMessage{
		BusinessID:       event.BusinessID,
		PhoneNumber:      event.PhoneNumber,
		ButtonText:       event.ButtonText,
		ContextMessageID: event.ContextMessageID,
		MessageID:        event.MessageID,
	})
	if err != nil {
		return fmt.Errorf("error serializando la respuesta de boton: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppButtonReplies, payload); err != nil {
		p.log.Error(ctx).Err(err).
			Str("button_text", event.ButtonText).
			Msg("[WhatsApp Publisher] - error publicando la respuesta de boton")
		return err
	}

	return nil
}
