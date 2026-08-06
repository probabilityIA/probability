package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/pay"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	GetBusinessBySlugFn     func(ctx context.Context, slug string) (*entities.BusinessPage, error)
	ListActiveProductsFn    func(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.PublicProduct, int64, error)
	GetProductByIDFn        func(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error)
	GetFeaturedProductsFn   func(ctx context.Context, businessID uint, limit int) ([]entities.PublicProduct, error)
	SaveContactSubmissionFn func(ctx context.Context, businessID uint, dto *dtos.ContactFormDTO) error

	IsIntegrationActiveFn    func(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error)
	GetTiendaIntegrationIDFn func(ctx context.Context, businessID uint) (uint, error)

	CreateCheckoutFn         func(ctx context.Context, checkout *entities.PublicCheckout) error
	GetCheckoutByReferenceFn func(ctx context.Context, reference string) (*entities.PublicCheckout, error)

	GetClientSessionFn      func(ctx context.Context, businessID uint, userID uint) (*entities.TiendaSession, error)
	GetUserBasicFn          func(ctx context.Context, userID uint) (*entities.TiendaSession, error)
	UserBelongsToBusinessFn func(ctx context.Context, businessID uint, userID uint) (bool, error)
	GetCustomerPricesFn     func(ctx context.Context, businessID uint, clientID uint) (map[string]float64, error)
	ListCustomerOrdersFn    func(ctx context.Context, businessID, userID, customerID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error)

	UpdateClientProfileFn   func(ctx context.Context, businessID, userID uint, name, phone, dni string) error
	UpdateUserAvatarFn      func(ctx context.Context, userID uint, avatarURL string) error
	GetUserAvatarFn         func(ctx context.Context, userID uint) (string, error)
	ListCustomerAddressesFn func(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error)
	CreateCustomerAddressFn func(ctx context.Context, businessID, customerID uint, addr *entities.CustomerAddress) error
	SetPrimaryAddressFn     func(ctx context.Context, businessID, customerID, addressID uint) error
	DeleteCustomerAddressFn func(ctx context.Context, businessID, customerID, addressID uint) error

	GetHiddenCategoriesFn  func(ctx context.Context, businessID uint) ([]string, error)
	ListPublicCategoriesFn func(ctx context.Context, businessID uint, hidden []string) ([]string, error)

	CreatedCheckout  *entities.PublicCheckout
	CreatedAddresses []*entities.CustomerAddress
	GateChecks       []uint
	HiddenPassed     [][]string
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) GetBusinessBySlug(ctx context.Context, slug string) (*entities.BusinessPage, error) {
	if m.GetBusinessBySlugFn != nil {
		return m.GetBusinessBySlugFn(ctx, slug)
	}
	return &entities.BusinessPage{ID: 26, Code: slug}, nil
}

