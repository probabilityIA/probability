package usecaseordermapping

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// ───────────────────────────────────────────
//
//	ORDER MAPPING USE CASE INTERFACE
//
// ───────────────────────────────────────────

// IOrderMappingUseCase define la interfaz para mapear y guardar órdenes canónicas
type IOrderMappingUseCase interface {
	// MapAndSaveOrder recibe una orden en formato canónico y la guarda en todas las tablas relacionadas
	MapAndSaveOrder(ctx context.Context, dto *domain.CanonicalOrderDTO) (*domain.OrderResponse, error)
}

// ───────────────────────────────────────────
//
//	ORDER MAPPING USE CASE IMPLEMENTATION
//
// ───────────────────────────────────────────

// UseCaseOrderMapping contiene el caso de uso para mapear y guardar órdenes canónicas
type UseCaseOrderMapping struct {
	repo domain.IRepository
}

// New crea una nueva instancia de UseCaseOrderMapping
func New(repo domain.IRepository) IOrderMappingUseCase {
	return &UseCaseOrderMapping{
		repo: repo,
	}
}
