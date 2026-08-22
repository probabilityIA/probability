package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (uc *vtexUseCase) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	_, cred, err := uc.resolveOrdersAccess(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	limit := filters.Limit
	if limit <= 0 {
		limit = 100
	}

	rawFilters := map[string]interface{}{}
	if filters.From != nil {
		rawFilters["created_at_min"] = filters.From.Format(time.RFC3339)
	}
	if filters.To != nil {
		rawFilters["created_at_max"] = filters.To.Format(time.RFC3339)
	}
	vtexFilters := buildVTEXFilters(rawFilters)

	rows := make([]orderscompare.ChannelOrder, 0, limit)
	page := 1
	perPage := 15

	for {
		result, err := uc.client.GetOrders(ctx, cred, page, perPage, vtexFilters)
		if err != nil {
			return nil, fmt.Errorf("obteniendo ordenes de VTEX: %w", err)
		}
		if len(result.List) == 0 {
			break
		}

		for _, summary := range result.List {
			rows = append(rows, orderscompare.ChannelOrder{
				ExternalID: summary.OrderID,
				Number:     summary.Sequence,
				Status:     mapper.MapVTEXStatus(summary.Status),
				RawStatus:  summary.Status,
				Total:      float64(summary.TotalValue) / 100,
				Currency:   summary.CurrencyCode,
				CreatedAt:  summary.CreationDate,
			})
			if len(rows) >= limit {
				return rows, nil
			}
		}

		if page >= result.Paging.Pages {
			break
		}
		page++
	}

	return rows, nil
}

func (uc *vtexUseCase) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	result := orderscompare.ImportResult{
		Queued:   []string{},
		Failed:   map[string]string{},
		NotFound: []string{},
	}

	integration, cred, err := uc.resolveOrdersAccess(ctx, integrationID)
	if err != nil {
		return result, err
	}

	for _, externalID := range externalIDs {
		order, rawJSON, err := uc.client.GetOrderByID(ctx, cred, externalID)
		if err != nil {
			uc.logger.Warn(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", externalID).
				Msg("No se pudo leer la orden de VTEX para importarla")
			result.NotFound = append(result.NotFound, externalID)
			continue
		}

		dto := mapper.MapVTEXOrderToProbability(order, rawJSON)
		dto.IntegrationID = integration.ID
		dto.BusinessID = integration.BusinessID
		skip, _ := orderscompare.SkipsInventory(dto.Status)
		dto.SkipInventory = skip

		if err := uc.publisher.Publish(ctx, dto); err != nil {
			result.Failed[externalID] = err.Error()
			continue
		}
		result.Queued = append(result.Queued, externalID)
	}

	return result, nil
}

func (uc *vtexUseCase) resolveOrdersAccess(ctx context.Context, integrationID string) (*domain.Integration, domain.Credential, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, domain.Credential{}, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.Credential{}, domain.ErrIntegrationNotFound
	}

	cred, err := uc.resolveCredential(ctx, integration, integrationID)
	if err != nil {
		return nil, domain.Credential{}, err
	}
	return integration, cred, nil
}