func (m *RepositoryMock) ListActiveProducts(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.PublicProduct, int64, error) {
	if m.ListActiveProductsFn != nil {
		return m.ListActiveProductsFn(ctx, businessID, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetProductByID(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
	if m.GetProductByIDFn != nil {
		return m.GetProductByIDFn(ctx, businessID, productID)
	}
	return &entities.PublicProduct{ID: productID, Price: 1000}, nil
}

func (m *RepositoryMock) GetFeaturedProducts(ctx context.Context, businessID uint, limit int) ([]entities.PublicProduct, error) {
	if m.GetFeaturedProductsFn != nil {
		return m.GetFeaturedProductsFn(ctx, businessID, limit)
	}
	return nil, nil
}

func (m *RepositoryMock) SaveContactSubmission(ctx context.Context, businessID uint, dto *dtos.ContactFormDTO) error {
	if m.SaveContactSubmissionFn != nil {
		return m.SaveContactSubmissionFn(ctx, businessID, dto)
	}
	return nil
}

func (m *RepositoryMock) IsIntegrationActive(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error) {
	m.GateChecks = append(m.GateChecks, integrationTypeID)
	if m.IsIntegrationActiveFn != nil {
		return m.IsIntegrationActiveFn(ctx, businessID, integrationTypeID)
	}
	return true, nil
}

func (m *RepositoryMock) GetTiendaIntegrationID(ctx context.Context, businessID uint) (uint, error) {
	if m.GetTiendaIntegrationIDFn != nil {
		return m.GetTiendaIntegrationIDFn(ctx, businessID)
	}
	return 5, nil
}

func (m *RepositoryMock) CreateCheckout(ctx context.Context, checkout *entities.PublicCheckout) error {
	m.CreatedCheckout = checkout
	if m.CreateCheckoutFn != nil {
		return m.CreateCheckoutFn(ctx, checkout)
	}
	return nil
}

func (m *RepositoryMock) GetCheckoutByReference(ctx context.Context, reference string) (*entities.PublicCheckout, error) {
	if m.GetCheckoutByReferenceFn != nil {
		return m.GetCheckoutByReferenceFn(ctx, reference)
	}
	return &entities.PublicCheckout{Reference: reference, Status: entities.CheckoutStatusAgreed}, nil
}

func (m *RepositoryMock) GetClientSession(ctx context.Context, businessID uint, userID uint) (*entities.TiendaSession, error) {
	if m.GetClientSessionFn != nil {
		return m.GetClientSessionFn(ctx, businessID, userID)
	}
	return &entities.TiendaSession{Name: "Ana", CustomerID: 7}, nil
}

func (m *RepositoryMock) GetUserBasic(ctx context.Context, userID uint) (*entities.TiendaSession, error) {
	if m.GetUserBasicFn != nil {
		return m.GetUserBasicFn(ctx, userID)
	}
	return &entities.TiendaSession{Name: "Basico"}, nil
}

func (m *RepositoryMock) UserBelongsToBusiness(ctx context.Context, businessID uint, userID uint) (bool, error) {
	if m.UserBelongsToBusinessFn != nil {
		return m.UserBelongsToBusinessFn(ctx, businessID, userID)
	}
	return true, nil
}

func (m *RepositoryMock) GetCustomerPrices(ctx context.Context, businessID uint, clientID uint) (map[string]float64, error) {
	if m.GetCustomerPricesFn != nil {
		return m.GetCustomerPricesFn(ctx, businessID, clientID)
	}
	return map[string]float64{}, nil
}

func (m *RepositoryMock) ListCustomerOrders(ctx context.Context, businessID, userID, customerID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error) {
	if m.ListCustomerOrdersFn != nil {
		return m.ListCustomerOrdersFn(ctx, businessID, userID, customerID, page, pageSize)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) UpdateClientProfile(ctx context.Context, businessID, userID uint, name, phone, dni string) error {
	if m.UpdateClientProfileFn != nil {
		return m.UpdateClientProfileFn(ctx, businessID, userID, name, phone, dni)
	}
	return nil
}

func (m *RepositoryMock) UpdateUserAvatar(ctx context.Context, userID uint, avatarURL string) error {
	if m.UpdateUserAvatarFn != nil {
		return m.UpdateUserAvatarFn(ctx, userID, avatarURL)
	}
	return nil
}

func (m *RepositoryMock) GetUserAvatar(ctx context.Context, userID uint) (string, error) {
	if m.GetUserAvatarFn != nil {
		return m.GetUserAvatarFn(ctx, userID)
	}
	return "", nil
}

func (m *RepositoryMock) ListCustomerAddresses(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error) {
	if m.ListCustomerAddressesFn != nil {
		return m.ListCustomerAddressesFn(ctx, businessID, customerID)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateCustomerAddress(ctx context.Context, businessID, customerID uint, addr *entities.CustomerAddress) error {
	m.CreatedAddresses = append(m.CreatedAddresses, addr)
	if m.CreateCustomerAddressFn != nil {
		return m.CreateCustomerAddressFn(ctx, businessID, customerID, addr)
	}
	return nil
}

func (m *RepositoryMock) SetPrimaryAddress(ctx context.Context, businessID, customerID, addressID uint) error {
	if m.SetPrimaryAddressFn != nil {
		return m.SetPrimaryAddressFn(ctx, businessID, customerID, addressID)
	}
	return nil
}

func (m *RepositoryMock) DeleteCustomerAddress(ctx context.Context, businessID, customerID, addressID uint) error {
	if m.DeleteCustomerAddressFn != nil {
		return m.DeleteCustomerAddressFn(ctx, businessID, customerID, addressID)
	}
	return nil
}

func (m *RepositoryMock) GetHiddenCategories(ctx context.Context, businessID uint) ([]string, error) {
	if m.GetHiddenCategoriesFn != nil {
		return m.GetHiddenCategoriesFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListPublicCategories(ctx context.Context, businessID uint, hidden []string) ([]string, error) {
	m.HiddenPassed = append(m.HiddenPassed, hidden)
	if m.ListPublicCategoriesFn != nil {
		return m.ListPublicCategoriesFn(ctx, businessID, hidden)
	}
	return nil, nil
}

type BoldGatewayMock struct {
	BoldGenerateSignatureForReferenceFn func(ctx context.Context, businessID uint, amount float64, currency, referencePrefix string) (*pay.BoldSignature, error)
	PublishAgreedStorefrontOrderFn      func(ctx context.Context, reference string) error

	PublishedRefs []string
}

var _ ports.IBoldGateway = (*BoldGatewayMock)(nil)

func (m *BoldGatewayMock) BoldGenerateSignatureForReference(ctx context.Context, businessID uint, amount float64, currency, referencePrefix string) (*pay.BoldSignature, error) {
	if m.BoldGenerateSignatureForReferenceFn != nil {
		return m.BoldGenerateSignatureForReferenceFn(ctx, businessID, amount, currency, referencePrefix)
	}
	return &pay.BoldSignature{}, nil
}

func (m *BoldGatewayMock) PublishAgreedStorefrontOrder(ctx context.Context, reference string) error {
	m.PublishedRefs = append(m.PublishedRefs, reference)
	if m.PublishAgreedStorefrontOrderFn != nil {
		return m.PublishAgreedStorefrontOrderFn(ctx, reference)
	}
	return nil
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger { return &SilentLogger{} }

func (l *SilentLogger) nop() zerolog.Logger { return zerolog.Nop() }

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Info() }
func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Error() }
func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Warn() }
func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Debug() }
func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Fatal() }
func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Panic() }
func (l *SilentLogger) With() zerolog.Context                       { n := l.nop(); return n.With() }
func (l *SilentLogger) WithService(service string) log.ILogger      { return l }
func (l *SilentLogger) WithModule(module string) log.ILogger        { return l }
func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger  { return l }
