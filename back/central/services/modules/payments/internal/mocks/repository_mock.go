package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/ports"
)

type RepositoryMock struct {
	CreatePaymentMethodFn        func(ctx context.Context, method *entities.PaymentMethod) error
	GetPaymentMethodByIDFn       func(ctx context.Context, id uint) (*entities.PaymentMethod, error)
	GetPaymentMethodByCodeFn     func(ctx context.Context, code string) (*entities.PaymentMethod, error)
	ListPaymentMethodsFn         func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.PaymentMethod, int64, error)
	UpdatePaymentMethodFn        func(ctx context.Context, method *entities.PaymentMethod) error
	DeletePaymentMethodFn        func(ctx context.Context, id uint) error
	TogglePaymentMethodActiveFn  func(ctx context.Context, id uint) (*entities.PaymentMethod, error)
	PaymentMethodExistsFn        func(ctx context.Context, code string) (bool, error)
	PaymentMethodHasActiveMapsFn func(ctx context.Context, id uint) (bool, error)

	CreatePaymentMappingFn                        func(ctx context.Context, mapping *entities.PaymentMethodMapping) error
	GetPaymentMappingByIDFn                       func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error)
	GetPaymentMappingByIDWithMethodFn             func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error)
	ListPaymentMappingsFn                         func(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error)
	ListPaymentMappingsWithMethodsFn              func(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error)
	UpdatePaymentMappingFn                        func(ctx context.Context, mapping *entities.PaymentMethodMapping) error
	DeletePaymentMappingFn                        func(ctx context.Context, id uint) error
	GetPaymentMappingsByIntegrationTypeFn         func(ctx context.Context, integrationType string) ([]entities.PaymentMethodMapping, error)
	GetPaymentMappingsByIntegrationTypeWithMethFn func(ctx context.Context, integrationType string) ([]entities.PaymentMethodMapping, error)
	TogglePaymentMappingActiveFn                  func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error)
	PaymentMappingExistsFn                        func(ctx context.Context, integrationType, originalMethod string) (bool, error)

	ListPaymentStatusesFn       func(ctx context.Context, isActive *bool) ([]entities.PaymentStatus, error)
	ListChannelPaymentMethodsFn func(ctx context.Context, integrationType *string, isActive *bool) ([]entities.ChannelPaymentMethod, error)
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) CreatePaymentMethod(ctx context.Context, method *entities.PaymentMethod) error {
	if m.CreatePaymentMethodFn != nil {
		return m.CreatePaymentMethodFn(ctx, method)
	}
	return nil
}

func (m *RepositoryMock) GetPaymentMethodByID(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
	if m.GetPaymentMethodByIDFn != nil {
		return m.GetPaymentMethodByIDFn(ctx, id)
	}
	return &entities.PaymentMethod{ID: id}, nil
}

func (m *RepositoryMock) GetPaymentMethodByCode(ctx context.Context, code string) (*entities.PaymentMethod, error) {
	if m.GetPaymentMethodByCodeFn != nil {
		return m.GetPaymentMethodByCodeFn(ctx, code)
	}
	return &entities.PaymentMethod{Code: code}, nil
}

