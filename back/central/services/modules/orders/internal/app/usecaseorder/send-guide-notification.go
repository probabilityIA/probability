package usecaseorder

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// SendGuideNotificationUseCase implementa el caso de uso de envio de notificacion de guia
type SendGuideNotificationUseCase struct {
	repository      ports.IRepository
	rabbitPublisher ports.IOrderRabbitPublisher
	log             log.ILogger
}

// NewSendGuideNotificationUseCase crea una nueva instancia del caso de uso
func NewSendGuideNotificationUseCase(
	repo ports.IRepository,
	rabbitPublisher ports.IOrderRabbitPublisher,
	logger log.ILogger,
) ports.ISendGuideNotificationUseCase {
	return &SendGuideNotificationUseCase{
		repository:      repo,
		rabbitPublisher: rabbitPublisher,
		log:             logger,
	}
}

// SendGuideNotification envia la notificacion de guia de envio por WhatsApp
func (uc *SendGuideNotificationUseCase) SendGuideNotification(ctx context.Context, orderID string) error {
	// 1. Obtener la orden
	order, err := uc.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		uc.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("Error getting order for guide notification")
		return fmt.Errorf("error getting order: %w", err)
	}

	// 2. Validar que tenga telefono del cliente
	if order.CustomerPhone == "" {
		uc.log.Warn().
			Str("order_id", orderID).
			Str("order_number", order.OrderNumber).
			Msg("Cannot send guide notification: order has no customer phone")
		return fmt.Errorf("order does not have customer phone")
	}

	// 3. Validar que tenga tracking number (guia generada)
	if order.TrackingNumber == nil || *order.TrackingNumber == "" {
		uc.log.Warn().
			Str("order_id", orderID).
			Str("order_number", order.OrderNumber).
			Msg("Cannot send guide notification: order has no tracking number")
		return fmt.Errorf("order does not have tracking number")
	}

	// 4. Publicar evento a RabbitMQ
	if err := uc.rabbitPublisher.PublishGuideNotificationRequested(ctx, order); err != nil {
		uc.log.Error().
			Err(err).
			Str("order_id", orderID).
			Str("order_number", order.OrderNumber).
			Msg("Error publishing guide notification event")
		return fmt.Errorf("error publishing guide notification: %w", err)
	}

	uc.log.Info().
		Str("order_id", orderID).
		Str("order_number", order.OrderNumber).
		Str("customer_phone", order.CustomerPhone).
		Str("tracking_number", *order.TrackingNumber).
		Msg("Guide notification published successfully")

	return nil
}
