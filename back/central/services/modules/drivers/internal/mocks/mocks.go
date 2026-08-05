package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/ports"
)

type ExistsCall struct {
	BusinessID     uint
	Identification string
	ExcludeID      *uint
}

type RepositoryMock struct {
	CreateFn                 func(ctx context.Context, driver *entities.Driver) (*entities.Driver, error)
	GetByIDFn                func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error)
	ListFn                   func(ctx context.Context, params dtos.ListDriversParams) ([]entities.Driver, int64, error)
	UpdateFn                 func(ctx context.Context, driver *entities.Driver) (*entities.Driver, error)
	DeleteFn                 func(ctx context.Context, businessID, driverID uint) error
	ExistsByIdentificationFn func(ctx context.Context, businessID uint, identification string, excludeID *uint) (bool, error)

	Created     *entities.Driver
	Updated     *entities.Driver
	ExistsCalls []ExistsCall
	Deleted     []uint
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) Create(ctx context.Context, driver *entities.Driver) (*entities.Driver, error) {
	m.Created = driver
	if m.CreateFn != nil {
		return m.CreateFn(ctx, driver)
	}
	driver.ID = 1
	return driver, nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, businessID, driverID)
	}
	return &entities.Driver{ID: driverID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) List(ctx context.Context, params dtos.ListDriversParams) ([]entities.Driver, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) Update(ctx context.Context, driver *entities.Driver) (*entities.Driver, error) {
	m.Updated = driver
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, driver)
	}
	return driver, nil
}

func (m *RepositoryMock) Delete(ctx context.Context, businessID, driverID uint) error {
	m.Deleted = append(m.Deleted, driverID)
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, businessID, driverID)
	}
	return nil
}

func (m *RepositoryMock) ExistsByIdentification(ctx context.Context, businessID uint, identification string, excludeID *uint) (bool, error) {
	m.ExistsCalls = append(m.ExistsCalls, ExistsCall{businessID, identification, excludeID})
	if m.ExistsByIdentificationFn != nil {
		return m.ExistsByIdentificationFn(ctx, businessID, identification, excludeID)
	}
	return false, nil
}
