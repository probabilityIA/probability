package repository

import (
	"context"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/ports"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

const defaultLocalLimit = 1000

func (r *repository) GetIntegration(ctx context.Context, integrationID uint) (*ports.Integration, error) {
	var row models.Integration
	err := r.db.Conn(ctx).
		Select("id", "name", "business_id", "integration_type_id", "is_active").
		Where("id = ?", integrationID).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ports.Integration{
		ID:              row.ID,
		Name:            row.Name,
		BusinessID:      row.BusinessID,
		IntegrationType: row.IntegrationTypeID,
		IsActive:        row.IsActive,
	}, nil
}

func (r *repository) ListLocalOrders(ctx context.Context, integrationID, businessID uint, from, to *time.Time, limit int) ([]dtos.LocalOrder, error) {
	type localRow struct {
		ID           string
		OrderNumber  string
		ExternalID   string
		Status       string
		TotalAmount  float64
		Currency     string
		CustomerName string
		OccurredAt   time.Time
		CreatedAt    time.Time
	}

	if limit <= 0 || limit > defaultLocalLimit {
		limit = defaultLocalLimit
	}

	query := r.db.Conn(ctx).
		Model(&models.Order{}).
		Select("id", "order_number", "external_id", "status", "total_amount", "currency", "customer_name", "occurred_at", "created_at").
		Where("integration_id = ? AND business_id = ? AND deleted_at IS NULL", integrationID, businessID)

	if from != nil {
		query = query.Where("COALESCE(occurred_at, created_at) >= ?", *from)
	}
	if to != nil {
		query = query.Where("COALESCE(occurred_at, created_at) <= ?", *to)
	}

	var rows []localRow
	if err := query.Order("COALESCE(occurred_at, created_at) DESC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	orders := make([]dtos.LocalOrder, 0, len(rows))
	for _, row := range rows {
		createdAt := row.OccurredAt
		if createdAt.IsZero() {
			createdAt = row.CreatedAt
		}
		orders = append(orders, dtos.LocalOrder{
			OrderID:      row.ID,
			OrderNumber:  row.OrderNumber,
			ExternalID:   row.ExternalID,
			Status:       row.Status,
			Total:        row.TotalAmount,
			Currency:     row.Currency,
			CustomerName: row.CustomerName,
			CreatedAt:    createdAt,
		})
	}
	return orders, nil
}
