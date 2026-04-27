package usecaseupdatestatus

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
)

func (uc *UseCaseUpdateStatus) validateStockForPicking(ctx context.Context, order *entities.ProbabilityOrder) error {
	if len(order.OrderItems) == 0 {
		return fmt.Errorf("%w: la orden no tiene productos", domainerrors.ErrInsufficientStock)
	}

	items := make([]dtos.StockCheckItem, 0, len(order.OrderItems))
	for _, it := range order.OrderItems {
		if it.ProductID == nil || *it.ProductID == "" || it.Quantity <= 0 {
			continue
		}
		items = append(items, dtos.StockCheckItem{
			ProductID:  *it.ProductID,
			ProductSKU: it.ProductSKU,
			Quantity:   it.Quantity,
		})
	}

	if len(items) == 0 {
		return fmt.Errorf("%w: la orden no tiene productos validos para verificar inventario", domainerrors.ErrInsufficientStock)
	}

	if order.BusinessID == nil || *order.BusinessID == 0 {
		return fmt.Errorf("%w: la orden no tiene business_id", domainerrors.ErrInsufficientStock)
	}

	results, err := uc.repo.CheckStockForOrder(ctx, *order.BusinessID, order.WarehouseID, items)
	if err != nil {
		return fmt.Errorf("error consultando inventario: %w", err)
	}

	var faltantes []string
	for _, r := range results {
		if !r.Sufficient {
			label := r.ProductSKU
			if label == "" {
				label = r.ProductID
			}
			faltantes = append(faltantes, fmt.Sprintf("%s (requerido %d, disponible %d)", label, r.Required, r.Available))
		}
	}

	if len(faltantes) > 0 {
		return fmt.Errorf("%w: %s", domainerrors.ErrInsufficientStock, strings.Join(faltantes, "; "))
	}

	return nil
}
