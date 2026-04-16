package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateAnnouncements(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.AnnouncementCategory{},
		&models.Announcement{},
		&models.AnnouncementImage{},
		&models.AnnouncementLink{},
		&models.AnnouncementTarget{},
		&models.AnnouncementView{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate announcement tables: %w", err)
	}

	r.db.Conn(ctx).Exec("DELETE FROM announcement_images WHERE deleted_at IS NOT NULL")

	if r.db.Conn(ctx).Migrator().HasColumn(&models.AnnouncementImage{}, "deleted_at") {
		r.db.Conn(ctx).Migrator().DropColumn(&models.AnnouncementImage{}, "deleted_at")
	}

	if err := r.seedAnnouncementCategories(ctx); err != nil {
		return fmt.Errorf("failed to seed announcement categories: %w", err)
	}

	return nil
}

func (r *Repository) seedAnnouncementCategories(ctx context.Context) error {
	categories := []models.AnnouncementCategory{
		{Code: "promotion", Name: "Promocion", Icon: "tag", Color: "#10b981"},
		{Code: "alert", Name: "Alerta", Icon: "alert-triangle", Color: "#ef4444"},
		{Code: "informative", Name: "Informativo", Icon: "info", Color: "#3b82f6"},
		{Code: "tutorial", Name: "Tutorial", Icon: "book-open", Color: "#8b5cf6"},
		{Code: "update", Name: "Actualizacion", Icon: "refresh-cw", Color: "#f59e0b"},
		{Code: "terms", Name: "Terminos y Condiciones", Icon: "file-text", Color: "#6b7280"},
	}

	for _, cat := range categories {
		var existing models.AnnouncementCategory
		result := r.db.Conn(ctx).Where("code = ?", cat.Code).First(&existing)
		if result.RowsAffected == 0 {
			if err := r.db.Conn(ctx).Create(&cat).Error; err != nil {
				return fmt.Errorf("failed to seed category %s: %w", cat.Code, err)
			}
		}
	}

	return nil
}
