package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/secondary/client/response"
)

// ProcessWebhookOrder procesa una orden recibida por webhook de WooCommerce.
// El topic indica el evento: order.created, order.updated, order.deleted, order.restored.
// El storeURL proviene del header X-WC-Webhook-Source y se usa para identificar la integración.
func (uc *wooCommerceUseCase) ProcessWebhookOrder(ctx context.Context, topic string, storeURL string, rawBody []byte) error {
	// order.deleted — solo logear, no procesar
	if topic == "order.deleted" {
		uc.logger.Info(ctx).
			Str("topic", topic).
			Str("store_url", storeURL).
			Msg("WooCommerce order deleted event received, skipping")
		return nil
	}

	// 1. Deserializar payload
	var orderResp response.WooOrderResponse
	if err := json.Unmarshal(rawBody, &orderResp); err != nil {
		return fmt.Errorf("deserializing webhook payload: %w", err)
	}

	order := orderResp.ToDomain()

	uc.logger.Info(ctx).
		Str("topic", topic).
		Str("store_url", storeURL).
		Str("order_number", order.Number).
		Int64("order_id", order.ID).
		Msg("Processing WooCommerce webhook order")

	// 2. Mapear a DTO canónico
	dto := mapper.MapWooOrderToProbability(&order, rawBody)

	// 3. Publicar a la cola
	if err := uc.publisher.Publish(ctx, dto); err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("order_number", order.Number).
			Str("topic", topic).
			Msg("Failed to publish webhook order")
		return fmt.Errorf("publishing webhook order: %w", err)
	}

	uc.logger.Info(ctx).
		Str("order_number", order.Number).
		Str("topic", topic).
		Msg("WooCommerce webhook order published successfully")

	return nil
}
