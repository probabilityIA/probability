package queue

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type noopPublisher struct {
	logger log.ILogger
}

func NewNoOpPublisher(logger log.ILogger) domain.SnapshotPublisher {
	return &noopPublisher{logger: logger.WithModule("tiktok.publisher")}
}

func (p *noopPublisher) PublishSnapshot(ctx context.Context, snapshot *domain.ShopSnapshot) error {
	p.logger.Warn(ctx).
		Uint("integration_id", snapshot.IntegrationID).
		Msg("RabbitMQ no disponible, snapshot de TikTok descartado")
	return nil
}
