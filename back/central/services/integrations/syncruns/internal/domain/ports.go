package domain

import "context"

type IRepository interface {
	IFindingsRepository
	IDataDiffRepository

	Save(ctx context.Context, run *SyncRun) error
	ListLastByBusiness(ctx context.Context, businessID uint) ([]SyncRun, error)
	ListDetail(ctx context.Context, query DetailQuery) (*DetailPage, error)
	IntegrationBelongsToBusiness(ctx context.Context, integrationID, businessID uint) (bool, error)
}

type IUseCase interface {
	Record(ctx context.Context, run *SyncRun) error
	ListLast(ctx context.Context, businessID uint) ([]SyncRun, error)
	ListDetail(ctx context.Context, query DetailQuery) (*DetailPage, error)
	Findings(ctx context.Context, businessID uint) (*FindingsReport, error)
	FindingItems(ctx context.Context, query FindingItemsQuery) (*FindingItemsPage, error)
	MatchMatrix(ctx context.Context, query MatrixQuery) (*MatrixPage, error)
	DataSummary(ctx context.Context, businessID uint) (*DataSummary, error)
	DataPreview(ctx context.Context, query DataPreviewQuery) (*DataPreview, error)
}
