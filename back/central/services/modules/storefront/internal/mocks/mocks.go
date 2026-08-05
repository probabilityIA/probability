package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type StaffCall struct {
	UserID     uint
	BusinessID uint
	RoleID     uint
}

type LinkCall struct {
	ClientID uint
	UserID   uint
}

type RepositoryMock struct {
	ListActiveProductsFn func(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error)
	GetProductByIDFn     func(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error)

	ListOrdersByUserIDFn    func(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error)
	GetOrderByIDAndUserIDFn func(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error)

	GetClientByUserIDFn func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error)

	GetBusinessByCodeFn           func(ctx context.Context, code string) (*entities.StorefrontBusiness, error)
	CreateUserFn                  func(ctx context.Context, user *entities.NewUser) (uint, error)
	CreateBusinessStaffFn         func(ctx context.Context, userID, businessID, roleID uint) error
	CreateClientFn                func(ctx context.Context, client *entities.StorefrontClient) error
	GetClienteFinalRoleIDFn       func(ctx context.Context) (uint, error)
	UserExistsByEmailFn           func(ctx context.Context, email string) (bool, error)
	GetUserAuthByEmailFn          func(ctx context.Context, email string) (uint, string, error)
	VerifyPasswordFn              func(hash, password string) bool
	StaffExistsFn                 func(ctx context.Context, userID, businessID uint) (bool, error)
	GetClientByBusinessAndEmailFn func(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error)
	LinkClientUserFn              func(ctx context.Context, clientID, userID uint) error

	GetRoleLevelByUserAndBusinessFn func(ctx context.Context, userID, businessID uint) (int, error)
	GetPlatformIntegrationIDFn      func(ctx context.Context, businessID uint) (uint, error)
	IsIntegrationActiveOrMissingFn  func(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error)

	CreatedUsers   []*entities.NewUser
	CreatedClients []*entities.StorefrontClient
	StaffCalls     []StaffCall
	LinkCalls      []LinkCall
	GateChecks     []uint
	ProductLookups []string
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) ListActiveProducts(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error) {
	if m.ListActiveProductsFn != nil {
		return m.ListActiveProductsFn(ctx, businessID, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetProductByID(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error) {
	m.ProductLookups = append(m.ProductLookups, productID)
	if m.GetProductByIDFn != nil {
		return m.GetProductByIDFn(ctx, businessID, productID)
	}
	return &entities.StorefrontProduct{ID: productID, Price: 1000, Currency: "COP"}, nil
}

func (m *RepositoryMock) ListOrdersByUserID(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error) {
	if m.ListOrdersByUserIDFn != nil {
		return m.ListOrdersByUserIDFn(ctx, businessID, userID, page, pageSize)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetOrderByIDAndUserID(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error) {
	if m.GetOrderByIDAndUserIDFn != nil {
		return m.GetOrderByIDAndUserIDFn(ctx, orderID, businessID, userID)
	}
	return &entities.StorefrontOrder{}, nil
}

func (m *RepositoryMock) GetClientByUserID(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
	if m.GetClientByUserIDFn != nil {
		return m.GetClientByUserIDFn(ctx, businessID, userID)
	}
	return &entities.StorefrontClient{ID: 1, BusinessID: businessID, UserID: &userID, Name: "Ana"}, nil
}

func (m *RepositoryMock) GetBusinessByCode(ctx context.Context, code string) (*entities.StorefrontBusiness, error) {
	if m.GetBusinessByCodeFn != nil {
		return m.GetBusinessByCodeFn(ctx, code)
	}
	return &entities.StorefrontBusiness{ID: 26, Code: code}, nil
}

func (m *RepositoryMock) CreateUser(ctx context.Context, user *entities.NewUser) (uint, error) {
	m.CreatedUsers = append(m.CreatedUsers, user)
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, user)
	}
	return 42, nil
}

func (m *RepositoryMock) CreateBusinessStaff(ctx context.Context, userID, businessID, roleID uint) error {
	m.StaffCalls = append(m.StaffCalls, StaffCall{userID, businessID, roleID})
	if m.CreateBusinessStaffFn != nil {
		return m.CreateBusinessStaffFn(ctx, userID, businessID, roleID)
	}
	return nil
}

func (m *RepositoryMock) CreateClient(ctx context.Context, client *entities.StorefrontClient) error {
	m.CreatedClients = append(m.CreatedClients, client)
	if m.CreateClientFn != nil {
		return m.CreateClientFn(ctx, client)
	}
	return nil
}

func (m *RepositoryMock) GetClienteFinalRoleID(ctx context.Context) (uint, error) {
	if m.GetClienteFinalRoleIDFn != nil {
		return m.GetClienteFinalRoleIDFn(ctx)
	}
	return 9, nil
}

func (m *RepositoryMock) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.UserExistsByEmailFn != nil {
		return m.UserExistsByEmailFn(ctx, email)
	}
	return false, nil
}

func (m *RepositoryMock) GetUserAuthByEmail(ctx context.Context, email string) (uint, string, error) {
	if m.GetUserAuthByEmailFn != nil {
		return m.GetUserAuthByEmailFn(ctx, email)
	}
	return 0, "", nil
}

func (m *RepositoryMock) VerifyPassword(hash, password string) bool {
	if m.VerifyPasswordFn != nil {
		return m.VerifyPasswordFn(hash, password)
	}
	return true
}

func (m *RepositoryMock) StaffExists(ctx context.Context, userID, businessID uint) (bool, error) {
	if m.StaffExistsFn != nil {
		return m.StaffExistsFn(ctx, userID, businessID)
	}
	return false, nil
}

func (m *RepositoryMock) GetClientByBusinessAndEmail(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error) {
	if m.GetClientByBusinessAndEmailFn != nil {
		return m.GetClientByBusinessAndEmailFn(ctx, businessID, email)
	}
	return nil, nil
}

func (m *RepositoryMock) LinkClientUser(ctx context.Context, clientID, userID uint) error {
	m.LinkCalls = append(m.LinkCalls, LinkCall{clientID, userID})
	if m.LinkClientUserFn != nil {
		return m.LinkClientUserFn(ctx, clientID, userID)
	}
	return nil
}

func (m *RepositoryMock) GetRoleLevelByUserAndBusiness(ctx context.Context, userID, businessID uint) (int, error) {
	if m.GetRoleLevelByUserAndBusinessFn != nil {
		return m.GetRoleLevelByUserAndBusinessFn(ctx, userID, businessID)
	}
	return 0, nil
}

func (m *RepositoryMock) GetPlatformIntegrationID(ctx context.Context, businessID uint) (uint, error) {
	if m.GetPlatformIntegrationIDFn != nil {
		return m.GetPlatformIntegrationIDFn(ctx, businessID)
	}
	return 5, nil
}

func (m *RepositoryMock) IsIntegrationActiveOrMissing(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error) {
	m.GateChecks = append(m.GateChecks, integrationTypeID)
	if m.IsIntegrationActiveOrMissingFn != nil {
		return m.IsIntegrationActiveOrMissingFn(ctx, businessID, integrationTypeID)
	}
	return true, nil
}

type PublisherMock struct {
	PublishOrderFn func(ctx context.Context, order []byte) error
	Published      [][]byte
}

var _ ports.IStorefrontPublisher = (*PublisherMock)(nil)

func (m *PublisherMock) PublishOrder(ctx context.Context, order []byte) error {
	m.Published = append(m.Published, order)
	if m.PublishOrderFn != nil {
		return m.PublishOrderFn(ctx, order)
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
