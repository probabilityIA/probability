package app

import (
	"context"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	HandleIncoming(ctx context.Context, dto domain.IncomingMessageDTO) error
}

type useCase struct {
	aiProvider           domain.IAIProvider
	sessionCache         domain.ISessionCache
	productRepo          domain.IProductRepository
	customerRepo         domain.ICustomerRepository
	responsePublisher    domain.IAIResponsePublisher
	orderPublisher       domain.IAIOrderPublisher
	configProvider       domain.IConfigProvider
	persistencePublisher domain.IAIPersistencePublisher
	pauseChecker         domain.IAIPauseChecker
	log                  log.ILogger
}

func New(
	aiProvider domain.IAIProvider,
	sessionCache domain.ISessionCache,
	productRepo domain.IProductRepository,
	customerRepo domain.ICustomerRepository,
	responsePublisher domain.IAIResponsePublisher,
	orderPublisher domain.IAIOrderPublisher,
	configProvider domain.IConfigProvider,
	persistencePublisher domain.IAIPersistencePublisher,
	pauseChecker domain.IAIPauseChecker,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		aiProvider:           aiProvider,
		sessionCache:         sessionCache,
		productRepo:          productRepo,
		customerRepo:         customerRepo,
		responsePublisher:    responsePublisher,
		orderPublisher:       orderPublisher,
		configProvider:       configProvider,
		persistencePublisher: persistencePublisher,
		pauseChecker:         pauseChecker,
		log:                  logger,
	}
}
