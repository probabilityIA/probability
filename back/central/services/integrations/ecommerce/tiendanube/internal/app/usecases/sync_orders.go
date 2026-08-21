package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func (uc *tiendanubeUseCase) publishOrder(ctx context.Context, integration *domain.Integration, order *domain.TiendanubeOrder, rawJSON []byte) error {
	dto := mapper.MapTiendanubeOrderToProbability(order, rawJSON)
	dto.IntegrationID = integration.ID
	dto.BusinessID = integration.BusinessID

	if err := uc.publisher.Publish(ctx, dto); err != nil {
		return fmt.Errorf("publicando la orden %s de Tiendanube: %w", order.ID, err)
	}
	return nil
}

func (uc *tiendanubeUseCase) SyncOrders(ctx context.Context, integrationID string, filters domain.OrderFilters) (int, error) {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return 0, err
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return 0, err
	}

	orders, err := uc.client.GetOrders(ctx, cred, filters)
	if err != nil {
		return 0, fmt.Errorf("obteniendo ordenes de Tiendanube: %w", err)
	}

	published := 0
	for i := range orders {
		order := orders[i]
		if err := uc.publishOrder(ctx, integration, &order, nil); err != nil {
			uc.logger.Error(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", order.ID).
				Msg("No se pudo publicar una orden de Tiendanube")
			continue
		}
		published++
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Int("fetched", len(orders)).
		Int("published", published).
		Msg("Sincronizacion de ordenes de Tiendanube completada")

	return published, nil
}

func (uc *tiendanubeUseCase) ProcessOrderEvent(ctx context.Context, integrationID, event, orderID string) error {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return err
	}

	order, rawJSON, err := uc.client.GetOrder(ctx, cred, orderID)
	if err != nil {
		return fmt.Errorf("obteniendo la orden %s de Tiendanube: %w", orderID, err)
	}

	if err := uc.publishOrder(ctx, integration, order, rawJSON); err != nil {
		return err
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("event", event).
		Str("order_id", orderID).
		Str("order_number", order.Number).
		Msg("Orden de Tiendanube publicada desde webhook")

	return nil
}

func IsOrderEvent(event string) bool {
	return strings.HasPrefix(strings.TrimSpace(event), "order/")
}
