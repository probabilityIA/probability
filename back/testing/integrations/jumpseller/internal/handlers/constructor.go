package handlers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/shared/log"
)

const (
	MockStoreCode  = "probability-mock"
	MockHooksToken = "mock-hooks-token-jumpseller"
)

type hook struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Event     string `json:"event"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type product struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	SKU            string         `json:"sku"`
	Price          float64        `json:"price"`
	Stock          int            `json:"stock"`
	StockUnlimited bool           `json:"stock_unlimited"`
	Status         string         `json:"status"`
	Weight         float64        `json:"weight"`
	Height         float64        `json:"height"`
	Width          float64        `json:"width"`
	Length         float64        `json:"length"`
	PackageFormat  string         `json:"package_format"`
	Images         []productImage `json:"images"`
}

type productImage struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type Handler struct {
	logger log.ILogger

	mu         sync.Mutex
	hooks      map[int64]*hook
	products   map[int64]*product
	nextHookID int64
	nextProdID int64
	nextOrder  int64
}

func New(logger log.ILogger) *Handler {
	h := &Handler{
		logger:     logger,
		hooks:      make(map[int64]*hook),
		products:   make(map[int64]*product),
		nextHookID: 1,
		nextProdID: 100,
		nextOrder:  5000,
	}
	h.products[100] = &product{
		ID:            100,
		Name:          "Producto de prueba Jumpseller",
		SKU:           "JS-MOCK-001",
		Price:         75000,
		Stock:         10,
		Status:        "available",
		Weight:        1.5,
		Height:        12,
		Width:         20,
		Length:        30,
		PackageFormat: "box",
		Images:        []productImage{{ID: 100, URL: "https://mock.probability.test/jumpseller/JS-MOCK-001.jpg"}},
	}
	h.nextProdID = 101
	return h
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)

	router.GET("/store/info.json", h.handleStoreInfo)

	router.GET("/orders.json", h.handleListOrders)
	router.GET("/orders/:id", h.handleGetOrder)
	router.PUT("/orders/:id", h.handleUpdateOrder)

	router.GET("/hooks.json", h.handleListHooks)
	router.POST("/hooks.json", h.handleCreateHook)
	router.DELETE("/hooks/:id", h.handleDeleteHook)

	router.GET("/products.json", h.handleListProducts)
	router.GET("/products/search.json", h.handleSearchProducts)
	router.POST("/products.json", h.handleCreateProduct)
	router.PUT("/products/:id", h.handleUpdateProduct)
	router.PUT("/products/:id/variants/:vid", h.handleUpdateVariant)

	router.POST("/simulate/order", h.handleSimulateOrder)
	router.POST("/mock/seed-products", h.handleSeedProducts)
}
