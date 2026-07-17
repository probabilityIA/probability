package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	ordersPageSize     = 100
	maxOrderPages      = 200
	pagePauseDefault   = 500 * time.Millisecond
	syncDateOnlyLayout = "2006-01-02"
)

func (uc *jumpsellerUseCase) emitOrderSyncEvent(ctx context.Context, integration *domain.Integration, eventType string, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	var businessID uint
	if integration.BusinessID != nil {
		businessID = *integration.BusinessID
	}
	_ = rabbitmq.PublishEvent(ctx, uc.rabbit, rabbitmq.EventEnvelope{
		Type:          eventType,
		Category:      "integration",
		BusinessID:    businessID,
		IntegrationID: integration.ID,
		Data:          data,
	})
}

func (uc *jumpsellerUseCase) SyncOrders(ctx context.Context, integrationID string) error {
	after := time.Now().AddDate(0, 0, -30)
	params := map[string]interface{}{
		"created_at_min": after.Format(time.RFC3339),
		"created_at_max": time.Now().Format(time.RFC3339),
	}
	return uc.SyncOrdersWithParams(ctx, integrationID, params)
}

func (uc *jumpsellerUseCase) SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error {
	integration, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	queryParams, err := buildQueryParams(params)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Msg("Rango de fechas invalido para sincronizar ordenes de Jumpseller")
		return err
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("base_url", cred.BaseURL).
		Bool("is_testing", integration.IsTesting).
		Msg("Starting Jumpseller order sync")

	go uc.syncOrdersAsync(context.Background(), integration, cred, queryParams)

	return nil
}

func (uc *jumpsellerUseCase) syncOrdersAsync(ctx context.Context, integration *domain.Integration, cred domain.Credential, params *domain.GetOrdersParams) {
	start := time.Now()
	uc.emitOrderSyncEvent(ctx, integration, "integration.sync.started", map[string]interface{}{})
	uc.emitOrderSyncEvent(ctx, integration, "jumpseller.orders.sync.started", map[string]interface{}{})

	totalSynced := 0
	failed := 0
	if params.PerPage == 0 {
		params.PerPage = ordersPageSize
	}

	for page := 1; page <= maxOrderPages; page++ {
		params.Page = page

		result, rawOrders, err := uc.client.GetOrders(ctx, cred, params)
		if err != nil {
			uc.logger.Error(ctx).Err(err).
				Int("page", page).
				Msg("Error fetching Jumpseller orders page")
			uc.emitOrderSyncEvent(ctx, integration, "integration.sync.failed", map[string]interface{}{
				"error": err.Error(),
			})
			uc.emitOrderSyncEvent(ctx, integration, "jumpseller.orders.sync.failed", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if len(result.Orders) == 0 {
			break
		}

		for i := range result.Orders {
			order := result.Orders[i]

			var rawJSON []byte
			if i < len(rawOrders) {
				rawJSON = rawOrders[i]
			}

			dto := mapper.MapJumpsellerOrderToProbability(&order, rawJSON)
			dto.IntegrationID = integration.ID
			dto.BusinessID = integration.BusinessID

			if err := uc.publisher.Publish(ctx, dto); err != nil {
				uc.logger.Error(ctx).Err(err).
					Str("order_number", dto.OrderNumber).
					Msg("Error publishing Jumpseller order")
				failed++
				uc.emitOrderSyncEvent(ctx, integration, "jumpseller.orders.sync.item", map[string]interface{}{
					"order_number":    dto.OrderNumber,
					"external_id":     dto.ExternalID,
					"original_status": dto.OriginalStatus,
					"action":          "failed",
				})
				continue
			}
			totalSynced++

			uc.emitOrderSyncEvent(ctx, integration, "jumpseller.orders.sync.item", map[string]interface{}{
				"order_number":    dto.OrderNumber,
				"external_id":     dto.ExternalID,
				"customer_name":   dto.CustomerName,
				"total":           dto.TotalAmount,
				"currency":        dto.Currency,
				"status":          dto.Status,
				"original_status": dto.OriginalStatus,
				"items":           len(dto.OrderItems),
				"action":          "imported",
			})
		}

		uc.emitOrderSyncEvent(ctx, integration, "jumpseller.orders.sync.progress", map[string]interface{}{
			"page":      page,
			"processed": totalSynced + failed,
			"imported":  totalSynced,
			"failed":    failed,
		})

		uc.logger.Info(ctx).
			Int("page", page).
			Int("orders_in_page", len(result.Orders)).
			Msg("Jumpseller orders page synced")

		if len(result.Orders) < params.PerPage {
			break
		}

		time.Sleep(pagePauseDefault)
	}

	uc.logger.Info(ctx).
		Int("total_synced", totalSynced).
		Uint("integration_id", integration.ID).
		Msg("Jumpseller order sync completed")

	uc.emitOrderSyncEvent(ctx, integration, "integration.sync.completed", map[string]interface{}{
		"total_fetched": totalSynced,
		"duration":      time.Since(start).Round(time.Millisecond).String(),
	})

	uc.emitOrderSyncEvent(ctx, integration, "jumpseller.orders.sync.completed", map[string]interface{}{
		"total_fetched": totalSynced,
		"imported":      totalSynced,
		"failed":        failed,
		"duration":      time.Since(start).Round(time.Millisecond).String(),
	})
}

func parseSyncDate(field, value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse(syncDateOnlyLayout, value); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("%s: %q no es una fecha valida (se espera YYYY-MM-DD o RFC3339)", field, value)
}

func buildQueryParams(params interface{}) (*domain.GetOrdersParams, error) {
	qp := &domain.GetOrdersParams{
		Statuses: []string{domain.StatusPaid, domain.StatusCanceled},
		PerPage:  ordersPageSize,
	}

	m, ok := params.(map[string]interface{})
	if !ok {
		return qp, nil
	}

	if v, ok := m["created_at_min"].(string); ok && v != "" {
		t, err := parseSyncDate("created_at_min", v)
		if err != nil {
			return nil, err
		}
		qp.After = &t
	}
	if v, ok := m["created_at_max"].(string); ok && v != "" {
		t, err := parseSyncDate("created_at_max", v)
		if err != nil {
			return nil, err
		}
		qp.Before = &t
	}
	if qp.After != nil && qp.Before == nil {
		now := time.Now()
		qp.Before = &now
	}

	if v, ok := m["status"].(string); ok && v != "" {
		qp.Statuses = []string{v}
	}

	return qp, nil
}
