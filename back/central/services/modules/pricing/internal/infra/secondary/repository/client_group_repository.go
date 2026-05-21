package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateClientGroup(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error) {
	model := &models.ClientGroup{
		BusinessID:  group.BusinessID,
		Name:        group.Name,
		Description: group.Description,
		Color:       group.Color,
		IsActive:    group.IsActive,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	group.ID = model.ID
	group.CreatedAt = model.CreatedAt
	group.UpdatedAt = model.UpdatedAt
	return group, nil
}

func (r *Repository) UpdateClientGroup(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error) {
	result := r.db.Conn(ctx).Model(&models.ClientGroup{}).
		Where("id = ? AND business_id = ?", group.ID, group.BusinessID).
		Updates(map[string]any{
			"name":        group.Name,
			"description": group.Description,
			"color":       group.Color,
			"is_active":   group.IsActive,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domainerrors.ErrGroupNotFound
	}
	return r.GetClientGroup(ctx, group.BusinessID, group.ID)
}

func (r *Repository) GetClientGroup(ctx context.Context, businessID, groupID uint) (*entities.ClientGroup, error) {
	var model models.ClientGroup
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", groupID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrGroupNotFound
		}
		return nil, err
	}

	var memberCount int64
	if err := r.db.Conn(ctx).Model(&models.ClientGroupMember{}).
		Where("client_group_id = ? AND business_id = ?", groupID, businessID).
		Count(&memberCount).Error; err != nil {
		return nil, err
	}

	return &entities.ClientGroup{
		ID:          model.ID,
		BusinessID:  model.BusinessID,
		Name:        model.Name,
		Description: model.Description,
		Color:       model.Color,
		IsActive:    model.IsActive,
		MemberCount: memberCount,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

type groupWithCount struct {
	models.ClientGroup
	MemberCount int64 `gorm:"column:member_count"`
}

func (r *Repository) ListClientGroups(ctx context.Context, params dtos.ListClientGroupsParams) ([]entities.ClientGroup, int64, error) {
	var total int64
	countQuery := r.db.Conn(ctx).Model(&models.ClientGroup{}).
		Where("business_id = ?", params.BusinessID)
	if params.Search != "" {
		countQuery = countQuery.Where("name ILIKE ?", "%"+params.Search+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []groupWithCount
	query := r.db.Conn(ctx).
		Table("client_group cg").
		Select("cg.*, (SELECT COUNT(*) FROM client_group_member cgm WHERE cgm.client_group_id = cg.id AND cgm.deleted_at IS NULL) AS member_count").
		Where("cg.business_id = ? AND cg.deleted_at IS NULL", params.BusinessID)
	if params.Search != "" {
		query = query.Where("cg.name ILIKE ?", "%"+params.Search+"%")
	}
	if err := query.Order("cg.name ASC").
		Offset(params.Offset()).Limit(params.PageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	groups := make([]entities.ClientGroup, len(rows))
	for i, row := range rows {
		groups[i] = entities.ClientGroup{
			ID:          row.ID,
			BusinessID:  row.BusinessID,
			Name:        row.Name,
			Description: row.Description,
			Color:       row.Color,
			IsActive:    row.IsActive,
			MemberCount: row.MemberCount,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		}
	}
	return groups, total, nil
}

func (r *Repository) DeleteClientGroup(ctx context.Context, businessID, groupID uint) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND business_id = ?", groupID, businessID).
			Delete(&models.ClientGroup{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domainerrors.ErrGroupNotFound
		}
		if err := tx.Unscoped().
			Where("client_group_id = ? AND business_id = ?", groupID, businessID).
			Delete(&models.ClientGroupMember{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().
			Where("client_group_id = ? AND business_id = ?", groupID, businessID).
			Delete(&models.CustomProductPrice{}).Error
	})
}

func (r *Repository) GroupNameExists(ctx context.Context, businessID, groupID uint, name string) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.ClientGroup{}).
		Where("business_id = ? AND LOWER(name) = LOWER(?)", businessID, name)
	if groupID > 0 {
		query = query.Where("id != ?", groupID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
