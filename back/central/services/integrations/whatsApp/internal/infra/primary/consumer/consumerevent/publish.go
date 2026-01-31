package consumerevent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/consumer/consumerevent/request"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/repository"
)

// publishConfirmationRequest publica un evento de confirmación a RabbitMQ
func (c *consumer) publishConfirmationRequest(
	ctx context.Context,
	order *request.OrderData,
	config *repository.NotificationConfigData, // ← Usa repositorio
) error {
	event := map[string]interface{}{
		"event_type":     "order.confirmation_requested",
		"order_id":       order.ID,
		"order_number":   order.OrderNumber,
		"customer_phone": order.CustomerPhone,
		"total_amount":   order.TotalAmount,
		"currency":       order.Currency,
		"template_name":  config.TemplateName, // ← Campo directo del DTO
		"language":       config.Language,
		"recipient_type": config.RecipientType,
		"business_id":    order.BusinessID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	return c.rabbitMQ.Publish(ctx, "orders.confirmation.requested", payload)
}
