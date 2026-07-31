package domain

import "context"

type ITikTokClient interface {
	GetClientToken(ctx context.Context, baseURL, clientKey, clientSecret string) (string, error)
	QueryShopInfo(ctx context.Context, cred Credential, shopName string, limit int) ([]ShopInfo, error)
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
}

type SnapshotPublisher interface {
	PublishSnapshot(ctx context.Context, snapshot *ShopSnapshot) error
}

type ISnapshotRepository interface {
	Save(ctx context.Context, snapshot *ShopSnapshot) error
	ListByIntegration(ctx context.Context, businessID, integrationID uint, page, pageSize int) ([]ShopSnapshot, int64, error)
}
