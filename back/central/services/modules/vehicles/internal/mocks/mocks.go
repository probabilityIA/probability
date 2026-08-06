package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/ports"
)

type ExistsCall struct {
	BusinessID uint
	Plate      string
	ExcludeID  *uint
}

type RepositoryMock struct {
	CreateFn               func(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error)
	GetByIDFn              func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error)
	ListFn                 func(ctx context.Context, params dtos.ListVehiclesParams) ([]entities.Vehicle, int64, error)
	UpdateFn               func(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error)
	DeleteFn               func(ctx context.Context, businessID, vehicleID uint) error
	ExistsByLicensePlateFn func(ctx context.Context, businessID uint, plate string, excludeID *uint) (bool, error)

	Created     *entities.Vehicle
	Updated     *entities.Vehicle
	ExistsCalls []ExistsCall
	Deleted     []uint
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) Create(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error) {
	m.Created = vehicle
	if m.CreateFn != nil {
		return m.CreateFn(ctx, vehicle)
	}
	vehicle.ID = 1
	return vehicle, nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, businessID, vehicleID)
	}
	return &entities.Vehicle{ID: vehicleID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) List(ctx context.Context, params dtos.ListVehiclesParams) ([]entities.Vehicle, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) Update(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error) {
	m.Updated = vehicle
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, vehicle)
	}
	return vehicle, nil
}

func (m *RepositoryMock) Delete(ctx context.Context, businessID, vehicleID uint) error {
	m.Deleted = append(m.Deleted, vehicleID)
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, businessID, vehicleID)
	}
	return nil
}

func (m *RepositoryMock) ExistsByLicensePlate(ctx context.Context, businessID uint, plate string, excludeID *uint) (bool, error) {
	m.ExistsCalls = append(m.ExistsCalls, ExistsCall{businessID, plate, excludeID})
	if m.ExistsByLicensePlateFn != nil {
		return m.ExistsByLicensePlateFn(ctx, businessID, plate, excludeID)
	}
	return false, nil
}
