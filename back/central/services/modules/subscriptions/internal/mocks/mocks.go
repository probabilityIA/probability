package mocks

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	CreateSubscriptionTypeFn func(ctx context.Context, subType *entities.SubscriptionType) error
	UpdateSubscriptionTypeFn func(ctx context.Context, subType *entities.SubscriptionType) error
	DeleteSubscriptionTypeFn func(ctx context.Context, id uint) error
	GetSubscriptionTypeFn    func(ctx context.Context, id uint) (*entities.SubscriptionType, error)
	ListSubscriptionTypesFn  func(ctx context.Context, activeOnly bool) ([]entities.SubscriptionType, error)
	ListCustomPlansFn        func(ctx context.Context, businessID *uint) ([]entities.SubscriptionType, error)

	CreateSubscriptionAndActivateFn        func(ctx context.Context, subscription *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error
	GetLatestByBusinessIDFn                func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error)
	GetSubscriptionUsageFn                 func(ctx context.Context, businessID uint) (*entities.SubscriptionUsage, error)
	ListByBusinessIDFn                     func(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error)
	GetSubscriptionByIDFn                  func(ctx context.Context, id uint) (*entities.BusinessSubscription, error)
	RevertSubscriptionAndRecalculateFn     func(ctx context.Context, subscriptionID uint) (*entities.BusinessSubscription, error)
	UpdateBusinessSubscriptionStatusFn     func(ctx context.Context, businessID uint, status string, endDate *time.Time) error
	MarkExpiredIfStillActiveFn             func(ctx context.Context, businessID uint, before time.Time) error
	UpdateSubscriptionDatesFn              func(ctx context.Context, id uint, startDate, endDate time.Time) error
	UpdateBusinessSubscriptionEndDateFn    func(ctx context.Context, businessID uint, endDate time.Time) error
	GetBusinessCurrentSubscriptionTypeIDFn func(ctx context.Context, businessID uint) (*uint, error)
	ListBusinessesExpiringBetweenFn        func(ctx context.Context, from, to time.Time) ([]uint, error)
	ListBusinessesJustExpiredFn            func(ctx context.Context, before time.Time) ([]uint, error)

	CreateOverrideFn          func(ctx context.Context, override *entities.BusinessModuleOverride) error
	DeleteOverrideFn          func(ctx context.Context, businessID uint, moduleCode string) error
	ListOverridesByBusinessFn func(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error)

	CreateAuditLogFn          func(ctx context.Context, log *entities.SubscriptionAuditLog) error
	ListAuditLogsByBusinessFn func(ctx context.Context, businessID uint, limit int) ([]entities.SubscriptionAuditLog, error)

	FindSuperAdminUserIDFn func(ctx context.Context) (uint, error)
	ResolveUserLabelFn     func(ctx context.Context, userID uint) (string, error)

	ListBusinessesForAdminFn func(ctx context.Context, page, pageSize int, search, statusFilter string) ([]entities.AdminBusinessRow, int64, error)
	GetAdminKPIsFn           func(ctx context.Context) (entities.AdminKPIs, error)
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) CreateSubscriptionType(ctx context.Context, subType *entities.SubscriptionType) error {
	if m.CreateSubscriptionTypeFn != nil {
		return m.CreateSubscriptionTypeFn(ctx, subType)
	}
	return nil
}

func (m *RepositoryMock) UpdateSubscriptionType(ctx context.Context, subType *entities.SubscriptionType) error {
	if m.UpdateSubscriptionTypeFn != nil {
		return m.UpdateSubscriptionTypeFn(ctx, subType)
	}
	return nil
}

