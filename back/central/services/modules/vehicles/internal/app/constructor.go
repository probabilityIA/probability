package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/ports"
)

type IUseCase interface {
	CreateVehicle(ctx context.Context, dto dtos.CreateVehicleDTO) (*entities.Vehicle, error)
	GetVehicle(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error)
	ListVehicles(ctx context.Context, params dtos.ListVehiclesParams) ([]entities.Vehicle, int64, error)
	UpdateVehicle(ctx context.Context, dto dtos.UpdateVehicleDTO) (*entities.Vehicle, error)
	DeleteVehicle(ctx context.Context, businessID, vehicleID uint) error
}

type UseCase struct {
	repo ports.IRepository
}

func New(repo ports.IRepository) IUseCase {
	return &UseCase{repo: repo}
}
