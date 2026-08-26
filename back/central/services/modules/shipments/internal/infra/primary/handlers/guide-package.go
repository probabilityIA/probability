package handlers

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/shippingpkg"
)

func isDefaultPackage(pkg map[string]interface{}) bool {
	weight, _ := pkg["weight"].(float64)
	height, _ := pkg["height"].(float64)
	width, _ := pkg["width"].(float64)
	length, _ := pkg["length"].(float64)
	return weight <= shippingpkg.DefaultWeightKg &&
		height <= shippingpkg.DefaultDimCm &&
		width <= shippingpkg.DefaultDimCm &&
		length <= shippingpkg.DefaultDimCm
}

func (h *Handlers) applyPackageConfig(ctx context.Context, businessID uint, orderID string, raw map[string]interface{}) {
	if orderID == "" {
		return
	}

	pkgs, ok := raw["packages"].([]interface{})
	if !ok || len(pkgs) == 0 {
		return
	}
	pkg, ok := pkgs[0].(map[string]interface{})
	if !ok || !isDefaultPackage(pkg) {
		return
	}

	items, warehouseID, err := h.uc.Repo().GetOrderPackageItems(ctx, orderID)
	if err != nil {
		return
	}

	var whPtr *uint
	if warehouseID > 0 {
		whPtr = &warehouseID
	}

	cfg, err := h.uc.Repo().GetBusinessPackageConfig(ctx, businessID, whPtr)
	if err != nil || cfg == nil {
		return
	}

	quantity := 0
	for _, it := range items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		quantity += qty
	}

	resolved := shippingpkg.Resolve(cfg.Strategy, cfg.Boxes, shippingpkg.PackageInput{
		TotalQuantity: quantity,
		Items:         items,
	})

	pkg["weight"] = resolved.Weight
	pkg["height"] = resolved.Height
	pkg["width"] = resolved.Width
	pkg["length"] = resolved.Length
}
