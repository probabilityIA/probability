package usecaseordermapping

import "context"

// publishSyncOrderRejected delega al adapter de infra la publicación del evento
func (uc *UseCaseOrderMapping) publishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, orderNumber, externalID, platform, reason, errMsg string) {
	if uc.integrationEventPublisher == nil {
		return
	}
	uc.integrationEventPublisher.PublishSyncOrderRejected(ctx, integrationID, businessID, orderNumber, externalID, platform, reason, errMsg)
}

// publishSyncOrderCreated delega al adapter de infra la publicación del evento
func (uc *UseCaseOrderMapping) publishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
	if uc.integrationEventPublisher == nil {
		return
	}
	uc.integrationEventPublisher.PublishSyncOrderCreated(ctx, integrationID, businessID, data)
}

// publishSyncOrderUpdated delega al adapter de infra la publicación del evento
func (uc *UseCaseOrderMapping) publishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
	if uc.integrationEventPublisher == nil {
		return
	}
	uc.integrationEventPublisher.PublishSyncOrderUpdated(ctx, integrationID, businessID, data)
}
