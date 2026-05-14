package app

import (
	"context"
	"errors"
	"sync"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
)

const probabilityMinSample = 5

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
	mainRepo ports.IRepository
	resolver ports.IResolver
}

func NewProbability(repo ports.IProbabilityRepository, mainRepo ports.IRepository, resolver ports.IResolver) ports.IProbabilityUseCase {
	return &ProbabilityUseCase{repo: repo, mainRepo: mainRepo, resolver: resolver}
}

func (uc *ProbabilityUseCase) GetProbabilityByCarrier(ctx context.Context, orderID string, businessID uint) ([]dtos.ProbabilityResult, error) {
	carriers, err := uc.repo.CarriersForBusiness(ctx, businessID)
	if err != nil {
		return nil, err
	}
	results := make([]dtos.ProbabilityResult, len(carriers))
	errCh := make(chan error, len(carriers))
	const workers = 8
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	for i, carrier := range carriers {
		i, carrier := i, carrier
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			req := dtos.ProbabilityRequest{BusinessID: businessID, OrderID: orderID, Carrier: carrier}
			res, err := uc.GetProbability(ctx, req)
			if err != nil {
				errCh <- err
				return
			}
			if res == nil {
				res = &dtos.ProbabilityResult{Found: false}
			}
			res.Carrier = carrier
			results[i] = *res
		}()
	}
	wg.Wait()
	close(errCh)
	for e := range errCh {
		if e != nil {
			return nil, e
		}
	}
	return results, nil
}

func (uc *ProbabilityUseCase) GetOrderZone(ctx context.Context, orderID string, businessID uint) (*entities.Geozone, error) {
	ancestors, err := uc.repo.AncestorsByOrderID(ctx, orderID, businessID)
	if err != nil {
		return nil, err
	}
	if ancestors == nil || ancestors.DeepestID == nil || *ancestors.DeepestID == 0 {
		return nil, nil
	}
	return uc.mainRepo.GetByID(ctx, *ancestors.DeepestID, true)
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

	out := &dtos.ProbabilityResult{Found: false, Carrier: req.Carrier}
	if req.Carrier != "" {
		if delivered, total, err := uc.repo.GlobalCarrierStats(ctx, req.Carrier); err == nil && total > 0 {
			r := float64(delivered) / float64(total)
			out.GlobalRate = &r
			out.GlobalTotal = total
		}
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
		out.Found = true
		out.DeliveryRate = &rate
		out.Level = lvl.level
		out.Stats = &dtos.ProbabilityLevelStats{
			GeozoneID:   *gid,
			GeozoneType: lvl.level,
			GeozoneName: name,
			Total:       agg.Total,
			Delivered:   agg.Delivered,
			Cancelled:   agg.Cancelled,
			Returned:    agg.Returned,
			InTransit:   agg.InTransit,
		}
		return out, nil
	}
	return out, nil
}
