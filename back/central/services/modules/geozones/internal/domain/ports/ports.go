package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, dto dtos.CreateGeozoneDTO) (*entities.Geozone, error)
	GetByID(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error)
	GetByCode(ctx context.Context, businessID uint, geozoneType, code string) (*entities.Geozone, error)
	List(ctx context.Context, params dtos.ListGeozonesParams) ([]entities.Geozone, int64, error)
	LookupContaining(ctx context.Context, params dtos.LookupParams) ([]entities.Geozone, error)
	Delete(ctx context.Context, id uint) error
	GetForDisplay(ctx context.Context, params dtos.DisplayParams) ([]dtos.DisplayFeature, error)
}

type IDisplayCache interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte) error
	FlushAll(ctx context.Context) error
}
