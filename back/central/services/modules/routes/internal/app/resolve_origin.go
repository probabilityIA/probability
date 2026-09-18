package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
)

func (uc *UseCase) resolveOriginWarehouse(ctx context.Context, route *entities.Route) (*entities.OriginWarehouse, error) {
	if route.OriginWarehouseID != nil && *route.OriginWarehouseID > 0 {
		wh, err := uc.repo.GetWarehouseOrigin(ctx, route.BusinessID, *route.OriginWarehouseID)
		if err != nil || wh != nil {
			return wh, err
		}
	}

	orderIDs := make([]string, 0, len(route.Stops))
	for _, s := range route.Stops {
		if s.OrderID != nil && *s.OrderID != "" {
			orderIDs = append(orderIDs, *s.OrderID)
		}
	}
	ids, err := uc.repo.GetOrdersWarehouseIDs(ctx, route.BusinessID, orderIDs)
	if err != nil {
		return nil, err
	}
	if len(ids) == 1 {
		wh, err := uc.repo.GetWarehouseOrigin(ctx, route.BusinessID, ids[0])
		if err != nil || wh != nil {
			return wh, err
		}
	}

	return uc.repo.GetDefaultWarehouseOrigin(ctx, route.BusinessID)
}

func (uc *UseCase) ensureRouteOrigin(ctx context.Context, route *entities.Route) error {
	if route.OriginLat != nil && route.OriginLng != nil {
		return nil
	}

	wh, err := uc.resolveOriginWarehouse(ctx, route)
	if err != nil {
		return err
	}
	if wh == nil {
		return domainerrors.ErrOriginMissing
	}

	applyOriginWarehouse(route, wh)

	if wh.Lat == nil || wh.Lng == nil {
		return domainerrors.ErrOriginWithoutLocation
	}
	return nil
}

func applyOriginWarehouse(route *entities.Route, wh *entities.OriginWarehouse) {
	id := wh.ID
	route.OriginWarehouseID = &id
	if route.OriginAddress == "" {
		route.OriginAddress = joinAddress(wh.Address, wh.City)
	}
	route.OriginLat = wh.Lat
	route.OriginLng = wh.Lng
}

func joinAddress(address, city string) string {
	switch {
	case address == "":
		return city
	case city == "":
		return address
	default:
		return address + ", " + city
	}
}

func sameOrigin(latA, lngA, latB, lngB *float64) bool {
	return sameCoord(latA, latB) && sameCoord(lngA, lngB)
}

func sameCoord(a, b *float64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}
