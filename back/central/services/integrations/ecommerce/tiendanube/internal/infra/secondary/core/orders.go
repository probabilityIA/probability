package core

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

const defaultSyncWindowDays = 30

func (t *TiendanubeCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	filters := domain.OrderFilters{
		CreatedAtMin: time.Now().AddDate(0, 0, -defaultSyncWindowDays).Format(time.RFC3339),
		CreatedAtMax: time.Now().Format(time.RFC3339),
	}
	_, err := t.useCase.SyncOrders(ctx, integrationID, filters)
	return err
}

func (t *TiendanubeCore) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	_, err := t.useCase.SyncOrders(ctx, integrationID, buildOrderFilters(params))
	return err
}

func buildOrderFilters(params interface{}) domain.OrderFilters {
	filters := domain.OrderFilters{
		CreatedAtMin: time.Now().AddDate(0, 0, -defaultSyncWindowDays).Format(time.RFC3339),
		CreatedAtMax: time.Now().Format(time.RFC3339),
	}

	m, ok := params.(map[string]interface{})
	if !ok {
		return filters
	}

	if v, ok := m["created_at_min"].(string); ok && v != "" {
		filters.CreatedAtMin = v
	}
	if v, ok := m["created_at_max"].(string); ok && v != "" {
		filters.CreatedAtMax = v
	}
	if v, ok := m["status"].(string); ok && v != "" {
		filters.Status = v
	}
	if v, ok := m["payment_status"].(string); ok && v != "" {
		filters.PaymentStatus = v
	}
	if v, ok := m["limit"].(float64); ok && v > 0 {
		filters.Limit = int(v)
	}
	if v, ok := m["limit"].(int); ok && v > 0 {
		filters.Limit = v
	}

	return filters
}
