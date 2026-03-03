package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

// IInventoryCache define la interfaz para el servicio de caché de inventario.
// Inyectada en el repositorio para write-through/read-through.
// Resiliente: retorna nil en cache miss (no error), funciona sin Redis.
type IInventoryCache interface {
	GetProductLevels(ctx context.Context, productID string, businessID uint) ([]entities.InventoryLevel, error)
	SetProductLevels(ctx context.Context, productID string, businessID uint, levels []entities.InventoryLevel) error
	InvalidateProduct(ctx context.Context, productID string, businessID uint) error
	GetLevel(ctx context.Context, productID string, warehouseID uint) (*entities.InventoryLevel, error)
	SetLevel(ctx context.Context, productID string, warehouseID uint, level *entities.InventoryLevel) error
	InvalidateLevel(ctx context.Context, productID string, warehouseID uint) error
}
