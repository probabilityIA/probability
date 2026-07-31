package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ITikTokUseCase interface {
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
	SyncShopInfo(ctx context.Context, integrationID string, businessID uint) (*domain.ShopSnapshot, error)
	SaveSnapshot(ctx context.Context, snapshot *domain.ShopSnapshot) error
	ListSnapshots(ctx context.Context, businessID, integrationID uint, page, pageSize int) ([]domain.ShopSnapshot, int64, error)
}

type tiktokUseCase struct {
	client    domain.ITikTokClient
	service   domain.IIntegrationService
	publisher domain.SnapshotPublisher
	repo      domain.ISnapshotRepository
	logger    log.ILogger
}

func New(
	client domain.ITikTokClient,
	service domain.IIntegrationService,
	publisher domain.SnapshotPublisher,
	repo domain.ISnapshotRepository,
	logger log.ILogger,
) ITikTokUseCase {
	return &tiktokUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		repo:      repo,
		logger:    logger.WithModule("tiktok"),
	}
}
