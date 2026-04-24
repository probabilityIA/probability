package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreatePutawayRule(ctx context.Context, rule *entities.PutawayRule) (*entities.PutawayRule, error) {
	m := mappers.PutawayRuleEntityToModel(rule)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.PutawayRuleModelToEntity(m), nil
}

func (r *Repository) ListPutawayRules(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error) {
	var ml []models.PutawayRule
	var total int64

	q := r.db.Conn(ctx).Model(&models.PutawayRule{}).Where("business_id = ?", params.BusinessID)
	if params.ActiveOnly {
		q = q.Where("is_active = true")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("priority DESC, id ASC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}

	rules := make([]entities.PutawayRule, len(ml))
	for i := range ml {
		rules[i] = *mappers.PutawayRuleModelToEntity(&ml[i])
	}
	return rules, total, nil
}

func (r *Repository) GetPutawayRuleByID(ctx context.Context, businessID, ruleID uint) (*entities.PutawayRule, error) {
	var m models.PutawayRule
	if err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", ruleID, businessID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrPutawayRuleNotFound
		}
		return nil, err
	}
	return mappers.PutawayRuleModelToEntity(&m), nil
}

func (r *Repository) UpdatePutawayRule(ctx context.Context, rule *entities.PutawayRule) (*entities.PutawayRule, error) {
	updates := map[string]any{
		"product_id":     rule.ProductID,
		"category_id":    rule.CategoryID,
		"target_zone_id": rule.TargetZoneID,
		"priority":       rule.Priority,
		"strategy":       rule.Strategy,
		"is_active":      rule.IsActive,
	}
	res := r.db.Conn(ctx).Model(&models.PutawayRule{}).
		Where("id = ? AND business_id = ?", rule.ID, rule.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrPutawayRuleNotFound
	}
	return r.GetPutawayRuleByID(ctx, rule.BusinessID, rule.ID)
}

func (r *Repository) DeletePutawayRule(ctx context.Context, businessID, ruleID uint) error {
	res := r.db.Conn(ctx).Where("id = ? AND business_id = ?", ruleID, businessID).Delete(&models.PutawayRule{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrPutawayRuleNotFound
	}
	return nil
}

func (r *Repository) FindApplicableRule(ctx context.Context, businessID uint, productID string) (*entities.PutawayRule, error) {
	var m models.PutawayRule
	err := r.db.Conn(ctx).
		Where("business_id = ? AND is_active = true", businessID).
		Where("product_id = ? OR product_id IS NULL", productID).
		Order("product_id IS NOT NULL DESC, priority DESC").
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrNoPutawayRule
		}
		return nil, err
	}
	return mappers.PutawayRuleModelToEntity(&m), nil
}

func (r *Repository) PickLocationInZone(ctx context.Context, zoneID uint) (uint, error) {
	type row struct {
		ID uint
	}
	var result row
	err := r.db.Conn(ctx).
		Table("warehouse_locations wl").
		Select("wl.id").
		Joins("INNER JOIN warehouse_rack_levels wrl ON wrl.id = wl.level_id AND wrl.deleted_at IS NULL").
		Joins("INNER JOIN warehouse_racks wr ON wr.id = wrl.rack_id AND wr.deleted_at IS NULL").
		Joins("INNER JOIN warehouse_aisles wa ON wa.id = wr.aisle_id AND wa.deleted_at IS NULL").
		Where("wa.zone_id = ? AND wl.deleted_at IS NULL AND wl.is_active = true", zoneID).
		Order("wl.priority DESC, wl.id ASC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	if result.ID == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return result.ID, nil
}

func (r *Repository) CreatePutawaySuggestion(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
	m := mappers.PutawaySuggestionEntityToModel(s)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.PutawaySuggestionModelToEntity(m), nil
}

func (r *Repository) GetPutawaySuggestionByID(ctx context.Context, businessID, id uint) (*entities.PutawaySuggestion, error) {
	var m models.PutawaySuggestion
	if err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrPutawaySuggestionNotFound
		}
		return nil, err
	}
	return mappers.PutawaySuggestionModelToEntity(&m), nil
}

func (r *Repository) ListPutawaySuggestions(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error) {
	var ml []models.PutawaySuggestion
	var total int64
	q := r.db.Conn(ctx).Model(&models.PutawaySuggestion{}).Where("business_id = ?", params.BusinessID)
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("id DESC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entities.PutawaySuggestion, len(ml))
	for i := range ml {
		out[i] = *mappers.PutawaySuggestionModelToEntity(&ml[i])
	}
	return out, total, nil
}

func (r *Repository) UpdatePutawaySuggestion(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
	updates := map[string]any{
		"status":             s.Status,
		"actual_location_id": s.ActualLocationID,
		"confirmed_at":       s.ConfirmedAt,
		"confirmed_by_id":    s.ConfirmedByID,
	}
	res := r.db.Conn(ctx).Model(&models.PutawaySuggestion{}).
		Where("id = ? AND business_id = ?", s.ID, s.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrPutawaySuggestionNotFound
	}
	return r.GetPutawaySuggestionByID(ctx, s.BusinessID, s.ID)
}
