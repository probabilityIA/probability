package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// GetIntegrationTypeByIntegrationID obtiene el type_id de una integración
// basado en el ID de la integración.
//
// Tabla consultada: integrations (gestionada por módulo integrations/core)
// Replicado localmente siguiendo regla de aislamiento de repositorios.
// Solo lectura (SELECT) — no modifica estado.
func (r *Repository) GetIntegrationTypeByIntegrationID(ctx context.Context, integrationID uint) (int, error) {
	var typeID int

	err := r.db.Conn(ctx).
		Model(&models.Integration{}).
		Select("integration_type_id").
		Where("id = ?", integrationID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&typeID).Error

	if err == gorm.ErrRecordNotFound {
		r.log.Warn(ctx).
			Uint("integration_id", integrationID).
			Msg("Integration not found")
		return 0, fmt.Errorf("integration not found: %d", integrationID)
	}

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Uint("integration_id", integrationID).
			Msg("Failed to get integration type")
		return 0, fmt.Errorf("failed to get integration type: %w", err)
	}

	r.log.Debug(ctx).
		Uint("integration_id", integrationID).
		Int("type_id", typeID).
		Msg("Got integration type")

	return typeID, nil
}
