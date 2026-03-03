package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// RepositoryMock - Mock del repositorio de configuraciones de notificaciones
type RepositoryMock struct {
	CreateFn                                   func(ctx context.Context, config *entities.IntegrationNotificationConfig) error
	UpdateFn                                   func(ctx context.Context, config *entities.IntegrationNotificationConfig) error
	GetByIDFn                                  func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error)
	ListFn                                     func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error)
	DeleteFn                                   func(ctx context.Context, id uint) error
	GetActiveConfigsByIntegrationAndTriggerFn  func(ctx context.Context, integrationID uint, trigger string) ([]entities.IntegrationNotificationConfig, error)
	SyncConfigsFn                              func(ctx context.Context, businessID uint, integrationID uint, toCreate []*entities.IntegrationNotificationConfig, toUpdate []*entities.IntegrationNotificationConfig, toDeleteIDs []uint) error
}

func (m *RepositoryMock) Create(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, config)
	}
	// Comportamiento por defecto: asignar ID 1
	config.ID = 1
	return nil
}

func (m *RepositoryMock) Update(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, config)
	}
	return nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *RepositoryMock) List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filters)
	}
	return []entities.IntegrationNotificationConfig{}, nil
}

func (m *RepositoryMock) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]entities.IntegrationNotificationConfig, error) {
	if m.GetActiveConfigsByIntegrationAndTriggerFn != nil {
		return m.GetActiveConfigsByIntegrationAndTriggerFn(ctx, integrationID, trigger)
	}
	return []entities.IntegrationNotificationConfig{}, nil
}

func (m *RepositoryMock) SyncConfigs(ctx context.Context, businessID uint, integrationID uint, toCreate []*entities.IntegrationNotificationConfig, toUpdate []*entities.IntegrationNotificationConfig, toDeleteIDs []uint) error {
	if m.SyncConfigsFn != nil {
		return m.SyncConfigsFn(ctx, businessID, integrationID, toCreate, toUpdate, toDeleteIDs)
	}
	return nil
}
