package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) UpsertDeviceToken(ctx context.Context, dto dtos.RegisterDeviceDTO) error {
	now := time.Now()

	var existing models.DeviceToken
	err := r.db.Conn(ctx).
		Where("token = ?", dto.Token).
		First(&existing).Error

	if err == nil {
		return r.db.Conn(ctx).Model(&models.DeviceToken{}).
			Where("id = ?", existing.ID).
			Updates(map[string]any{
				"user_id":      dto.UserID,
				"business_id":  dto.BusinessID,
				"platform":     dto.Platform,
				"app_version":  dto.AppVersion,
				"device_name":  dto.DeviceName,
				"is_active":    true,
				"last_seen_at": now,
				"updated_at":   now,
				"deleted_at":   nil,
			}).Error
	}

	device := models.DeviceToken{
		UserID:     dto.UserID,
		BusinessID: dto.BusinessID,
		Token:      dto.Token,
		Platform:   dto.Platform,
		AppVersion: dto.AppVersion,
		DeviceName: dto.DeviceName,
		IsActive:   true,
		LastSeenAt: &now,
	}
	return r.db.Conn(ctx).Create(&device).Error
}

func (r *Repository) DeactivateToken(ctx context.Context, userID uint, token string) error {
	return r.db.Conn(ctx).Model(&models.DeviceToken{}).
		Where("user_id = ? AND token = ?", userID, token).
		Update("is_active", false).Error
}

func (r *Repository) DeactivateTokens(ctx context.Context, tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}
	return r.db.Conn(ctx).Model(&models.DeviceToken{}).
		Where("token IN ?", tokens).
		Update("is_active", false).Error
}

func (r *Repository) ListActiveTokensByBusiness(ctx context.Context, businessID uint) ([]entities.DeviceToken, error) {
	var rows []models.DeviceToken
	err := r.db.Conn(ctx).
		Where("business_id = ? AND is_active = true AND deleted_at IS NULL", businessID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return toEntities(rows), nil
}

func (r *Repository) ListDevicesByUser(ctx context.Context, userID uint) ([]entities.DeviceToken, error) {
	var rows []models.DeviceToken
	err := r.db.Conn(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("last_seen_at DESC NULLS LAST").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return toEntities(rows), nil
}

func (r *Repository) UserBelongsToBusiness(ctx context.Context, userID, businessID uint) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Table("user_businesses").
		Where("user_id = ? AND business_id = ?", userID, businessID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	err = r.db.Conn(ctx).Model(&models.BusinessStaff{}).
		Where("user_id = ? AND business_id = ? AND deleted_at IS NULL", userID, businessID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func toEntities(rows []models.DeviceToken) []entities.DeviceToken {
	out := make([]entities.DeviceToken, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.DeviceToken{
			ID:         row.ID,
			UserID:     row.UserID,
			BusinessID: row.BusinessID,
			Token:      row.Token,
			Platform:   row.Platform,
			AppVersion: row.AppVersion,
			DeviceName: row.DeviceName,
			IsActive:   row.IsActive,
			LastSeenAt: row.LastSeenAt,
			CreatedAt:  row.CreatedAt,
		})
	}
	return out
}
