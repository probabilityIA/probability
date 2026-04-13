package repository

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) Create(ctx context.Context, announcement *entities.Announcement) (*entities.Announcement, error) {
	model := announcementFromEntity(announcement)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ctx, model.ID)
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*entities.Announcement, error) {
	var model models.Announcement
	err := r.db.Conn(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Links", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Targets").
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrAnnouncementNotFound
		}
		return nil, err
	}
	return announcementToEntity(&model), nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListAnnouncementsParams) ([]entities.Announcement, int64, error) {
	var total int64

	countQuery := r.db.Conn(ctx).Model(&models.Announcement{})
	if params.Status != "" {
		countQuery = countQuery.Where("status = ?", params.Status)
	}
	if params.CategoryID != nil {
		countQuery = countQuery.Where("category_id = ?", *params.CategoryID)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		countQuery = countQuery.Where("title ILIKE ? OR message ILIKE ?", like, like)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.Announcement
	query := r.db.Conn(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Links", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Targets")

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("title ILIKE ? OR message ILIKE ?", like, like)
	}

	offset := params.Offset()
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(params.PageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entities.Announcement, len(rows))
	for i, row := range rows {
		result[i] = *announcementToEntity(&row)
	}
	return result, total, nil
}

func (r *Repository) Update(ctx context.Context, announcement *entities.Announcement) (*entities.Announcement, error) {
	model := announcementFromEntity(announcement)
	model.ID = announcement.ID

	err := r.db.Conn(ctx).Model(&models.Announcement{}).Where("id = ?", announcement.ID).
		Updates(map[string]interface{}{
			"category_id":    model.CategoryID,
			"title":          model.Title,
			"message":        model.Message,
			"display_type":   model.DisplayType,
			"frequency_type": model.FrequencyType,
			"is_global":      model.IsGlobal,
			"status":         model.Status,
			"starts_at":      model.StartsAt,
			"ends_at":        model.EndsAt,
			"force_redisplay": model.ForceRedisplay,
		}).Error
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, announcement.ID)
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	result := r.db.Conn(ctx).Where("id = ?", id).Delete(&models.Announcement{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrAnnouncementNotFound
	}
	return nil
}

func (r *Repository) ChangeStatus(ctx context.Context, id uint, status entities.AnnouncementStatus) error {
	return r.db.Conn(ctx).Model(&models.Announcement{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}
