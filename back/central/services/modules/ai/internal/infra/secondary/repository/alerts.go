package repository

import (
	"context"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) SaveAlert(ctx context.Context, alert entities.Alert) (bool, error) {
	row := models.AssistantAlert{
		ID:               alert.ID,
		BusinessID:       alert.BusinessID,
		EventID:          alert.EventID,
		EventType:        alert.EventType,
		Severity:         alert.Severity,
		Title:            alert.Title,
		Body:             alert.Body,
		DestinationKey:   alert.DestinationKey,
		DestinationRoute: alert.DestinationRoute,
		ReferenceType:    alert.ReferenceType,
		ReferenceID:      alert.ReferenceID,
		CreatedAt:        alert.CreatedAt,
	}
	result := r.db.Conn(ctx).Omit(clause.Associations).Clauses(clause.OnConflict{DoNothing: true}).Create(&row)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *Repository) seenAt(ctx context.Context, businessID, userID uint) (time.Time, error) {
	var cursor models.AssistantAlertCursor
	err := r.db.Conn(ctx).Where("user_id = ? AND business_id = ?", userID, businessID).First(&cursor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return time.Time{}, nil
	}
	return cursor.SeenAt, err
}

func (r *Repository) ListAlerts(ctx context.Context, query dtos.AlertQuery) ([]entities.Alert, int64, error) {
	seen, err := r.seenAt(ctx, query.BusinessID, query.UserID)
	if err != nil {
		return nil, 0, err
	}
	base := r.db.Conn(ctx).Model(&models.AssistantAlert{}).Where("business_id = ?", query.BusinessID)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []models.AssistantAlert
	if err := base.Order("created_at DESC").Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	items := make([]entities.Alert, 0, len(rows))
	for _, row := range rows {
		items = append(items, toAlert(row, seen))
	}
	return items, total, nil
}

func (r *Repository) CountUnread(ctx context.Context, businessID, userID uint) (int64, *entities.Alert, error) {
	seen, err := r.seenAt(ctx, businessID, userID)
	if err != nil {
		return 0, nil, err
	}
	scope := r.db.Conn(ctx).Model(&models.AssistantAlert{}).Where("business_id = ? AND created_at > ?", businessID, seen)
	var count int64
	if err := scope.Count(&count).Error; err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, nil, nil
	}
	var row models.AssistantAlert
	if err := scope.Order("created_at DESC").First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return count, nil, nil
		}
		return 0, nil, err
	}
	latest := toAlert(row, seen)
	return count, &latest, nil
}

func (r *Repository) MarkSeen(ctx context.Context, businessID, userID uint, at time.Time) error {
	cursor := models.AssistantAlertCursor{UserID: userID, BusinessID: businessID, SeenAt: at, UpdatedAt: at}
	return r.db.Conn(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "business_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"seen_at", "updated_at"}),
	}).Create(&cursor).Error
}

func (r *Repository) RecentAlerts(ctx context.Context, businessID uint, limit int) ([]entities.Alert, error) {
	var rows []models.AssistantAlert
	err := r.db.Conn(ctx).Where("business_id = ?", businessID).Order("created_at DESC").Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	items := make([]entities.Alert, 0, len(rows))
	for _, row := range rows {
		items = append(items, toAlert(row, time.Time{}))
	}
	return items, nil
}

func (r *Repository) DeleteAlertsOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	result := r.db.Conn(ctx).Where("created_at < ?", cutoff).Delete(&models.AssistantAlert{})
	return result.RowsAffected, result.Error
}

func toAlert(row models.AssistantAlert, seen time.Time) entities.Alert {
	return entities.Alert{
		ID:               row.ID,
		BusinessID:       row.BusinessID,
		EventID:          row.EventID,
		EventType:        row.EventType,
		Severity:         row.Severity,
		Title:            row.Title,
		Body:             row.Body,
		DestinationKey:   row.DestinationKey,
		DestinationRoute: row.DestinationRoute,
		ReferenceType:    row.ReferenceType,
		ReferenceID:      row.ReferenceID,
		CreatedAt:        row.CreatedAt,
		Unread:           row.CreatedAt.After(seen),
	}
}

func (r *Repository) LastAlertOfType(ctx context.Context, businessID uint, eventType string) (*entities.Alert, error) {
	var row models.AssistantAlert
	err := r.db.Conn(ctx).Where("business_id = ? AND event_type = ?", businessID, eventType).Order("created_at DESC").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	alert := toAlert(row, time.Time{})
	return &alert, nil
}
