package repository

import (
	"context"
	stderrors "errors"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}

func (r *Repository) Create(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
	model := &models.ShippingMargin{
		BusinessID:      m.BusinessID,
		CarrierCode:     m.CarrierCode,
		CarrierName:     m.CarrierName,
		MarginAmount:    m.MarginAmount,
		InsuranceMargin: m.InsuranceMargin,
		IsActive:        m.IsActive,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	m.ID = model.ID
	m.CreatedAt = model.CreatedAt
	m.UpdatedAt = model.UpdatedAt
	return m, nil
}

func (r *Repository) GetByID(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
	var model models.ShippingMargin
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", id, businessID).
		First(&model).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrShippingMarginNotFound
		}
		return nil, err
	}
	return modelToEntity(&model), nil
}

func (r *Repository) GetByBusinessAndCarrier(ctx context.Context, businessID uint, carrierCode string) (*entities.ShippingMargin, error) {
	var model models.ShippingMargin
	err := r.db.Conn(ctx).
		Where("business_id = ? AND carrier_code = ?", businessID, carrierCode).
		First(&model).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrShippingMarginNotFound
		}
		return nil, err
	}
	return modelToEntity(&model), nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListShippingMarginsParams) ([]entities.ShippingMargin, int64, error) {
	var modelsList []models.ShippingMargin
	var total int64

	query := r.db.Conn(ctx).Model(&models.ShippingMargin{}).
		Where("business_id = ?", params.BusinessID)

	if params.CarrierCode != "" {
		query = query.Where("carrier_code = ?", params.CarrierCode)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).
		Order("carrier_code ASC").
		Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	out := make([]entities.ShippingMargin, len(modelsList))
	for i, m := range modelsList {
		out[i] = *modelToEntity(&m)
	}
	return out, total, nil
}

func (r *Repository) Update(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
	updates := map[string]interface{}{
		"carrier_name":     m.CarrierName,
		"margin_amount":    m.MarginAmount,
		"insurance_margin": m.InsuranceMargin,
		"is_active":        m.IsActive,
	}
	res := r.db.Conn(ctx).Model(&models.ShippingMargin{}).
		Where("id = ? AND business_id = ?", m.ID, m.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrShippingMarginNotFound
	}
	return r.GetByID(ctx, m.BusinessID, m.ID)
}

func (r *Repository) Delete(ctx context.Context, businessID, id uint) error {
	res := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", id, businessID).
		Delete(&models.ShippingMargin{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrShippingMarginNotFound
	}
	return nil
}

func (r *Repository) ExistsByCarrier(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
	var count int64
	q := r.db.Conn(ctx).Model(&models.ShippingMargin{}).
		Where("business_id = ? AND carrier_code = ?", businessID, carrierCode)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func modelToEntity(m *models.ShippingMargin) *entities.ShippingMargin {
	return &entities.ShippingMargin{
		ID:              m.ID,
		BusinessID:      m.BusinessID,
		CarrierCode:     m.CarrierCode,
		CarrierName:     m.CarrierName,
		MarginAmount:    m.MarginAmount,
		InsuranceMargin: m.InsuranceMargin,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