func (m *RepositoryMock) DeleteSubscriptionType(ctx context.Context, id uint) error {
	if m.DeleteSubscriptionTypeFn != nil {
		return m.DeleteSubscriptionTypeFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) GetSubscriptionType(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
	if m.GetSubscriptionTypeFn != nil {
		return m.GetSubscriptionTypeFn(ctx, id)
	}
	return &entities.SubscriptionType{ID: id, Active: true, Price: 100}, nil
}

func (m *RepositoryMock) ListSubscriptionTypes(ctx context.Context, activeOnly bool) ([]entities.SubscriptionType, error) {
	if m.ListSubscriptionTypesFn != nil {
		return m.ListSubscriptionTypesFn(ctx, activeOnly)
	}
	return nil, nil
}

func (m *RepositoryMock) ListCustomPlans(ctx context.Context, businessID *uint) ([]entities.SubscriptionType, error) {
	if m.ListCustomPlansFn != nil {
		return m.ListCustomPlansFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateSubscriptionAndActivate(ctx context.Context, subscription *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
	if m.CreateSubscriptionAndActivateFn != nil {
		return m.CreateSubscriptionAndActivateFn(ctx, subscription, subscriptionTypeID, endDate)
	}
	return nil
}

func (m *RepositoryMock) GetLatestByBusinessID(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
	if m.GetLatestByBusinessIDFn != nil {
		return m.GetLatestByBusinessIDFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) GetSubscriptionUsage(ctx context.Context, businessID uint) (*entities.SubscriptionUsage, error) {
	if m.GetSubscriptionUsageFn != nil {
		return m.GetSubscriptionUsageFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListByBusinessID(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error) {
	if m.ListByBusinessIDFn != nil {
		return m.ListByBusinessIDFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) GetSubscriptionByID(ctx context.Context, id uint) (*entities.BusinessSubscription, error) {
	if m.GetSubscriptionByIDFn != nil {
		return m.GetSubscriptionByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *RepositoryMock) RevertSubscriptionAndRecalculate(ctx context.Context, subscriptionID uint) (*entities.BusinessSubscription, error) {
	if m.RevertSubscriptionAndRecalculateFn != nil {
		return m.RevertSubscriptionAndRecalculateFn(ctx, subscriptionID)
	}
	return nil, nil
}

func (m *RepositoryMock) UpdateBusinessSubscriptionStatus(ctx context.Context, businessID uint, status string, endDate *time.Time) error {
	if m.UpdateBusinessSubscriptionStatusFn != nil {
		return m.UpdateBusinessSubscriptionStatusFn(ctx, businessID, status, endDate)
	}
	return nil
}

func (m *RepositoryMock) MarkExpiredIfStillActive(ctx context.Context, businessID uint, before time.Time) error {
	if m.MarkExpiredIfStillActiveFn != nil {
		return m.MarkExpiredIfStillActiveFn(ctx, businessID, before)
	}
	return nil
}

func (m *RepositoryMock) UpdateSubscriptionDates(ctx context.Context, id uint, startDate, endDate time.Time) error {
	if m.UpdateSubscriptionDatesFn != nil {
		return m.UpdateSubscriptionDatesFn(ctx, id, startDate, endDate)
	}
	return nil
}

func (m *RepositoryMock) UpdateBusinessSubscriptionEndDate(ctx context.Context, businessID uint, endDate time.Time) error {
	if m.UpdateBusinessSubscriptionEndDateFn != nil {
		return m.UpdateBusinessSubscriptionEndDateFn(ctx, businessID, endDate)
	}
	return nil
}

func (m *RepositoryMock) GetBusinessCurrentSubscriptionTypeID(ctx context.Context, businessID uint) (*uint, error) {
	if m.GetBusinessCurrentSubscriptionTypeIDFn != nil {
		return m.GetBusinessCurrentSubscriptionTypeIDFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListBusinessesExpiringBetween(ctx context.Context, from, to time.Time) ([]uint, error) {
	if m.ListBusinessesExpiringBetweenFn != nil {
		return m.ListBusinessesExpiringBetweenFn(ctx, from, to)
	}
	return nil, nil
}

func (m *RepositoryMock) ListBusinessesJustExpired(ctx context.Context, before time.Time) ([]uint, error) {
	if m.ListBusinessesJustExpiredFn != nil {
		return m.ListBusinessesJustExpiredFn(ctx, before)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateOverride(ctx context.Context, override *entities.BusinessModuleOverride) error {
	if m.CreateOverrideFn != nil {
		return m.CreateOverrideFn(ctx, override)
	}
	return nil
}

func (m *RepositoryMock) DeleteOverride(ctx context.Context, businessID uint, moduleCode string) error {
	if m.DeleteOverrideFn != nil {
		return m.DeleteOverrideFn(ctx, businessID, moduleCode)
	}
	return nil
}

func (m *RepositoryMock) ListOverridesByBusiness(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
	if m.ListOverridesByBusinessFn != nil {
		return m.ListOverridesByBusinessFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateAuditLog(ctx context.Context, log *entities.SubscriptionAuditLog) error {
	if m.CreateAuditLogFn != nil {
		return m.CreateAuditLogFn(ctx, log)
	}
	return nil
}

func (m *RepositoryMock) ListAuditLogsByBusiness(ctx context.Context, businessID uint, limit int) ([]entities.SubscriptionAuditLog, error) {
	if m.ListAuditLogsByBusinessFn != nil {
		return m.ListAuditLogsByBusinessFn(ctx, businessID, limit)
	}
	return nil, nil
}

func (m *RepositoryMock) FindSuperAdminUserID(ctx context.Context) (uint, error) {
	if m.FindSuperAdminUserIDFn != nil {
		return m.FindSuperAdminUserIDFn(ctx)
	}
	return 1, nil
}

func (m *RepositoryMock) ResolveUserLabel(ctx context.Context, userID uint) (string, error) {
	if m.ResolveUserLabelFn != nil {
		return m.ResolveUserLabelFn(ctx, userID)
	}
	return "Sistema", nil
}

func (m *RepositoryMock) ListBusinessesForAdmin(ctx context.Context, page, pageSize int, search, statusFilter string) ([]entities.AdminBusinessRow, int64, error) {
	if m.ListBusinessesForAdminFn != nil {
		return m.ListBusinessesForAdminFn(ctx, page, pageSize, search, statusFilter)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetAdminKPIs(ctx context.Context) (entities.AdminKPIs, error) {
	if m.GetAdminKPIsFn != nil {
		return m.GetAdminKPIsFn(ctx)
	}
	return entities.AdminKPIs{}, nil
}

type DebitCall struct {
	BusinessID uint
	Amount     float64
	Reference  string
	Concept    string
	UserID     uint
}

type WalletDebiterMock struct {
	GetBalanceFn func(ctx context.Context, businessID uint) (float64, error)
	DebitFn      func(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error
	DebitCalls   []DebitCall
}

var _ ports.IWalletDebiter = (*WalletDebiterMock)(nil)

func (m *WalletDebiterMock) GetBalance(ctx context.Context, businessID uint) (float64, error) {
	if m.GetBalanceFn != nil {
		return m.GetBalanceFn(ctx, businessID)
	}
	return 0, nil
}

func (m *WalletDebiterMock) Debit(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error {
	m.DebitCalls = append(m.DebitCalls, DebitCall{
		BusinessID: businessID, Amount: amount, Reference: reference, Concept: concept, UserID: userID,
	})
	if m.DebitFn != nil {
		return m.DebitFn(ctx, businessID, amount, reference, concept, userID)
	}
	return nil
}

type AlertCall struct {
	BusinessID  uint
	Title       string
	Message     string
	CreatedByID uint
	Daily       bool
}

type AnnouncementsGatewayMock struct {
	CreateBusinessAlertFn     func(ctx context.Context, businessID uint, title, message string, createdByID uint, daily bool) (uint, error)
	FindActiveBusinessAlertFn func(ctx context.Context, businessID uint, title string) (*uint, error)
	DeactivateAnnouncementFn  func(ctx context.Context, id uint) error
	CreatedAlerts             []AlertCall
	DeactivatedIDs            []uint
}

var _ ports.IAnnouncementsGateway = (*AnnouncementsGatewayMock)(nil)

func (m *AnnouncementsGatewayMock) CreateBusinessAlert(ctx context.Context, businessID uint, title, message string, createdByID uint, daily bool) (uint, error) {
	m.CreatedAlerts = append(m.CreatedAlerts, AlertCall{
		BusinessID: businessID, Title: title, Message: message, CreatedByID: createdByID, Daily: daily,
	})
	if m.CreateBusinessAlertFn != nil {
		return m.CreateBusinessAlertFn(ctx, businessID, title, message, createdByID, daily)
	}
	return 1, nil
}

func (m *AnnouncementsGatewayMock) FindActiveBusinessAlert(ctx context.Context, businessID uint, title string) (*uint, error) {
	if m.FindActiveBusinessAlertFn != nil {
		return m.FindActiveBusinessAlertFn(ctx, businessID, title)
	}
	return nil, nil
}

func (m *AnnouncementsGatewayMock) DeactivateAnnouncement(ctx context.Context, id uint) error {
	m.DeactivatedIDs = append(m.DeactivatedIDs, id)
	if m.DeactivateAnnouncementFn != nil {
		return m.DeactivateAnnouncementFn(ctx, id)
	}
	return nil
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
