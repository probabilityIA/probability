package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (p *responsePublisher) PublishResponse(ctx context.Context, phoneNumber string, businessID uint, text string) error {
	dto := domain.AIResponseDTO{
		PhoneNumber:  phoneNumber,
		ResponseText: text,
		BusinessID:   businessID,
		Timestamp:    time.Now().Unix(),
	}

	payload, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("error serializing AI response: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppAIResponse, payload); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("phone", phoneNumber).
			Msg("Error publicando respuesta AI a WhatsApp")
		return fmt.Errorf("error publishing AI response: %w", err)
	}

	p.log.Info(ctx).
		Str("phone", phoneNumber).
		Uint("business_id", businessID).
		Msg("Respuesta AI publicada a whatsapp.ai.response")

	return nil
}
