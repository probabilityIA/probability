package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const prunePorCampo = `
DELETE FROM product_field_changes
WHERE id IN (
    SELECT id FROM product_field_changes
    WHERE product_id = ? AND field = ?
    ORDER BY created_at DESC, id DESC
    OFFSET ?
)`

func (r *Repository) RecordFieldChanges(ctx context.Context, changes []domain.FieldChange) error {
	if len(changes) == 0 {
		return nil
	}

	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		for _, change := range changes {
			historia := models.ProductFieldChange{
				ProductID:     change.ProductID,
				Field:         change.Field,
				BusinessID:    change.BusinessID,
				OldValue:      change.OldValue,
				Truncated:     change.Truncated,
				Source:        change.Origin.Source,
				IntegrationID: change.Origin.IntegrationID,
				UserID:        change.Origin.UserID,
				BatchID:       change.BatchID,
				CreatedAt:     change.ChangedAt,
			}
			if err := tx.Create(&historia).Error; err != nil {
				return err
			}

			if err := tx.Exec(prunePorCampo, change.ProductID, change.Field, domain.MaxChangesPerField).Error; err != nil {
				return err
			}

			origen := models.ProductFieldOrigin{
				ProductID:     change.ProductID,
				Field:         change.Field,
				BusinessID:    change.BusinessID,
				Source:        change.Origin.Source,
				IntegrationID: change.Origin.IntegrationID,
				UserID:        change.Origin.UserID,
				ChangedAt:     change.ChangedAt,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "product_id"}, {Name: "field"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"source", "integration_id", "user_id", "changed_at", "updated_at", "deleted_at",
				}),
			}).Create(&origen).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

const origenesQuery = `
SELECT
    o.field,
    o.source,
    o.integration_id,
    COALESCE(i.name, '')       AS integration_name,
    o.user_id,
    COALESCE(u.name, '')       AS user_name,
    o.changed_at,
    COALESCE(h.old_value, '')  AS previous_value,
    h.created_at               AS previous_at
FROM product_field_origins o
LEFT JOIN integrations i ON i.id = o.integration_id AND i.deleted_at IS NULL
LEFT JOIN "user" u ON u.id = o.user_id AND u.deleted_at IS NULL
LEFT JOIN LATERAL (
    SELECT c.old_value, c.created_at
    FROM product_field_changes c
    WHERE c.product_id = o.product_id AND c.field = o.field
    ORDER BY c.created_at DESC, c.id DESC
    LIMIT 1
) h ON true
WHERE o.product_id = ? AND o.business_id = ? AND o.deleted_at IS NULL
ORDER BY o.changed_at DESC`

func (r *Repository) GetProductFieldOrigins(ctx context.Context, businessID uint, productID string) ([]domain.FieldOriginView, error) {
	var rows []struct {
		Field           string
		Source          string
		IntegrationID   *uint
		IntegrationName string
		UserID          *uint
		UserName        string
		ChangedAt       time.Time
		PreviousValue   string
		PreviousAt      *time.Time
	}

	if err := r.db.Conn(ctx).Raw(origenesQuery, productID, businessID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	origins := make([]domain.FieldOriginView, 0, len(rows))
	for _, row := range rows {
		origins = append(origins, domain.FieldOriginView{
			Field:           row.Field,
			Source:          row.Source,
			IntegrationID:   row.IntegrationID,
			IntegrationName: row.IntegrationName,
			UserID:          row.UserID,
			UserName:        row.UserName,
			ChangedAt:       row.ChangedAt,
			PreviousValue:   row.PreviousValue,
			PreviousAt:      row.PreviousAt,
		})
	}
	return origins, nil
}
