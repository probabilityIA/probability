package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	CreateSubscriptionType(ctx context.Context, dto dtos.CreateSubscriptionTypeDTO) (*entities.SubscriptionType, error)
	UpdateSubscriptionType(ctx context.Context, dto dtos.UpdateSubscriptionTypeDTO) (*entities.SubscriptionType, error)
	DeleteSubscriptionType(ctx context.Context, id uint) error
	GetSubscriptionType(ctx context.Context, id uint) (*entities.SubscriptionType, error)
	ListSubscriptionTypes(ctx context.Context, activeOnly bool) ([]entities.SubscriptionType, error)

	CreateCustomPlan(ctx context.Context, dto dtos.CreateCustomPlanDTO, actorUserID uint) (*entities.SubscriptionType, error)
	ListCustomPlans(ctx context.Context, businessID *uint) ([]entities.SubscriptionType, error)

	PurchaseSubscription(ctx context.Context, dto dtos.PurchaseSubscriptionDTO) (*entities.BusinessSubscription, error)
	RegisterPayment(ctx context.Context, dto dtos.RegisterPaymentDTO, actorUserID uint) (*entities.BusinessSubscription, error)
	EditSubscriptionDates(ctx context.Context, dto dtos.EditSubscriptionDatesDTO, actorUserID uint) (*entities.BusinessSubscription, error)
	DisableSubscription(ctx context.Context, businessID uint, actorUserID uint) error
	ReactivateSubscription(ctx context.Context, businessID uint, actorUserID uint) error
	ExtendCourtesy(ctx context.Context, dto dtos.ExtendCourtesyDTO, actorUserID uint) (*entities.BusinessSubscription, error)
	RevertPayment(ctx context.Context, subscriptionID uint, actorUserID uint) (*entities.BusinessSubscription, error)
	ListPaymentHistory(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error)
	ListAuditLogs(ctx context.Context, businessID uint, limit int) ([]entities.SubscriptionAuditLog, error)
	GetBusinessSubscription(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error)
	GetBusinessSubscriptionMeta(ctx context.Context, businessID uint) (*entities.BusinessSubscriptionMeta, error)
	SetAutoPaymentEnabled(ctx context.Context, businessID uint, enabled bool, actorUserID uint) error
	GetPaymentWindow(ctx context.Context, businessID uint) (*time.Time, *time.Time, error)
	GetSubscriptionUsage(ctx context.Context, businessID uint) (*entities.SubscriptionUsage, error)
	CheckShipmentOverage(ctx context.Context, businessID uint) (blocked bool, reason string, fee float64, err error)
	AcceptShipmentOverage(ctx context.Context, businessID uint) error
	PayShipmentOverage(ctx context.Context, businessID uint, actorUserID uint) error

	GrantOverride(ctx context.Context, dto dtos.GrantOverrideDTO) error
	RevokeOverride(ctx context.Context, businessID uint, moduleCode string, actorUserID uint) error
	ListOverrides(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error)

	HasModuleAccess(ctx context.Context, businessID uint, moduleCode string) (bool, error)
	GetModuleCodes() []string
	GetModuleCatalog() []dtos.ModuleInfo
	GetAccessibleModules(ctx context.Context, businessID uint) ([]string, error)
	EcommerceChannelLimit(ctx context.Context, businessID uint) (int, error)

	CheckExpiringSubscriptions(ctx context.Context) error
	DowngradeExpiredTrials(ctx context.Context) error
	SettleFreePlanCycles(ctx context.Context) error
	AssignTrialSubscription(ctx context.Context, businessID uint)

	ListAdminBusinesses(ctx context.Context, page, pageSize int, search, statusFilter string) ([]entities.AdminBusinessRow, int64, error)
	GetAdminKPIs(ctx context.Context) (entities.AdminKPIs, error)
}

type UseCase struct {
	repo          ports.IRepository
	wallet        ports.IWalletDebiter
	announcements ports.IAnnouncementsGateway
	log           log.ILogger
	systemUserID  uint
}

func New(repo ports.IRepository, wallet ports.IWalletDebiter, announcements ports.IAnnouncementsGateway, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, wallet: wallet, announcements: announcements, log: logger}
}
