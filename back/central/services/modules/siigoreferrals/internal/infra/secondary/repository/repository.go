package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func referralToModel(r *entities.SiigoReferral) *models.SiigoReferral {
	return &models.SiigoReferral{
		Name:       r.Name,
		Email:      r.Email,
		Phone:      r.Phone,
		OrderRange: r.OrderRange,
	}
}

func (r *Repository) CreateReferral(ctx context.Context, referral *entities.SiigoReferral) error {
	model := referralToModel(referral)
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return err
	}
	referral.ID = model.ID
	referral.CreatedAt = model.CreatedAt
	return nil
}

func modelToReferral(m *models.SiigoReferral) entities.SiigoReferral {
	return entities.SiigoReferral{
		ID:         m.ID,
		Name:       m.Name,
		Email:      m.Email,
		Phone:      m.Phone,
		OrderRange: m.OrderRange,
		CreatedAt:  m.CreatedAt,
	}
}

func (r *Repository) ListReferrals(ctx context.Context, page, pageSize int) ([]entities.SiigoReferral, int64, error) {
	var models_ []models.SiigoReferral
	var total int64

	if err := r.db.Conn(ctx).Model(&models.SiigoReferral{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Conn(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models_).Error; err != nil {
		return nil, 0, err
	}

	referrals := make([]entities.SiigoReferral, 0, len(models_))
	for i := range models_ {
		referrals = append(referrals, modelToReferral(&models_[i]))
	}
	return referrals, total, nil
}
