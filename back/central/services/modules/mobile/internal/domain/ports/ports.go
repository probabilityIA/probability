package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/entities"
)

type IRepository interface {
	GetOrderFull(ctx context.Context, businessID uint, orderID string) (*entities.OrderFull, error)
}
