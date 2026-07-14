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

	CreateBusinessSubscription(ctx context.Context, subscription *entities.BusinessSubscription) error
	GetLatestByBusinessID(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error)
	ListByBusinessID(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error)
	UpdateBusinessCurrentSubscriptionType(ctx context.Context, businessID uint, subscriptionTypeID uint, status string, endDate time.Time) error
	UpdateBusinessSubscriptionStatus(ctx context.Context, businessID uint, status string, endDate *time.Time) error
	GetBusinessCurrentSubscriptionTypeID(ctx context.Context, businessID uint) (*uint, error)
	ListBusinessesExpiringBetween(ctx context.Context, from, to time.Time) ([]uint, error)
	ListBusinessesJustExpired(ctx context.Context, before time.Time) ([]uint, error)

	CreateOverride(ctx context.Context, override *entities.BusinessModuleOverride) error
	DeleteOverride(ctx context.Context, businessID uint, moduleCode string) error
	ListOverridesByBusiness(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error)

	FindSuperAdminUserID(ctx context.Context) (uint, error)
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