func (m *RepositoryMock) ListPaymentMethods(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.PaymentMethod, int64, error) {
	if m.ListPaymentMethodsFn != nil {
		return m.ListPaymentMethodsFn(ctx, page, pageSize, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) UpdatePaymentMethod(ctx context.Context, method *entities.PaymentMethod) error {
	if m.UpdatePaymentMethodFn != nil {
		return m.UpdatePaymentMethodFn(ctx, method)
	}
	return nil
}

func (m *RepositoryMock) DeletePaymentMethod(ctx context.Context, id uint) error {
	if m.DeletePaymentMethodFn != nil {
		return m.DeletePaymentMethodFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) TogglePaymentMethodActive(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
	if m.TogglePaymentMethodActiveFn != nil {
		return m.TogglePaymentMethodActiveFn(ctx, id)
	}
	return &entities.PaymentMethod{ID: id}, nil
}

func (m *RepositoryMock) PaymentMethodExists(ctx context.Context, code string) (bool, error) {
	if m.PaymentMethodExistsFn != nil {
		return m.PaymentMethodExistsFn(ctx, code)
	}
	return false, nil
}

func (m *RepositoryMock) PaymentMethodHasActiveMappings(ctx context.Context, id uint) (bool, error) {
	if m.PaymentMethodHasActiveMapsFn != nil {
		return m.PaymentMethodHasActiveMapsFn(ctx, id)
	}
	return false, nil
}

func (m *RepositoryMock) CreatePaymentMapping(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
	if m.CreatePaymentMappingFn != nil {
		return m.CreatePaymentMappingFn(ctx, mapping)
	}
	return nil
}

func (m *RepositoryMock) GetPaymentMappingByID(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
	if m.GetPaymentMappingByIDFn != nil {
		return m.GetPaymentMappingByIDFn(ctx, id)
	}
	return &entities.PaymentMethodMapping{ID: id}, nil
}

func (m *RepositoryMock) GetPaymentMappingByIDWithMethod(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
	if m.GetPaymentMappingByIDWithMethodFn != nil {
		return m.GetPaymentMappingByIDWithMethodFn(ctx, id)
	}
	return &entities.PaymentMethodMapping{ID: id}, nil
}

func (m *RepositoryMock) ListPaymentMappings(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error) {
	if m.ListPaymentMappingsFn != nil {
		return m.ListPaymentMappingsFn(ctx, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) ListPaymentMappingsWithMethods(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error) {
	if m.ListPaymentMappingsWithMethodsFn != nil {
		return m.ListPaymentMappingsWithMethodsFn(ctx, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) UpdatePaymentMapping(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
	if m.UpdatePaymentMappingFn != nil {
		return m.UpdatePaymentMappingFn(ctx, mapping)
	}
	return nil
}

func (m *RepositoryMock) DeletePaymentMapping(ctx context.Context, id uint) error {
	if m.DeletePaymentMappingFn != nil {
		return m.DeletePaymentMappingFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) GetPaymentMappingsByIntegrationType(ctx context.Context, integrationType string) ([]entities.PaymentMethodMapping, error) {
	if m.GetPaymentMappingsByIntegrationTypeFn != nil {
		return m.GetPaymentMappingsByIntegrationTypeFn(ctx, integrationType)
	}
	return nil, nil
}

func (m *RepositoryMock) GetPaymentMappingsByIntegrationTypeWithMethods(ctx context.Context, integrationType string) ([]entities.PaymentMethodMapping, error) {
	if m.GetPaymentMappingsByIntegrationTypeWithMethFn != nil {
		return m.GetPaymentMappingsByIntegrationTypeWithMethFn(ctx, integrationType)
	}
	return nil, nil
}

func (m *RepositoryMock) TogglePaymentMappingActive(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
	if m.TogglePaymentMappingActiveFn != nil {
		return m.TogglePaymentMappingActiveFn(ctx, id)
	}
	return &entities.PaymentMethodMapping{ID: id}, nil
}

func (m *RepositoryMock) PaymentMappingExists(ctx context.Context, integrationType, originalMethod string) (bool, error) {
	if m.PaymentMappingExistsFn != nil {
		return m.PaymentMappingExistsFn(ctx, integrationType, originalMethod)
	}
	return false, nil
}

func (m *RepositoryMock) ListPaymentStatuses(ctx context.Context, isActive *bool) ([]entities.PaymentStatus, error) {
	if m.ListPaymentStatusesFn != nil {
		return m.ListPaymentStatusesFn(ctx, isActive)
	}
	return nil, nil
}

func (m *RepositoryMock) ListChannelPaymentMethods(ctx context.Context, integrationType *string, isActive *bool) ([]entities.ChannelPaymentMethod, error) {
	if m.ListChannelPaymentMethodsFn != nil {
		return m.ListChannelPaymentMethodsFn(ctx, integrationType, isActive)
	}
	return nil, nil
}
