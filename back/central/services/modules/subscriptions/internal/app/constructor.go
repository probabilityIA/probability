package app

import (
	"context"

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

	PurchaseSubscription(ctx context.Context, dto dtos.PurchaseSubscriptionDTO) (*entities.BusinessSubscription, error)
	RegisterPayment(ctx context.Context, dto dtos.RegisterPaymentDTO) (*entities.BusinessSubscription, error)
	DisableSubscription(ctx context.Context, businessID uint) error
	GetBusinessSubscription(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error)

	GrantOverride(ctx context.Context, dto dtos.GrantOverrideDTO) error
	RevokeOverride(ctx context.Context, businessID uint, moduleCode string) error
	ListOverrides(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error)

	HasModuleAccess(ctx context.Context, businessID uint, moduleCode string) (bool, error)
	GetModuleCodes() []string

	CheckExpiringSubscriptions(ctx context.Context) error
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
