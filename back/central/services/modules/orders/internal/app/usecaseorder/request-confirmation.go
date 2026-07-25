package usecaseorder

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RequestConfirmationUseCase struct {
	repository      ports.IRepository
	rabbitPublisher ports.IOrderRabbitPublisher
	log             log.ILogger
}

func NewRequestConfirmationUseCase(
	repo ports.IRepository,
	rabbitPublisher ports.IOrderRabbitPublisher,
	logger log.ILogger,
) ports.IRequestConfirmationUseCase {
	return &RequestConfirmationUseCase{
		repository:      repo,
		rabbitPublisher: rabbitPublisher,
		log:             logger,
	}
}

func (uc *RequestConfirmationUseCase) RequestConfirmation(ctx context.Context, orderID string, businessID uint) error {
	order, err := uc.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		uc.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("Error getting order for confirmation request")
		return fmt.Errorf("error getting order: %w", err)
	}

	if order == nil {
		return fmt.Errorf("order not found")
	}

	if order.BusinessID == nil || *order.BusinessID != businessID {
		uc.log.Warn().
			Str("order_id", orderID).
			Uint("business_id", businessID).
			Msg("Cannot request confirmation: order belongs to another business")
		return fmt.Errorf("la orden no pertenece al negocio")
	}

	if order.CustomerPhone == "" {
		uc.log.Warn().
			Str("order_id", orderID).
			Str("order_number", order.OrderNumber).
			Msg("Cannot request confirmation: order has no customer phone")
		return fmt.Errorf("order does not have customer phone")
	}

	if order.IsConfirmed != nil && *order.IsConfirmed {
		uc.log.Warn().
			Str("order_id", orderID).
			Str("order_number", order.OrderNumber).
			Msg("Cannot request confirmation: order already confirmed")
		return fmt.Errorf("order is already confirmed")
	}

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
