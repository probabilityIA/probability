package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
)

func (uc *UseCase) GetDriver(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
	return uc.repo.GetByID(ctx, businessID, driverID)
}
