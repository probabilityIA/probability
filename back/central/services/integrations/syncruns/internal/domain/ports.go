package domain

import "context"

type IRepository interface {
	Save(ctx context.Context, run *SyncRun) error
	ListLastByBusiness(ctx context.Context, businessID uint) ([]SyncRun, error)
	IntegrationBelongsToBusiness(ctx context.Context, integrationID, businessID uint) (bool, error)
}

type IUseCase interface {
	Record(ctx context.Context, run *SyncRun) error
	ListLast(ctx context.Context, businessID uint) ([]SyncRun, error)
}
