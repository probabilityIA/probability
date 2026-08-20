package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

type IRepository interface {
	CreateSubscriptionType(ctx context.Context, subType *entities.SubscriptionType) error
	UpdateSubscriptionType(ctx context.Context, subType *entities.SubscriptionType) error
	DeleteSubscriptionType(ctx context.Context, id uint) error
	GetSubscriptionType(ctx context.Context, id uint) (*entities.SubscriptionType, error)
	ListSubscriptionTypes(ctx context.Context, activeOnly bool) ([]entities.SubscriptionType, error)
	ListCustomPlans(ctx context.Context, businessID *uint) ([]entities.SubscriptionType, error)

	CreateSubscriptionAndActivate(ctx context.Context, subscription *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error
	GetLatestByBusinessID(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error)
	GetSubscriptionUsage(ctx context.Context, businessID uint) (*entities.SubscriptionUsage, error)
	ListByBusinessID(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error)
	GetSubscriptionByID(ctx context.Context, id uint) (*entities.BusinessSubscription, error)
	RevertSubscriptionAndRecalculate(ctx context.Context, subscriptionID uint) (*entities.BusinessSubscription, error)
	UpdateBusinessSubscriptionStatus(ctx context.Context, businessID uint, status string, endDate *time.Time) error
	UpdateSubscriptionDates(ctx context.Context, id uint, startDate, endDate time.Time) error
	UpdateBusinessSubscriptionEndDate(ctx context.Context, businessID uint, endDate time.Time) error
	GetBusinessCurrentSubscriptionTypeID(ctx context.Context, businessID uint) (*uint, error)
	ListBusinessesExpiringBetween(ctx context.Context, from, to time.Time) ([]uint, error)
	ListBusinessesJustExpired(ctx context.Context, before time.Time) ([]uint, error)
	MarkExpiredIfStillActive(ctx context.Context, businessID uint, before time.Time) error

	CreateOverride(ctx context.Context, override *entities.BusinessModuleOverride) error
	DeleteOverride(ctx context.Context, businessID uint, moduleCode string) error
	ListOverridesByBusiness(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error)

	CreateAuditLog(ctx context.Context, log *entities.SubscriptionAuditLog) error
	ListAuditLogsByBusiness(ctx context.Context, businessID uint, limit int) ([]entities.SubscriptionAuditLog, error)

	FindSuperAdminUserID(ctx context.Context) (uint, error)
	ResolveUserLabel(ctx context.Context, userID uint) (string, error)

	ListBusinessesForAdmin(ctx context.Context, page, pageSize int, search, statusFilter string) ([]entities.AdminBusinessRow, int64, error)
	GetAdminKPIs(ctx context.Context) (entities.AdminKPIs, error)
}

type IWalletDebiter interface {
	GetBalance(ctx context.Context, businessID uint) (float64, error)
	Debit(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error
}

type IAnnouncementsGateway interface {
	CreateBusinessAlert(ctx context.Context, businessID uint, title, message string, createdByID uint, daily bool) (uint, error)
	FindActiveBusinessAlert(ctx context.Context, businessID uint, title string) (*uint, error)
	DeactivateAnnouncement(ctx context.Context, id uint) error
}
