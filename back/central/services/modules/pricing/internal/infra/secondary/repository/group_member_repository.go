package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type clientRow struct {
	ID        uint    `gorm:"column:id"`
	Name      string  `gorm:"column:name"`
	Email     *string `gorm:"column:email"`
	Phone     string  `gorm:"column:phone"`
	Dni       *string `gorm:"column:dni"`
	GroupID   *uint   `gorm:"column:group_id"`
	GroupName *string `gorm:"column:group_name"`
}

func clientRowToSummary(row clientRow) entities.ClientSummary {
	summary := entities.ClientSummary{
		ID:      row.ID,
		Name:    row.Name,
		Phone:   row.Phone,
		GroupID: row.GroupID,
	}
	if row.Email != nil {
		summary.Email = *row.Email
	}
	if row.Dni != nil {
		summary.Dni = *row.Dni
	}
	if row.GroupName != nil {
		summary.GroupName = *row.GroupName
	}
	return summary
}

func (r *Repository) ListGroupMembers(ctx context.Context, params dtos.ListGroupMembersParams) ([]entities.ClientSummary, int64, error) {
	base := r.db.Conn(ctx).
		Table("client_group_member cgm").
		Joins("JOIN client c ON c.id = cgm.client_id AND c.deleted_at IS NULL").
		Where("cgm.client_group_id = ? AND cgm.business_id = ? AND cgm.deleted_at IS NULL", params.ClientGroupID, params.BusinessID)
	if params.Search != "" {
		like := "%" + params.Search + "%"
		base = base.Where("c.name ILIKE ? OR c.email ILIKE ? OR c.phone ILIKE ? OR c.dni ILIKE ?", like, like, like, like)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []clientRow
	if err := base.
		Select("c.id, c.name, c.email, c.phone, c.dni").
		Order("c.name ASC").
		Offset(params.Offset()).Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	members := make([]entities.ClientSummary, len(rows))
	for i, row := range rows {
		members[i] = clientRowToSummary(row)
	}
	return members, total, nil
}

func (r *Repository) ListAvailableClients(ctx context.Context, params dtos.ListAvailableClientsParams) ([]entities.ClientSummary, int64, error) {
	base := r.db.Conn(ctx).
		Table("client c").
		Joins("LEFT JOIN client_group_member cgm ON cgm.client_id = c.id AND cgm.deleted_at IS NULL").
		Joins("LEFT JOIN client_group cg ON cg.id = cgm.client_group_id AND cg.deleted_at IS NULL").
		Where("c.business_id = ? AND c.deleted_at IS NULL", params.BusinessID)
	if params.OnlyUngrouped {
		base = base.Where("cgm.id IS NULL")
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		base = base.Where("c.name ILIKE ? OR c.email ILIKE ? OR c.phone ILIKE ? OR c.dni ILIKE ?", like, like, like, like)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []clientRow
	if err := base.
		Select("c.id, c.name, c.email, c.phone, c.dni, cgm.client_group_id AS group_id, cg.name AS group_name").
		Order("c.name ASC").
		Offset(params.Offset()).Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	clients := make([]entities.ClientSummary, len(rows))
	for i, row := range rows {
		clients[i] = clientRowToSummary(row)
	}
	return clients, total, nil
}

func (r *Repository) AddGroupMembers(ctx context.Context, dto dtos.AddGroupMembersDTO) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		var groupCount int64
		if err := tx.Model(&models.ClientGroup{}).
			Where("id = ? AND business_id = ?", dto.ClientGroupID, dto.BusinessID).
			Count(&groupCount).Error; err != nil {
			return err
		}
		if groupCount == 0 {
			return domainerrors.ErrGroupNotFound
		}

		var validIDs []uint
		if err := tx.Model(&models.Client{}).
			Where("business_id = ? AND id IN ?", dto.BusinessID, dto.ClientIDs).
			Pluck("id", &validIDs).Error; err != nil {
			return err
		}
		if len(validIDs) == 0 {
			return domainerrors.ErrClientNotFound
		}

		if err := tx.Unscoped().
			Where("business_id = ? AND client_id IN ?", dto.BusinessID, validIDs).
			Delete(&models.ClientGroupMember{}).Error; err != nil {
			return err
		}

		memberships := make([]models.ClientGroupMember, len(validIDs))
		for i, clientID := range validIDs {
			memberships[i] = models.ClientGroupMember{
				BusinessID:    dto.BusinessID,
				ClientGroupID: dto.ClientGroupID,
				ClientID:      clientID,
			}
		}
		return tx.Create(&memberships).Error
	})
}

func (r *Repository) RemoveGroupMember(ctx context.Context, businessID, groupID, clientID uint) error {
	result := r.db.Conn(ctx).Unscoped().
		Where("business_id = ? AND client_group_id = ? AND client_id = ?", businessID, groupID, clientID).
		Delete(&models.ClientGroupMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrClientNotFound
	}
	return nil
}

func (r *Repository) GetClientGroupID(ctx context.Context, businessID, clientID uint) (*uint, error) {
	var member models.ClientGroupMember
	err := r.db.Conn(ctx).
		Where("business_id = ? AND client_id = ?", businessID, clientID).
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member.ClientGroupID, nil
}
