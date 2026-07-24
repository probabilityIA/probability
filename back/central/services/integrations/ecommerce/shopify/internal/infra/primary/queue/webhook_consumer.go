package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/mappers"
	clientresponse "github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/response"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type WebhookMessage struct {
	Topic      string          `json:"topic"`
	ShopDomain string          `json:"shop_domain"`
	IsTest     bool            `json:"is_test"`
	Body       json.RawMessage `json:"body"`
}

type WebhookConsumer struct {
	queue   rabbitmq.IQueue
	useCase usecases.IShopifyUseCase
	logger  log.ILogger
}

func NewWebhookConsumer(queue rabbitmq.IQueue, useCase usecases.IShopifyUseCase, logger log.ILogger) *WebhookConsumer {
	return &WebhookConsumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("shopify.webhook_consumer"),
	}
}

func (c *WebhookConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueWebhooksShopifyReceived, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de webhooks Shopify")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueWebhooksShopifyReceived, func(body []byte) error {
			var msg WebhookMessage
			if err := json.Unmarshal(body, &msg); err != nil {
				c.logger.Error(ctx).Err(err).Msg("Mensaje de webhook Shopify invalido, descartado")
				return nil
			}
			DispatchWebhook(context.Background(), c.useCase, c.logger, msg.Topic, msg.ShopDomain, msg.Body, msg.IsTest)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de webhooks Shopify")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de webhooks Shopify iniciado")
}

func DispatchWebhook(ctx context.Context, useCase usecases.IShopifyUseCase, logger log.ILogger, topic string, shopDomain string, bodyBytes []byte, isTest bool) {
	var orderResp clientresponse.Order
	if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
		logger.Error(ctx).
			Err(err).
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Msg("Error al parsear payload de Shopify a Order")
		return
	}

	mapped := mappers.MapOrderResponseToShopifyOrder(orderResp, bodyBytes, nil, 0, "shopify")
	shopifyOrder := &mapped

	var err error
	switch topic {
	case "orders/create":
		err = useCase.CreateOrder(ctx, shopDomain, shopifyOrder, bodyBytes, isTest)
	case "orders/paid":
		err = useCase.ProcessOrderPaid(ctx, shopDomain, shopifyOrder)
	case "orders/updated":
		err = useCase.ProcessOrderUpdated(ctx, shopDomain, shopifyOrder)
	case "orders/cancelled":
		err = useCase.ProcessOrderCancelled(ctx, shopDomain, shopifyOrder)
	case "orders/fulfilled":
		err = useCase.ProcessOrderFulfilled(ctx, shopDomain, shopifyOrder)
	case "orders/partially_fulfilled":
		err = useCase.ProcessOrderPartiallyFulfilled(ctx, shopDomain, shopifyOrder)
	default:
		logger.Info(ctx).
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Msg("Topic no manejado, ignorando webhook")
		return
	}

	if err != nil {
		logger.Error(ctx).
			Err(err).
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("Error al procesar webhook de Shopify")
	} else {
		logger.Info(ctx).
			Str("topic", topic).
			Str("shop_domain", shopDomain).
			Str("order_id", orderResp.Name).
			Msg("Webhook procesado exitosamente")
	}
}
