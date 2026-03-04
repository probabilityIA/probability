package app

import (
	"context"
	"fmt"
)

func (uc *useCase) DeleteAllOrders(ctx context.Context, businessID uint) (int64, error) {
	if businessID == 0 {
		return 0, fmt.Errorf("business_id is required")
	}
	return uc.repo.DeleteAllOrders(ctx, businessID)
}
