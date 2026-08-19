package handlers

import (
	"strings"
	"time"
)

type seedItem struct {
	name  string
	sku   string
	price float64
	stock int
}

var defaultSeed = []seedItem{
	{name: "Camiseta Padel White - XL / Blanco", sku: "8013XL", price: 89000, stock: 12},
	{name: "Pantaloneta Padel Negra - M / Negro", sku: "8019M", price: 75000, stock: 4},
	{name: "Camiseta Motion Black - L / Negro", sku: "8227L", price: 92000, stock: 0},
	{name: "Proteina Whey 2lb", sku: "52205", price: 180000, stock: 25},
	{name: "Producto solo en Tiendanube", sku: "TN-ONLY-001", price: 45000, stock: 7},
}

func (h *Handler) seedDefaults() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.seedDefaultsLocked()
}

func (h *Handler) seedDefaultsLocked() {
	for _, item := range defaultSeed {
		stock := item.stock
		h.addProductLocked(item.name, item.sku, item.price, &stock)
	}
}

func (h *Handler) addProductLocked(name, sku string, price float64, stock *int) *product {
	now := time.Now().Format(time.RFC3339)

	nuevo := &product{
		ID:          h.nextProductID,
		Name:        localized{"es": strings.TrimSpace(name)},
		Description: localized{"es": ""},
		Published:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.nextProductID++

	nueva := variant{
		ID:        h.nextVariantID,
		ProductID: nuevo.ID,
		SKU:       strings.TrimSpace(sku),
		Price:     money(price),
		Weight:    money(0),
		Depth:     money(0),
		Width:     money(0),
		Height:    money(0),
	}
	h.nextVariantID++

	if stock != nil {
		nueva.Stock = *stock
		nueva.StockManagement = true
	}

	nuevo.Variants = []variant{nueva}
	nuevo.Images = []image{{
		ID:        h.nextImageID,
		ProductID: nuevo.ID,
		Src:       "https://mock.probability.test/tiendanube/" + strings.TrimSpace(sku) + ".jpg",
		Position:  1,
	}}
	h.nextImageID++

	h.products[nuevo.ID] = nuevo
	h.order = append(h.order, nuevo.ID)
	return nuevo
}
