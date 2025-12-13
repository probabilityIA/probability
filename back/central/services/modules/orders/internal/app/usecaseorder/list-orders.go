package usecaseorder

import (
	"context"
	"fmt"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder/mapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// ListOrders obtiene una lista paginada de órdenes con filtros
func (uc *UseCaseOrder) ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.OrdersListResponse, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Obtener órdenes del repositorio
	orders, total, err := uc.repo.ListOrders(ctx, page, pageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing orders: %w", err)
	}

	// Mapear a respuestas resumidas
	orderSummaries := make([]domain.OrderSummary, len(orders))
	for i, order := range orders {
		orderSummaries[i] = mapper.ToOrderSummary(&order)
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &domain.OrdersListResponse{
		Data:       orderSummaries,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
