package app

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
)

const probabilityMinSample = 20

type levelEntry struct {
	column   string
	level    string
	idGetter func(*entities.GeozoneAncestors) *uint
}

var probabilityLevels = []levelEntry{
	{"geozone_barrio_id", "barrio", func(a *entities.GeozoneAncestors) *uint { return a.BarrioID }},
	{"geozone_neighborhood_id", "neighborhood", func(a *entities.GeozoneAncestors) *uint { return a.NeighborhoodID }},
	{"geozone_admin_district_id", "admin_district", func(a *entities.GeozoneAncestors) *uint { return a.AdminDistrictID }},
	{"geozone_locality_id", "locality", func(a *entities.GeozoneAncestors) *uint { return a.LocalityID }},
	{"geozone_city_id", "city", func(a *entities.GeozoneAncestors) *uint { return a.CityID }},
	{"geozone_state_id", "state", func(a *entities.GeozoneAncestors) *uint { return a.StateID }},
	{"geozone_country_id", "country", func(a *entities.GeozoneAncestors) *uint { return a.CountryID }},
}

type ProbabilityUseCase struct {
	repo     ports.IProbabilityRepository
	resolver ports.IResolver
}

func NewProbability(repo ports.IProbabilityRepository, resolver ports.IResolver) ports.IProbabilityUseCase {
	return &ProbabilityUseCase{repo: repo, resolver: resolver}
}

func (uc *ProbabilityUseCase) GetProbability(ctx context.Context, req dtos.ProbabilityRequest) (*dtos.ProbabilityResult, error) {
	var ancestors *entities.GeozoneAncestors
	var err error

	if req.OrderID != "" {
		ancestors, err = uc.repo.AncestorsByOrderID(ctx, req.OrderID, req.BusinessID)
		if err != nil {
			return nil, err
		}
		if ancestors == nil || (ancestors.DeepestID == nil && ancestors.CityID == nil && ancestors.StateID == nil) {
			return &dtos.ProbabilityResult{Found: false}, nil
		}
	} else if req.Lat != nil && req.Lng != nil {
		ancestors, err = uc.resolver.Resolve(ctx, *req.Lat, *req.Lng, req.BusinessID)
		if err != nil {
			return nil, err
		}
		if ancestors == nil {
			return &dtos.ProbabilityResult{Found: false}, nil
		}
	} else {
		return nil, errors.New("either order_id or lat/lng required")
	}

	for _, lvl := range probabilityLevels {
		gid := lvl.idGetter(ancestors)
		if gid == nil || *gid == 0 {
			continue
		}
		agg, err := uc.repo.AggregateAtLevel(ctx, req.BusinessID, lvl.column, *gid, req.Carrier)
		if err != nil {
			return nil, err
		}
		if agg.Total < probabilityMinSample {
			continue
		}
		rate := float64(agg.Delivered) / float64(agg.Total)
		name, _, _ := uc.repo.GeozoneNameAndType(ctx, *gid)
		return &dtos.ProbabilityResult{
			Found:        true,
			DeliveryRate: &rate,
			Level:        lvl.level,
			Carrier:      req.Carrier,
			Stats: &dtos.ProbabilityLevelStats{
				GeozoneID:   *gid,
				GeozoneType: lvl.level,
				GeozoneName: name,
				Total:       agg.Total,
				Delivered:   agg.Delivered,
				Cancelled:   agg.Cancelled,
				Returned:    agg.Returned,
				InTransit:   agg.InTransit,
			},
		}, nil
	}
	return &dtos.ProbabilityResult{Found: false}, nil
}
