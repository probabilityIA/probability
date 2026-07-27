package handlers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/shared/log"
)

const MockAccountName = "probabilitymock"

type sku struct {
	ID              int        `json:"Id"`
	ProductID       int        `json:"ProductId"`
	Name            string     `json:"Name"`
	RefID           string     `json:"RefId"`
	IsActive        bool       `json:"IsActive"`
	WeightKg        float64    `json:"WeightKg"`
	Length          float64    `json:"Length"`
	Width           float64    `json:"Width"`
	Height          float64    `json:"Height"`
	MeasurementUnit string     `json:"MeasurementUnit"`
	ImageURL        string     `json:"ImageUrl"`
	Images          []skuImage `json:"Images"`
}

type skuImage struct {
	ImageURL string `json:"ImageUrl"`
}

type warehouse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type hookConfig struct {
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers"`
	Statuses []string          `json:"statuses"`
}

type Handler struct {
	logger log.ILogger

	mu          sync.Mutex
	skus        map[int]*sku
	warehouses  []*warehouse
	stock       map[string]int
	hook        *hookConfig
	nextSKUID   int
	nextOrderID int
	webhookURL  string
}

func New(logger log.ILogger, webhookURL string) *Handler {
	h := &Handler{
		logger:      logger,
		skus:        make(map[int]*sku),
		stock:       make(map[string]int),
		nextSKUID:   1001,
		nextOrderID: 5000,
		webhookURL:  webhookURL,
	}
	h.warehouses = []*warehouse{
		{ID: "1_1", Name: "Bodega Mock VTEX Principal", IsActive: true},
		{ID: "1_2", Name: "Bodega Mock VTEX Secundaria", IsActive: true},
	}
	h.seed("Camiseta Mock VTEX", "VTEX-MOCK-001", 10)
	h.seed("Gorra Mock VTEX", "VTEX-MOCK-002", 25)
	h.seed("Producto VTEX sin RefId", "", 4)
	return h
}

func (h *Handler) seed(name, refID string, quantity int) *sku {
	id := h.nextSKUID
	h.nextSKUID++

	item := &sku{
		ID:              id,
		ProductID:       id + 90000,
		Name:            name,
		RefID:           refID,
		IsActive:        true,
		WeightKg:        1.2,
		Length:          30,
		Width:           20,
		Height:          10,
		MeasurementUnit: "un",
		ImageURL:        "https://mock.probability.test/vtex/" + refID + ".jpg",
		Images:          []skuImage{{ImageURL: "https://mock.probability.test/vtex/" + refID + ".jpg"}},
	}
	h.skus[id] = item
	h.stock[stockKey(id, h.warehouses[0].ID)] = quantity
	return item
}

func stockKey(skuID int, warehouseID string) string {
	return warehouseID + ":" + itoa(skuID)
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	buf := make([]byte, 0, 12)
	for v > 0 {
		buf = append([]byte{byte('0' + v%10)}, buf...)
		v /= 10
	}
	if neg {
		return "-" + string(buf)
	}
	return string(buf)
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)

	catalog := router.Group("/api/catalog_system/pvt")
	{
		catalog.GET("/products/GetProductAndSkuIds", h.handleGetProductAndSkuIds)
		catalog.GET("/sku/stockkeepingunitids", h.handleListSKUIDs)
		catalog.GET("/sku/stockkeepingunitidbyrefid/:refid", h.handleSKUIDByRefID)
	}

	router.GET("/api/catalog/pvt/stockkeepingunit/:id", h.handleGetSKU)
	router.GET("/api/catalog-seller-portal/skus/_search", h.handleSellerSKUSearch)

	logistics := router.Group("/api/logistics/pvt")
	{
		logistics.GET("/configuration/warehouses", h.handleListWarehouses)
		logistics.PUT("/inventory/skus/:sku/warehouses/:warehouse", h.handleUpdateInventory)
	}

	oms := router.Group("/api/oms/pvt")
	{
		oms.GET("/orders", h.handleListOrders)
		oms.GET("/orders/:id", h.handleGetOrder)
	}

	hook := router.Group("/api/orders/hook")
	{
		hook.GET("/config", h.handleGetHook)
		hook.POST("/config", h.handleSetHook)
		hook.PUT("/config", h.handleSetHook)
		hook.DELETE("/config", h.handleDeleteHook)
	}

	router.POST("/simulate/order", h.handleSimulateOrder)
	router.POST("/mock/seed-products", h.handleSeedProducts)
}
