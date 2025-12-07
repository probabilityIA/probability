package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/app/usecaseordermapping"
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// UseCases contiene todos los casos de uso del módulo orders
// Mantiene compatibilidad con el código existente y expone los nuevos casos de uso
type UseCases struct {
	repo domain.IRepository

	// Casos de uso modulares
	OrderCRUD    *usecaseorder.UseCaseOrder
	OrderMapping usecaseordermapping.IOrderMappingUseCase
}

// New crea una nueva instancia de UseCases
func New(repo domain.IRepository) *UseCases {
	return &UseCases{
		repo:         repo,
		OrderCRUD:    usecaseorder.New(repo),
		OrderMapping: usecaseordermapping.New(repo),
	}
}

// ───────────────────────────────────────────
// MÉTODOS DE COMPATIBILIDAD - Delegar al CRUD
// ───────────────────────────────────────────

// CreateOrder delega al caso de uso CRUD
func (uc *UseCases) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.OrderResponse, error) {
	return uc.OrderCRUD.CreateOrder(ctx, req)
}

// GetOrderByID delega al caso de uso CRUD
func (uc *UseCases) GetOrderByID(ctx context.Context, id string) (*domain.OrderResponse, error) {
	return uc.OrderCRUD.GetOrderByID(ctx, id)
}

// ListOrders delega al caso de uso CRUD
func (uc *UseCases) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.OrdersListResponse, error) {
	return uc.OrderCRUD.ListOrders(ctx, page, pageSize, filters)
}

// UpdateOrder delega al caso de uso CRUD
func (uc *UseCases) UpdateOrder(ctx context.Context, id string, req *domain.UpdateOrderRequest) (*domain.OrderResponse, error) {
	return uc.OrderCRUD.UpdateOrder(ctx, id, req)
}

// DeleteOrder delega al caso de uso CRUD
func (uc *UseCases) DeleteOrder(ctx context.Context, id string) error {
	return uc.OrderCRUD.DeleteOrder(ctx, id)
}
