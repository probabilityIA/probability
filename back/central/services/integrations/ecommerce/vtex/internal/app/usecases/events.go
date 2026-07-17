package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const syncProgressBatch = 25

const maxReportedFailedSKUs = 50

type failedSKUs struct {
	skus  []string
	total int
}

func (f *failedSKUs) add(sku string) {
	f.total++
	if sku == "" || len(f.skus) >= maxReportedFailedSKUs {
		return
	}
	f.skus = append(f.skus, sku)
}

func (f *failedSKUs) count() int {
	return f.total
}

func (f *failedSKUs) list() []string {
	if f.skus == nil {
		return []string{}
	}
	return f.skus
}

func (f *failedSKUs) truncated() int {
	hidden := f.total - len(f.skus)
	if hidden < 0 {
		return 0
	}
	return hidden
}

func (uc *vtexUseCase) emitSyncEvent(ctx context.Context, businessID, integrationID uint, eventType string, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	_ = rabbitmq.PublishEvent(ctx, uc.rabbit, rabbitmq.EventEnvelope{
		Type:          eventType,
		Category:      "integration",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Data:          data,
	})
}

func (uc *vtexUseCase) maybeProductProgress(ctx context.Context, businessID, integrationID uint, correlationID, direction string, processed, total, created, updated, failed int) {
	if processed%syncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "vtex.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      direction,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
}
