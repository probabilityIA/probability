package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type whatsAppTemplateRequest struct {
	BusinessID   uint     `json:"business_id"`
	Phone        string   `json:"phone"`
	TemplateName string   `json:"template_name"`
	Language     string   `json:"language"`
	Parameters   []string `json:"parameters"`
}

type WhatsAppTemplateRequester struct {
	rabbit rabbitmq.IQueue
}

func NewWhatsAppTemplateRequester(rabbit rabbitmq.IQueue) ports.IWhatsAppTemplateRequester {
	return &WhatsAppTemplateRequester{rabbit: rabbit}
}

func (r *WhatsAppTemplateRequester) RequestTemplate(ctx context.Context, businessID uint, phone, templateName string, parameters []string) error {
	if r.rabbit == nil {
		return nil
	}
	body, err := json.Marshal(whatsAppTemplateRequest{
		BusinessID:   businessID,
		Phone:        phone,
		TemplateName: templateName,
		Language:     "es",
		Parameters:   parameters,
	})
	if err != nil {
		return fmt.Errorf("serializar solicitud de plantilla: %w", err)
	}
	return r.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppFlowSends, body)
}
