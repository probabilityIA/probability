package queue

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const eventoCambioEstado = "order.status_changed"

type orderEventMessage struct {
	EventType string         `json:"event_type"`
	OrderID   string         `json:"order_id"`
	Timestamp time.Time      `json:"timestamp"`
	Order     *orderSnapshot `json:"order"`
	Changes   map[string]any `json:"changes,omitempty"`
}

type orderSnapshot struct {
	Platform       string `json:"platform"`
	Status         string `json:"status"`
	ExternalID     string `json:"external_id"`
	IntegrationID  uint   `json:"integration_id"`
	TrackingNumber string `json:"tracking_number"`
	TrackingURL    string `json:"tracking_url"`
	Carrier        string `json:"carrier"`
	ShippingCity   string `json:"shipping_city"`
}

type OrderStatusConsumer struct {
	queue  rabbitmq.IQueue
	uc     usecases.ITiendanubeUseCase
	logger log.ILogger
}

func NewOrderStatusConsumer(queue rabbitmq.IQueue, uc usecases.ITiendanubeUseCase, logger log.ILogger) *OrderStatusConsumer {
	return &OrderStatusConsumer{
		queue:  queue,
		uc:     uc,
		logger: logger.WithModule("tiendanube.status_consumer"),
	}
}

func (c *OrderStatusConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ no disponible, el push-back de estados a Tiendanube queda apagado")
		return
	}

	if err := c.queue.DeclareExchange(rabbitmq.ExchangeOrderEvents, "fanout", true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("No se pudo declarar el exchange de eventos de ordenes")
		return
	}
	if err := c.queue.DeclareQueue(rabbitmq.QueueOrdersToTiendanube, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("No se pudo declarar la cola de ordenes hacia Tiendanube")
		return
	}
	if err := c.queue.BindQueue(rabbitmq.QueueOrdersToTiendanube, rabbitmq.ExchangeOrderEvents, ""); err != nil {
		c.logger.Error(ctx).Err(err).Msg("No se pudo enlazar la cola de ordenes hacia Tiendanube")
		return
	}

	c.logger.Info(ctx).Str("queue", rabbitmq.QueueOrdersToTiendanube).Msg("Consumidor de push-back de estados a Tiendanube iniciado")

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueOrdersToTiendanube, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("El consumidor de estados de Tiendanube se detuvo con error")
		}
	}()
}

func (c *OrderStatusConsumer) handleMessage(ctx context.Context, body []byte) {
	var msg orderEventMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("No se pudo leer el evento de orden")
		return
	}

	if msg.EventType != eventoCambioEstado {
		return
	}
	if msg.Order == nil || !esTiendanube(msg.Order.Platform) {
		return
	}
	if msg.Order.ExternalID == "" || msg.Order.IntegrationID == 0 {
		return
	}

	estado, _ := msg.Changes["current_status"].(string)
	if estado == "" {
		estado = msg.Order.Status
	}
	if estado == "" {
		return
	}

	integrationID := strconv.FormatUint(uint64(msg.Order.IntegrationID), 10)
	update := domain.OrderStatusUpdate{
		ExternalOrderID: msg.Order.ExternalID,
		Status:          estado,
		TrackingNumber:  msg.Order.TrackingNumber,
		TrackingURL:     msg.Order.TrackingURL,
		Carrier:         msg.Order.Carrier,
		City:            msg.Order.ShippingCity,
		HappenedAt:      msg.Timestamp,
	}

	if err := c.uc.UpdateOrderStatus(ctx, integrationID, update); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("order_id", msg.OrderID).
			Str("external_id", msg.Order.ExternalID).
			Str("status", estado).
			Msg("No se pudo empujar el estado de la orden a Tiendanube")
		return
	}

	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Str("external_id", msg.Order.ExternalID).
		Str("status", estado).
		Msg("Estado de la orden empujado a Tiendanube")
}

func esTiendanube(platform string) bool {
	return strings.Contains(strings.ToLower(platform), "tiendanube")
}
