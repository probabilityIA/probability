package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	catalogoIntegrationTypeID = 30
	catalogoDisplayName       = "Catalogo"
	catalogoDescription       = "Catalogo con login de cliente final: el negocio carga el stock disponible y sus clientes hacen pedidos limitados a esa cantidad."
)

func (r *Repository) migrateCatalogoRename(ctx context.Context) error {
	if err := r.db.Conn(ctx).Model(&models.IntegrationType{}).
		Where("id = ?", catalogoIntegrationTypeID).
		Updates(map[string]interface{}{
			"name":        catalogoDisplayName,
			"description": catalogoDescription,
		}).Error; err != nil {
		return fmt.Errorf("migrateCatalogoRename: %w", err)
	}
	return nil
}
