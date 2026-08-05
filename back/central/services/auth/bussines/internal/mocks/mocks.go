package mocks

import (
	"context"
	"mime/multipart"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type BusinessRepositoryMock struct {
	GetBusinessesFn             func(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.Business, int64, error)
	GetBusinessByIDFn           func(ctx context.Context, id uint) (*domain.Business, error)
	GetBusinessByCodeFn         func(ctx context.Context, code string) (*domain.Business, error)
	GetBusinessByCustomDomainFn func(ctx context.Context, d string) (*domain.Business, error)
	CreateBusinessFn            func(ctx context.Context, business domain.Business) (uint, error)
	UpdateBusinessFn            func(ctx context.Context, id uint, business domain.Business) (string, error)
	DeleteBusinessFn            func(ctx context.Context, id uint) (string, error)

	GetBusinessTypesFn      func(ctx context.Context) ([]domain.BusinessType, error)
	GetBusinessTypeByIDFn   func(ctx context.Context, id uint) (*domain.BusinessType, error)
	GetBusinessTypeByCodeFn func(ctx context.Context, code string) (*domain.BusinessType, error)
	GetBusinessTypeByNameFn func(ctx context.Context, name string) (*domain.BusinessType, error)
	CreateBusinessTypeFn    func(ctx context.Context, businessType domain.BusinessType) (string, error)
	UpdateBusinessTypeFn    func(ctx context.Context, id uint, businessType domain.BusinessType) (string, error)
	DeleteBusinessTypeFn    func(ctx context.Context, id uint) (string, error)

	GetBusinessTypeResourcesPermittedFn func(ctx context.Context, businessTypeID uint) ([]domain.BusinessTypeResourcePermitted, error)
	GetResourceByIDFn                   func(ctx context.Context, resourceID uint) (*domain.Resource, error)

	GetBusinessesWithConfiguredResourcesPaginatedFn func(ctx context.Context, page, perPage int, businessID *uint, businessTypeID *uint) ([]domain.BusinessWithConfiguredResourcesResponse, int64, error)
	GetBusinessByIDWithConfiguredResourcesFn        func(ctx context.Context, businessID uint) (*domain.BusinessWithConfiguredResourcesResponse, error)
	ToggleBusinessResourceActiveFn                  func(ctx context.Context, businessID uint, resourceID uint, active bool) error
	ToggleBusinessActiveFn                          func(ctx context.Context, businessID uint, active bool) error
	CreatePlatformIntegrationFn                     func(ctx context.Context, businessID uint) error

	GetExistingOrderPrefixesFn func(ctx context.Context) ([]string, error)

	GetFiscalProfileFn    func(ctx context.Context, businessID uint) (*domain.FiscalProfile, error)
	UpsertFiscalProfileFn func(ctx context.Context, profile *domain.FiscalProfile) error
}

var _ domain.IBusinessRepository = (*BusinessRepositoryMock)(nil)

func (m *BusinessRepositoryMock) GetBusinesses(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.Business, int64, error) {
	if m.GetBusinessesFn != nil {
		return m.GetBusinessesFn(ctx, page, perPage, name, businessTypeID, isActive)
	}
	return nil, 0, nil
}

func (m *BusinessRepositoryMock) GetBusinessByID(ctx context.Context, id uint) (*domain.Business, error) {
	if m.GetBusinessByIDFn != nil {
		return m.GetBusinessByIDFn(ctx, id)
	}
	return &domain.Business{ID: id}, nil
}

func (m *BusinessRepositoryMock) GetBusinessByCode(ctx context.Context, code string) (*domain.Business, error) {
	if m.GetBusinessByCodeFn != nil {
		return m.GetBusinessByCodeFn(ctx, code)
	}
	return nil, nil
}

func (m *BusinessRepositoryMock) GetBusinessByCustomDomain(ctx context.Context, d string) (*domain.Business, error) {
	if m.GetBusinessByCustomDomainFn != nil {
		return m.GetBusinessByCustomDomainFn(ctx, d)
	}
	return nil, nil
}

func (m *BusinessRepositoryMock) CreateBusiness(ctx context.Context, business domain.Business) (uint, error) {
	if m.CreateBusinessFn != nil {
		return m.CreateBusinessFn(ctx, business)
	}
	return 1, nil
}

func (m *BusinessRepositoryMock) UpdateBusiness(ctx context.Context, id uint, business domain.Business) (string, error) {
	if m.UpdateBusinessFn != nil {
		return m.UpdateBusinessFn(ctx, id, business)
	}
	return "", nil
}

func (m *BusinessRepositoryMock) DeleteBusiness(ctx context.Context, id uint) (string, error) {
	if m.DeleteBusinessFn != nil {
		return m.DeleteBusinessFn(ctx, id)
	}
	return "", nil
}

func (m *BusinessRepositoryMock) GetBusinessTypes(ctx context.Context) ([]domain.BusinessType, error) {
	if m.GetBusinessTypesFn != nil {
		return m.GetBusinessTypesFn(ctx)
	}
	return nil, nil
}

func (m *BusinessRepositoryMock) GetBusinessTypeByID(ctx context.Context, id uint) (*domain.BusinessType, error) {
	if m.GetBusinessTypeByIDFn != nil {
		return m.GetBusinessTypeByIDFn(ctx, id)
	}
	return &domain.BusinessType{ID: id}, nil
}

func (m *BusinessRepositoryMock) GetBusinessTypeByCode(ctx context.Context, code string) (*domain.BusinessType, error) {
	if m.GetBusinessTypeByCodeFn != nil {
		return m.GetBusinessTypeByCodeFn(ctx, code)
	}
	return nil, nil
}

func (m *BusinessRepositoryMock) GetBusinessTypeByName(ctx context.Context, name string) (*domain.BusinessType, error) {
	if m.GetBusinessTypeByNameFn != nil {
		return m.GetBusinessTypeByNameFn(ctx, name)
	}
	return nil, nil
}

func (m *BusinessRepositoryMock) CreateBusinessType(ctx context.Context, businessType domain.BusinessType) (string, error) {
	if m.CreateBusinessTypeFn != nil {
		return m.CreateBusinessTypeFn(ctx, businessType)
	}
	return "", nil
}

func (m *BusinessRepositoryMock) UpdateBusinessType(ctx context.Context, id uint, businessType domain.BusinessType) (string, error) {
	if m.UpdateBusinessTypeFn != nil {
		return m.UpdateBusinessTypeFn(ctx, id, businessType)
	}
	return "", nil
}

func (m *BusinessRepositoryMock) DeleteBusinessType(ctx context.Context, id uint) (string, error) {
	if m.DeleteBusinessTypeFn != nil {
		return m.DeleteBusinessTypeFn(ctx, id)
	}
	return "", nil
}

func (m *BusinessRepositoryMock) GetBusinessTypeResourcesPermitted(ctx context.Context, businessTypeID uint) ([]domain.BusinessTypeResourcePermitted, error) {
	if m.GetBusinessTypeResourcesPermittedFn != nil {
		return m.GetBusinessTypeResourcesPermittedFn(ctx, businessTypeID)
	}
	return nil, nil
}

func (m *BusinessRepositoryMock) GetResourceByID(ctx context.Context, resourceID uint) (*domain.Resource, error) {
	if m.GetResourceByIDFn != nil {
		return m.GetResourceByIDFn(ctx, resourceID)
	}
	return &domain.Resource{ID: resourceID}, nil
}

func (m *BusinessRepositoryMock) GetBusinessesWithConfiguredResourcesPaginated(ctx context.Context, page, perPage int, businessID *uint, businessTypeID *uint) ([]domain.BusinessWithConfiguredResourcesResponse, int64, error) {
	if m.GetBusinessesWithConfiguredResourcesPaginatedFn != nil {
		return m.GetBusinessesWithConfiguredResourcesPaginatedFn(ctx, page, perPage, businessID, businessTypeID)
	}
	return nil, 0, nil
}

func (m *BusinessRepositoryMock) GetBusinessByIDWithConfiguredResources(ctx context.Context, businessID uint) (*domain.BusinessWithConfiguredResourcesResponse, error) {
	if m.GetBusinessByIDWithConfiguredResourcesFn != nil {
		return m.GetBusinessByIDWithConfiguredResourcesFn(ctx, businessID)
	}
	return &domain.BusinessWithConfiguredResourcesResponse{ID: businessID}, nil
}

func (m *BusinessRepositoryMock) ToggleBusinessResourceActive(ctx context.Context, businessID uint, resourceID uint, active bool) error {
	if m.ToggleBusinessResourceActiveFn != nil {
		return m.ToggleBusinessResourceActiveFn(ctx, businessID, resourceID, active)
	}
	return nil
}

func (m *BusinessRepositoryMock) ToggleBusinessActive(ctx context.Context, businessID uint, active bool) error {
	if m.ToggleBusinessActiveFn != nil {
		return m.ToggleBusinessActiveFn(ctx, businessID, active)
	}
	return nil
}

func (m *BusinessRepositoryMock) CreatePlatformIntegration(ctx context.Context, businessID uint) error {
	if m.CreatePlatformIntegrationFn != nil {
		return m.CreatePlatformIntegrationFn(ctx, businessID)
	}
	return nil
}

func (m *BusinessRepositoryMock) GetExistingOrderPrefixes(ctx context.Context) ([]string, error) {
	if m.GetExistingOrderPrefixesFn != nil {
		return m.GetExistingOrderPrefixesFn(ctx)
	}
	return nil, nil
}

func (m *BusinessRepositoryMock) GetFiscalProfile(ctx context.Context, businessID uint) (*domain.FiscalProfile, error) {
	if m.GetFiscalProfileFn != nil {
		return m.GetFiscalProfileFn(ctx, businessID)
	}
	return nil, nil
}

func (m *BusinessRepositoryMock) UpsertFiscalProfile(ctx context.Context, profile *domain.FiscalProfile) error {
	if m.UpsertFiscalProfileFn != nil {
		return m.UpsertFiscalProfileFn(ctx, profile)
	}
	return nil
}

type S3ServiceMock struct {
	UploadImageFn   func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	GetImageURLFn   func(filename string) string
	DeleteImageFn   func(ctx context.Context, filename string) error
	ImageExistsFn   func(ctx context.Context, filename string) (bool, error)
	UploadImageCall []string
	DeleteImageCall []string
}

var _ domain.IS3Service = (*S3ServiceMock)(nil)

func (m *S3ServiceMock) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	m.UploadImageCall = append(m.UploadImageCall, folder)
	if m.UploadImageFn != nil {
		return m.UploadImageFn(ctx, file, folder)
	}
	return "", nil
}

