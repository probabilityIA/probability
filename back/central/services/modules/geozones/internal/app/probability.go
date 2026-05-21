package app

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/repository"
)

type ProbabilityUseCase struct {
	repo     ports.IProbabilityRepository
	mainRepo ports.IRepository
	resolver ports.IResolver
	cache    ports.IProbabilityCache
}

func NewProbability(repo ports.IProbabilityRepository, mainRepo ports.IRepository, resolver ports.IResolver, cache ports.IProbabilityCache) ports.IProbabilityUseCase {
	return &ProbabilityUseCase{repo: repo, mainRepo: mainRepo, resolver: resolver, cache: cache}
}

func (uc *ProbabilityUseCase) GetProbabilityByCarrier(ctx context.Context, orderID string, businessID uint) ([]dtos.ProbabilityResult, error) {
	if uc.cache != nil {
		if cached, ok := uc.cache.GetByOrder(ctx, businessID, orderID); ok {
			return cached, nil
		}
	}
	ancestors, err := uc.repo.AncestorsByOrderID(ctx, orderID, businessID)
	if err != nil {
		return nil, err
	}
	if ancestors == nil {
		return []dtos.ProbabilityResult{}, nil
	}
	results, err := uc.repo.ProbabilityByOrder(ctx, ancestors)
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []dtos.ProbabilityResult{}
	}
	applyBaselineCascade(results)
	if uc.cache != nil {
		_ = uc.cache.SetByOrder(ctx, businessID, orderID, results)
	}
	return results, nil
}

func applyBaselineCascade(results []dtos.ProbabilityResult) {
	for i := range results {
		r := &results[i]
		if r.DeliveryRate != nil && *r.DeliveryRate > 0.01 {
			continue
		}
		if r.GlobalRate != nil && r.GlobalTotal > 0 && *r.GlobalRate > 0.01 {
			rate := *r.GlobalRate
			r.DeliveryRate = &rate
			r.IsEstimated = true
			r.EstimateSource = "global_carrier"
			continue
		}
		rate := baselineForCarrier(r.Carrier)
		r.DeliveryRate = &rate
		r.IsEstimated = true
		r.EstimateSource = "carrier_baseline"
	}
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
	} else if req.Lat != nil && req.Lng != nil {
		ancestors, err = uc.resolver.Resolve(ctx, *req.Lat, *req.Lng, req.BusinessID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("either order_id or lat/lng required")
	}

	if ancestors == nil {
		return &dtos.ProbabilityResult{Found: false, Carrier: req.Carrier}, nil
	}

	if req.Carrier == "" {
		results, err := uc.repo.ProbabilityByOrder(ctx, ancestors)
		if err != nil {
			return nil, err
		}
		for i := range results {
			if results[i].Found {
				return &results[i], nil
			}
		}
		if len(results) > 0 {
			return &results[0], nil
		}
		return &dtos.ProbabilityResult{Found: false}, nil
	}

	res, err := uc.repo.ProbabilityForCarrier(ctx, ancestors, repository.NormalizeCarrierKey(req.Carrier))
	if err != nil {
		return nil, err
	}
	if res == nil {
		res = &dtos.ProbabilityResult{Found: false, Carrier: req.Carrier}
	}
	if res.Carrier == "" {
		res.Carrier = req.Carrier
	}
	tmp := []dtos.ProbabilityResult{*res}
	applyBaselineCascade(tmp)
	return &tmp[0], nil
}
