package usecaseorder

import (
	"context"
	"fmt"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder/mapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseOrder) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*dtos.OrdersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	orders, total, err := uc.repo.ListOrders(ctx, page, pageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing orders: %w", err)
	}

	orderSummaries := make([]dtos.OrderSummary, len(orders))
	for i, order := range orders {
		orderSummaries[i] = mapper.ToOrderSummary(&order)
	}

	uc.attachNotificationCounters(ctx, filters, orderSummaries)

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dtos.OrdersListResponse{
		Data:       orderSummaries,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (uc *UseCaseOrder) attachNotificationCounters(ctx context.Context, filters map[string]interface{}, summaries []dtos.OrderSummary) {
	if len(summaries) == 0 {
		return
	}

	businessID, ok := filters["business_id"].(uint)
	if !ok || businessID == 0 {
		return
	}

	keys := make([]entities.OrderNotificationKey, 0, len(summaries))
	for _, s := range summaries {
		keys = append(keys, entities.OrderNotificationKey{
			OrderID:       s.ID,
			OrderNumber:   s.OrderNumber,
			IntegrationID: s.IntegrationID,
		})
	}

	counters, err := uc.repo.GetOrderNotificationCounters(ctx, businessID, keys)
	if err != nil {
		uc.logger.Error().Err(err).
			Uint("business_id", businessID).
			Msg("Error getting order notification counters")
		return
	}

	for i := range summaries {
		found, ok := counters[summaries[i].ID]
		if !ok {
			continue
		}
		summaries[i].Notifications = toNotificationCounterDTOs(found)
	}
}

func toNotificationCounterDTOs(counters []entities.OrderNotificationCounter) []dtos.NotificationCounter {
	out := make([]dtos.NotificationCounter, 0, len(counters))
	for _, c := range counters {
		out = append(out, dtos.NotificationCounter{
			Channel:  c.Channel,
			Sent:     c.Sent,
			Expected: c.Expected,
			Failed:   c.Failed,
			State:    c.State,
		})
	}
	return out
}
