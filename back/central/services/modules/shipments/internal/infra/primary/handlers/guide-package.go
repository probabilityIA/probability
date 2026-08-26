package handlers

import (
	"context"
	"strconv"

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
	autoPackage, _ := raw["auto_package"].(bool)
	delete(raw, "auto_package")

	pkgs, ok := raw["packages"].([]interface{})
	if !ok || len(pkgs) == 0 {
		return
	}
	pkg, ok := pkgs[0].(map[string]interface{})
	if !ok || (!autoPackage && !isDefaultPackage(pkg)) {
		return
	}
	cartWeight := 0.0
	if autoPackage {
		cartWeight, _ = pkg["weight"].(float64)
	}

	var items []shippingpkg.PackageItem
	whPtr := requestWarehouseID(raw)
	if orderID != "" {
		orderItems, warehouseID, err := h.uc.Repo().GetOrderPackageItems(ctx, orderID)
		if err != nil {
			return
		}
		items = orderItems
		if warehouseID > 0 {
			whPtr = &warehouseID
		}
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
		CartWeightKg:  cartWeight,
		Items:         items,
	})

	setPackageDims(pkg, resolved)
}

func parseQuoteItems(raw map[string]interface{}) []shippingpkg.PackageItem {
	rawItems, ok := raw["items"].([]interface{})
	if !ok || len(rawItems) == 0 {
		return nil
	}
	items := make([]shippingpkg.PackageItem, 0, len(rawItems))
	for _, ri := range rawItems {
		m, ok := ri.(map[string]interface{})
		if !ok {
			continue
		}
		sku, _ := m["sku"].(string)
		if sku == "" {
			continue
		}
		qty := 1
		if q, ok := m["quantity"].(float64); ok && q > 0 {
			qty = int(q)
		}
		items = append(items, shippingpkg.PackageItem{SKU: sku, Quantity: qty})
	}
	return items
}

func (h *Handlers) applyPackageFromItems(ctx context.Context, businessID uint, items []shippingpkg.PackageItem, raw map[string]interface{}) {
	pkgs, ok := raw["packages"].([]interface{})
	if !ok || len(pkgs) == 0 {
		return
	}
	pkg, ok := pkgs[0].(map[string]interface{})
	if !ok {
		return
	}

	skus := make([]string, 0, len(items))
	for _, it := range items {
		skus = append(skus, it.SKU)
	}
	if dims, err := h.uc.Repo().GetProductDimensionsBySKUs(ctx, businessID, skus); err == nil {
		for i := range items {
			d, ok := dims[items[i].SKU]
			if !ok {
				continue
			}
			items[i].Length = d.Length
			items[i].Width = d.Width
			items[i].Height = d.Height
			items[i].Weight = d.Weight
		}
	}

	strategy := ""
	var boxes []shippingpkg.Box
	if cfg, err := h.uc.Repo().GetBusinessPackageConfig(ctx, businessID, requestWarehouseID(raw)); err == nil && cfg != nil {
		strategy = cfg.Strategy
		boxes = cfg.Boxes
	}

	setPackageDims(pkg, shippingpkg.Resolve(strategy, boxes, shippingpkg.PackageInput{Items: items}))
}

func requestWarehouseID(raw map[string]interface{}) *uint {
	switch v := raw["warehouse_id"].(type) {
	case float64:
		if v > 0 {
			id := uint(v)
			return &id
		}
	case string:
		if n, err := strconv.ParseUint(v, 10, 64); err == nil && n > 0 {
			id := uint(n)
			return &id
		}
	}
	return nil
}

func setPackageDims(pkg map[string]interface{}, resolved shippingpkg.ResolvedPackage) {
	pkg["weight"] = resolved.Weight
	pkg["height"] = resolved.Height
	pkg["width"] = resolved.Width
	pkg["length"] = resolved.Length
}
