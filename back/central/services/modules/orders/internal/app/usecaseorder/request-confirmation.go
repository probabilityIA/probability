package usecaseorder

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IRequestConfirmationUseCase define la interfaz del caso de uso
type IRequestConfirmationUseCase interface {
	RequestConfirmation(ctx context.Context, orderID string) error
}

// RequestConfirmationUseCase implementa el caso de uso de solicitud de confirmación
type RequestConfirmationUseCase struct {
	repository      ports.IRepository
	rabbitPublisher ports.IOrderRabbitPublisher
	log             log.ILogger
}

// NewRequestConfirmationUseCase crea una nueva instancia del caso de uso
func NewRequestConfirmationUseCase(
	repo ports.IRepository,
	rabbitPublisher ports.IOrderRabbitPublisher,
	logger log.ILogger,
) IRequestConfirmationUseCase {
	return &RequestConfirmationUseCase{
		repository:      repo,
		rabbitPublisher: rabbitPublisher,
		log:             logger,
	}
}

// RequestConfirmation solicita confirmación de una orden
func (uc *RequestConfirmationUseCase) RequestConfirmation(ctx context.Context, orderID string) error {
	// 1. Obtener la orden
	order, err := uc.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		uc.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("Error getting order for confirmation request")
		return fmt.Errorf("error getting order: %w", err)
	}

	// 2. Validar que tenga teléfono del cliente
	if order.CustomerPhone == "" {
		uc.log.Warn().
			Str("order_id", orderID).
			Str("order_number", order.OrderNumber).
			Msg("Cannot request confirmation: order has no customer phone")
		return fmt.Errorf("order does not have customer phone")
	}

	// 3. Validar que no esté ya confirmada
	if order.IsConfirmed != nil && *order.IsConfirmed {
		uc.log.Warn().
			Str("order_id", orderID).
			Str("order_number", order.OrderNumber).
			Msg("Cannot request confirmation: order already confirmed")
		return fmt.Errorf("order is already confirmed")
	}

	// 4. Publicar evento a RabbitMQ
	if err := uc.rabbitPublisher.PublishConfirmationRequested(ctx, order); err != nil {
		uc.log.Error().
			Err(err).
			Str("order_id", orderID).
			Str("order_number", order.OrderNumber).
			Msg("Error publishing confirmation request event")
		return fmt.Errorf("error publishing confirmation request: %w", err)
	}

	uc.log.Info().
		Str("order_id", orderID).
		Str("order_number", order.OrderNumber).
		Str("customer_phone", order.CustomerPhone).
		Msg("Confirmation request published successfully")

	return nil
}
