package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/ports"
)

type IUseCase interface {
	CreateDriver(ctx context.Context, dto dtos.CreateDriverDTO) (*entities.Driver, error)
	GetDriver(ctx context.Context, businessID, driverID uint) (*entities.Driver, error)
	ListDrivers(ctx context.Context, params dtos.ListDriversParams) ([]entities.Driver, int64, error)
	UpdateDriver(ctx context.Context, dto dtos.UpdateDriverDTO) (*entities.Driver, error)
	DeleteDriver(ctx context.Context, businessID, driverID uint) error
}

type UseCase struct {
	repo ports.IRepository
}

func New(repo ports.IRepository) IUseCase {
	return &UseCase{repo: repo}
}
