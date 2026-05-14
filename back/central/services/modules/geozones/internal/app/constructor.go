package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
)

type IUseCase interface {
	Create(ctx context.Context, dto dtos.CreateGeozoneDTO) (*entities.Geozone, error)
	BulkImport(ctx context.Context, dto dtos.BulkImportDTO) (*dtos.BulkImportResult, error)
	Get(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error)
	List(ctx context.Context, params dtos.ListGeozonesParams) ([]entities.Geozone, int64, error)
	Lookup(ctx context.Context, params dtos.LookupParams) ([]entities.Geozone, error)
	Delete(ctx context.Context, id uint) error
	GetForDisplay(ctx context.Context, geozoneType string, zoom int, bbox *dtos.Bbox, parentID *uint) ([]byte, string, error)
	FlushDisplayCache(ctx context.Context) error
}

type UseCase struct {
	repo  ports.IRepository
	cache ports.IDisplayCache
}

func New(repo ports.IRepository, cache ports.IDisplayCache) IUseCase {
	return &UseCase{repo: repo, cache: cache}
}
