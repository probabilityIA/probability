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

type LevelAggregate struct {
	Total      int64
	Delivered  int64
	Cancelled  int64
	Returned   int64
	InTransit  int64
}

type IProbabilityRepository interface {
	AncestorsByOrderID(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error)
	AggregateAtLevel(ctx context.Context, businessID uint, levelColumn string, geozoneID uint, carrier string) (LevelAggregate, error)
	GeozoneNameAndType(ctx context.Context, geozoneID uint) (string, string, error)
	CarriersForBusiness(ctx context.Context, businessID uint) ([]string, error)
	GlobalCarrierStats(ctx context.Context, carrier string) (delivered int64, total int64, err error)
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
