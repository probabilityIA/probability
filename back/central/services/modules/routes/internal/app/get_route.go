package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

func (uc *UseCase) GetRoute(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
	return uc.repo.GetRouteByID(ctx, businessID, routeID)
}
