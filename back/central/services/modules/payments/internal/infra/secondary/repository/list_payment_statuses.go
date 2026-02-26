package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ListPaymentStatuses(ctx context.Context, isActive *bool) ([]entities.PaymentStatus, error) {
	var statuses []models.PaymentStatus

	query := r.db.Conn(ctx).Model(&models.PaymentStatus{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("code ASC").Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.PaymentStatus, len(statuses))
	for i, s := range statuses {
		result[i] = entities.PaymentStatus{
			ID:          s.ID,
			Code:        s.Code,
			Name:        s.Name,
			Description: s.Description,
			Category:    s.Category,
			IsActive:    s.IsActive,
			Icon:        s.Icon,
			Color:       s.Color,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		}
		if s.DeletedAt.Valid {
			result[i].DeletedAt = &s.DeletedAt.Time
		}
	}

	return result, nil
}
