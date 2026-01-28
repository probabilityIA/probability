package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// IOrderRabbitPublisher define la interfaz para publicar eventos de órdenes a RabbitMQ
type IOrderRabbitPublisher interface {
	PublishConfirmationRequested(ctx context.Context, order *domain.ProbabilityOrder) error
}

// OrderRabbitPublisher implementa el publicador de eventos de órdenes
type OrderRabbitPublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// NewOrderRabbitPublisher crea un nuevo publisher de RabbitMQ
func NewOrderRabbitPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) IOrderRabbitPublisher {
	return &OrderRabbitPublisher{
		rabbit: rabbit,
		log:    logger,
	}
}

// PublishConfirmationRequested publica un evento cuando una orden requiere confirmación
func (p *OrderRabbitPublisher) PublishConfirmationRequested(ctx context.Context, order *domain.ProbabilityOrder) error {
	// Construir resumen de items
	itemsSummary := buildItemsSummary(order)

	// Construir dirección de envío
	shippingAddress := buildShippingAddress(order)

	// Construir evento
	event := map[string]interface{}{
		"event_type":       "order.confirmation_requested",
		"order_id":         order.ID,
		"order_number":     order.OrderNumber,
		"business_id":      order.BusinessID,
		"customer_name":    order.CustomerName,
		"customer_phone":   order.CustomerPhone,
		"customer_email":   order.CustomerEmail,
		"total_amount":     order.TotalAmount,
		"currency":         order.Currency,
		"items_summary":    itemsSummary,
		"shipping_address": shippingAddress,
		"payment_method_id": order.PaymentMethodID,
		"integration_id":   order.IntegrationID,
		"platform":         order.Platform,
		"timestamp":        time.Now().Unix(),
	}

	// Serializar a JSON
	payload, err := json.Marshal(event)
	if err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Msg("Error marshaling confirmation event")
		return fmt.Errorf("error marshaling event: %w", err)
	}

	// Publicar a RabbitMQ
	if err := p.rabbit.Publish(ctx, "orders.confirmation.requested", payload); err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("queue", "orders.confirmation.requested").
			Msg("Error publishing confirmation event to RabbitMQ")
		return fmt.Errorf("error publishing to RabbitMQ: %w", err)
	}

	p.log.Info().
		Str("order_id", order.ID).
		Str("order_number", order.OrderNumber).
		Str("queue", "orders.confirmation.requested").
		Msg("Confirmation event published successfully")

	return nil
}

// buildItemsSummary construye un resumen de los items de la orden
func buildItemsSummary(order *domain.ProbabilityOrder) string {
	if len(order.OrderItems) == 0 {
		return "Sin items"
	}

	summary := ""
	for i, item := range order.OrderItems {
		if i > 0 {
			summary += ", "
		}
		summary += fmt.Sprintf("%dx %s", item.Quantity, item.ProductName)
	}

	return summary
}

// buildShippingAddress construye un resumen de la dirección de envío
func buildShippingAddress(order *domain.ProbabilityOrder) string {
	if order.ShippingStreet == "" && order.ShippingCity == "" {
		return "Sin dirección"
	}

	address := order.ShippingStreet
	if order.ShippingCity != "" {
		if address != "" {
			address += ", "
		}
		address += order.ShippingCity
	}
	if order.ShippingState != "" {
		if address != "" {
			address += ", "
		}
		address += order.ShippingState
	}

	return address
}
