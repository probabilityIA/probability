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
	ResolveAncestors(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error)
}

type IResolver interface {
	Resolve(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error)
}

type IProbabilityRepository interface {
	AncestorsByOrderID(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error)
	ProbabilityByOrder(ctx context.Context, ancestors *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error)
	ProbabilityForCarrier(ctx context.Context, ancestors *entities.GeozoneAncestors, carrierKey string) (*dtos.ProbabilityResult, error)
	RefreshAggregates(ctx context.Context) error
}

type IProbabilityUseCase interface {
	GetProbability(ctx context.Context, req dtos.ProbabilityRequest) (*dtos.ProbabilityResult, error)
	GetOrderZone(ctx context.Context, orderID string, businessID uint) (*entities.Geozone, error)
	GetProbabilityByCarrier(ctx context.Context, orderID string, businessID uint) ([]dtos.ProbabilityResult, error)
}

type IDisplayCache interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte) error
	FlushAll(ctx context.Context) error
}

type IProbabilityCache interface {
	GetByOrder(ctx context.Context, businessID uint, orderID string) ([]dtos.ProbabilityResult, bool)
	SetByOrder(ctx context.Context, businessID uint, orderID string, results []dtos.ProbabilityResult) error
	InvalidateOrder(ctx context.Context, businessID uint, orderID string) error
}
