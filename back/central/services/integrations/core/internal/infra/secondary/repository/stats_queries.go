package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) GetIntegrationStats(ctx context.Context, businessID uint) ([]domain.IntegrationStats, error) {
	var rows []models.IntegrationStat
	if err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]domain.IntegrationStats, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.IntegrationStats{
			IntegrationID:    row.IntegrationID,
			OrdersCount:      row.OrdersTotal,
			OrdersInProgress: row.OrdersInProgress,
			OrdersDelivered:  row.OrdersDelivered,
			OrdersCancelled:  row.OrdersCancelled,
			OrdersReturned:   row.OrdersReturned,
			ProductsCount:    row.ProductsCount,
			LastOrderAt:      row.LastOrderAt,
		})
	}
	return result, nil
}
