package usecases

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func (uc *vtexUseCase) SyncOrders(ctx context.Context, integrationID string) error {
	after := time.Now().AddDate(0, 0, -30)
	params := map[string]interface{}{
		"created_at_min": after.Format(time.RFC3339),
	}
	return uc.SyncOrdersWithParams(ctx, integrationID, params)
}

func (uc *vtexUseCase) SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	cred, err := uc.resolveCredential(ctx, integration, integrationID)
	if err != nil {
		return err
	}

	filters := buildVTEXFilters(params)

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("account", cred.AccountName).
		Msg("Starting VTEX order sync")

	go uc.syncOrdersAsync(context.Background(), integration, cred, filters)

	return nil
}

func (uc *vtexUseCase) syncOrdersAsync(ctx context.Context, integration *domain.Integration, cred domain.Credential, filters map[string]string) {
	totalSynced := 0
	page := 1
	perPage := 15

	for {
		result, err := uc.client.GetOrders(ctx, cred, page, perPage, filters)
		if err != nil {
			uc.logger.Error(ctx).Err(err).
				Int("page", page).
				Msg("Error fetching VTEX orders page")
			break
		}

		if len(result.List) == 0 {
			break
		}

		for _, summary := range result.List {
			order, rawJSON, err := uc.client.GetOrderByID(ctx, cred, summary.OrderID)
			if err != nil {
				uc.logger.Error(ctx).Err(err).
					Str("order_id", summary.OrderID).
					Msg("Error fetching VTEX order detail")
				continue
			}

			dto := mapper.MapVTEXOrderToProbability(order, rawJSON)
			dto.IntegrationID = integration.ID
			dto.BusinessID = integration.BusinessID

			if err := uc.publisher.Publish(ctx, dto); err != nil {
				uc.logger.Error(ctx).Err(err).
					Str("order_id", summary.OrderID).
					Msg("Error publishing VTEX order")
				continue
			}
			totalSynced++
		}

		uc.logger.Info(ctx).
			Int("page", page).
			Int("orders_in_page", len(result.List)).
			Int("total", result.Paging.Total).
			Msg("VTEX orders page synced")

		if page >= result.Paging.Pages {
			break
		}

		page++
		time.Sleep(500 * time.Millisecond)
	}

	uc.logger.Info(ctx).
		Int("total_synced", totalSynced).
		Uint("integration_id", integration.ID).
		Msg("VTEX order sync completed")
}

func buildVTEXFilters(params interface{}) map[string]string {
	filters := make(map[string]string)

	m, ok := params.(map[string]interface{})
	if !ok {
		return filters
	}

	var dateFrom, dateTo string
	if v, ok := m["created_at_min"].(string); ok && v != "" {
		dateFrom = v
	}
	if v, ok := m["created_at_max"].(string); ok && v != "" {
		dateTo = v
	} else {
		dateTo = time.Now().Format(time.RFC3339)
	}

	if dateFrom != "" {
		filters["f_creationDate"] = url.QueryEscape(fmt.Sprintf("creationDate:[%s TO %s]", dateFrom, dateTo))
	}

	if v, ok := m["status"].(string); ok && v != "" {
		filters["f_status"] = v
	}

	return filters
}
