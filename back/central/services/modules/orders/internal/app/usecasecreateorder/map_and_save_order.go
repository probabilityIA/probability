package usecasecreateorder

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
)

func (uc *UseCaseCreateOrder) MapAndSaveOrder(ctx context.Context, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error) {
	if dto.IntegrationID == 0 && !dto.IsManualOrder {
		return nil, errors.New("integration_id is required")
	}
	if dto.BusinessID == nil || *dto.BusinessID == 0 {
		return nil, errors.New("business_id is required")
	}

	exists, err := uc.repo.OrderExists(ctx, dto.ExternalID, dto.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error checking if order exists: %w", err)
	}
	if exists {
		existingOrder, err := uc.repo.GetOrderByExternalID(ctx, dto.ExternalID, dto.IntegrationID)
		if err != nil {
			return nil, fmt.Errorf("error getting existing order: %w", err)
		}
		return uc.updateUseCase.UpdateOrder(ctx, existingOrder, dto)
	}

	client, err := uc.GetOrCreateCustomer(ctx, *dto.BusinessID, dto)
	if err != nil {
		return nil, fmt.Errorf("error processing customer: %w", err)
	}
	var clientID *uint
	if client != nil {
		clientID = &client.ID
	}

	if client != nil && dto.ClientGroupID != nil && *dto.ClientGroupID > 0 {
		if err := uc.repo.AssignClientToGroup(ctx, *dto.BusinessID, *dto.ClientGroupID, client.ID); err != nil {
			uc.logger.Warn(ctx).Err(err).
				Uint("client_id", client.ID).
				Uint("client_group_id", *dto.ClientGroupID).
				Msg("No se pudo asignar el cliente al grupo de precios")
		}
	}

	if dto.CustomerName == "" && (dto.CustomerFirstName != "" || dto.CustomerLastName != "") {
		dto.CustomerName = fmt.Sprintf("%s %s", dto.CustomerFirstName, dto.CustomerLastName)
	}

	statusMapping := uc.mapOrderStatuses(ctx, dto)

	order := uc.buildOrderEntity(dto, clientID, statusMapping)

	uc.hydrateBusinessName(ctx, order)

	uc.assignPaymentMethodID(order, dto)

	uc.syncIsPaidFromPaymentStatus(ctx, order, statusMapping.PaymentStatusID)

	uc.populateOrderFields(order, dto)

	uc.geocodeOrderIfNeeded(ctx, order)

	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		if errors.Is(err, domainerrors.ErrOrderAlreadyExists) {
			existingOrder, gerr := uc.repo.GetOrderByExternalID(ctx, dto.ExternalID, dto.IntegrationID)
			if gerr != nil {
				return nil, fmt.Errorf("error getting existing order after create conflict: %w", gerr)
			}
			return uc.updateUseCase.UpdateOrder(ctx, existingOrder, dto)
		}
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	if dto.BusinessID != nil {
		if err := uc.repo.ResolveOrderGeozone(ctx, order.ID, *dto.BusinessID); err != nil {
			uc.logger.Warn(ctx).Err(err).Str("order_id", order.ID).Msg("Failed to resolve order geozone")
		}
	}

	if err := uc.saveRelatedEntities(ctx, order, dto); err != nil {
		return nil, err
	}

	uc.publishOrderEvents(ctx, order, dto.IsManualOrder)

	return uc.mapOrderToResponse(order), nil
}

func (uc *UseCaseCreateOrder) hydrateBusinessName(ctx context.Context, order *entities.ProbabilityOrder) {
	if order == nil || order.BusinessName != "" || order.BusinessID == nil || *order.BusinessID == 0 {
		return
	}
	name, err := uc.repo.GetBusinessNameByID(ctx, *order.BusinessID)
	if err != nil {
		uc.logger.Warn(ctx).Err(err).Uint("business_id", *order.BusinessID).Str("order_id", order.ID).Msg("No se pudo resolver el nombre del negocio para la notificacion")
		return
	}
	order.BusinessName = name
}
