package repository

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSystemProviderFlag(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.IntegrationType{}); err != nil {
		return err
	}

	return r.db.Conn(ctx).Exec(`
		UPDATE integration_types
		SET is_system_provider = true
		WHERE deleted_at IS NULL
		  AND code IN ('envioclick', 'enviame', 'mipaquete', 'shipit')
	`).Error
}
