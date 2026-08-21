package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/entities"
)

func (u *UseCase) GetOrderFull(ctx context.Context, businessID uint, orderID string) (*entities.OrderFull, error) {
	return u.repo.GetOrderFull(ctx, businessID, orderID)
}
