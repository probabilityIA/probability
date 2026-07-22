package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

type IntegrationStatsItem struct {
	IntegrationID    uint       `json:"integration_id"`
	OrdersCount      int64      `json:"orders_count"`
	OrdersInProgress int64      `json:"orders_in_progress"`
	OrdersDelivered  int64      `json:"orders_delivered"`
	OrdersCancelled  int64      `json:"orders_cancelled"`
	OrdersReturned   int64      `json:"orders_returned"`
	ProductsCount    int64      `json:"products_count"`
	LastOrderAt      *time.Time `json:"last_order_at,omitempty"`
}

type IntegrationStatsResponse struct {
	Success bool                   `json:"success"`
	Data    []IntegrationStatsItem `json:"data"`
}

func ToIntegrationStatsResponse(stats []domain.IntegrationStats) IntegrationStatsResponse {
	items := make([]IntegrationStatsItem, 0, len(stats))
	for _, s := range stats {
		items = append(items, IntegrationStatsItem{
			IntegrationID:    s.IntegrationID,
			OrdersCount:      s.OrdersCount,
			OrdersInProgress: s.OrdersInProgress,
			OrdersDelivered:  s.OrdersDelivered,
			OrdersCancelled:  s.OrdersCancelled,
			OrdersReturned:   s.OrdersReturned,
			ProductsCount:    s.ProductsCount,
			LastOrderAt:      s.LastOrderAt,
		})
	}
	return IntegrationStatsResponse{Success: true, Data: items}
}
