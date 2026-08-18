package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateTrazabilidadUsuario(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.Shipment{},
		&models.ShippingQuote{},
		&models.Order{},
	); err != nil {
		return fmt.Errorf("automigrate trazabilidad de usuario: %w", err)
	}
	return nil
}