func (m *S3ServiceMock) GetImageURL(filename string) string {
	if m.GetImageURLFn != nil {
		return m.GetImageURLFn(filename)
	}
	return ""
}

func (m *S3ServiceMock) DeleteImage(ctx context.Context, filename string) error {
	m.DeleteImageCall = append(m.DeleteImageCall, filename)
	if m.DeleteImageFn != nil {
		return m.DeleteImageFn(ctx, filename)
	}
	return nil
}

func (m *S3ServiceMock) ImageExists(ctx context.Context, filename string) (bool, error) {
	if m.ImageExistsFn != nil {
		return m.ImageExistsFn(ctx, filename)
	}
	return false, nil
}

type ConfigMock struct {
	Values map[string]string
}

var _ env.IConfig = (*ConfigMock)(nil)

func NewConfigMock(values map[string]string) *ConfigMock {
	if values == nil {
		values = map[string]string{}
	}
	return &ConfigMock{Values: values}
}

func (m *ConfigMock) Get(key string) string {
	return m.Values[key]
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger {
	return &SilentLogger{}
}

func (l *SilentLogger) nop() zerolog.Logger {
	return zerolog.Nop()
}

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Info()
}

func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Error()
}

func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Warn()
}

func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Debug()
}

func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Fatal()
}

func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Panic()
}

func (l *SilentLogger) With() zerolog.Context {
	n := l.nop()
	return n.With()
}

func (l *SilentLogger) WithService(service string) log.ILogger {
	return l
}

func (l *SilentLogger) WithModule(module string) log.ILogger {
	return l
}

func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger {
	return l
}
