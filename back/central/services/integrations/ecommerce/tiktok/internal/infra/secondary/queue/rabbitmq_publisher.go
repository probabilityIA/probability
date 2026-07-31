package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type SnapshotMessage struct {
	BusinessID           uint   `json:"business_id"`
	IntegrationID        uint   `json:"integration_id"`
	ShopID               int64  `json:"shop_id"`
	ShopName             string `json:"shop_name"`
	ShopRating           string `json:"shop_rating"`
	ShopReviewCount      int64  `json:"shop_review_count"`
	ItemSoldCount        int64  `json:"item_sold_count"`
	ShopPerformanceValue int64  `json:"shop_performance_value"`
}

type rabbitPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

func New(queue rabbitmq.IQueue, logger log.ILogger) domain.SnapshotPublisher {
	return &rabbitPublisher{
		queue:  queue,
		logger: logger.WithModule("tiktok.publisher"),
	}
}

func (p *rabbitPublisher) PublishSnapshot(ctx context.Context, snapshot *domain.ShopSnapshot) error {
	if err := p.queue.DeclareQueue(rabbitmq.QueueTiktokShopSnapshots, true); err != nil {
		return fmt.Errorf("declaring tiktok snapshots queue: %w", err)
	}

	msg := SnapshotMessage{
		BusinessID:           snapshot.BusinessID,
		IntegrationID:        snapshot.IntegrationID,
		ShopID:               snapshot.ShopID,
		ShopName:             snapshot.ShopName,
		ShopRating:           snapshot.ShopRating,
		ShopReviewCount:      snapshot.ShopReviewCount,
		ItemSoldCount:        snapshot.ItemSoldCount,
		ShopPerformanceValue: snapshot.ShopPerformanceValue,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling tiktok snapshot: %w", err)
	}

	return p.queue.Publish(ctx, rabbitmq.QueueTiktokShopSnapshots, data)
}
