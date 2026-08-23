package handlers

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/shared/log"
)

const (
	MockStoreID     = "1234567"
	MockAccessToken = "mock-tiendanube-token"
)

type localized map[string]string

type variant struct {
	ID              int64       `json:"id"`
	ProductID       int64       `json:"product_id"`
	SKU             string      `json:"sku"`
	Barcode         string      `json:"barcode"`
	Price           string      `json:"price"`
	Stock           interface{} `json:"stock"`
	StockManagement bool        `json:"stock_management"`
	Weight          string      `json:"weight"`
	Depth           string      `json:"depth"`
	Width           string      `json:"width"`
	Height          string      `json:"height"`
}

type image struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Src       string `json:"src"`
	Position  int    `json:"position"`
}

type product struct {
	ID          int64     `json:"id"`
	Name        localized `json:"name"`
	Description localized `json:"description"`
	Published   bool      `json:"published"`
	Variants    []variant `json:"variants"`
	Images      []image   `json:"images"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type Handler struct {
	logger log.ILogger

	mu            sync.Mutex
	products      map[int64]*product
	order         []int64
	nextProductID int64
	nextVariantID int64
	nextImageID   int64

	orders      map[int64]*mockOrder
	orderIDs    []int64
	nextOrderID int64
}

func New(logger log.ILogger) *Handler {
	h := &Handler{
		logger:        logger,
		products:      make(map[int64]*product),
		order:         make([]int64, 0),
		nextProductID: 2001,
		nextVariantID: 9001,
		nextImageID:   7001,
		orders:        make(map[int64]*mockOrder),
		orderIDs:      make([]int64, 0),
		nextOrderID:   5001,
	}
	h.seedDefaults()
	h.mu.Lock()
	h.seedOrdersLocked()
	h.mu.Unlock()
	return h
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)

	api := router.Group("/:version/:store_id", h.authMiddleware())
	{
		api.GET("/store", h.handleStore)
		api.GET("/products", h.handleListProducts)
		api.POST("/products", h.handleCreateProduct)
		api.GET("/products/:id", h.handleGetProduct)
		api.PUT("/products/:id", h.handleUpdateProduct)
		api.PUT("/products/:id/variants/:vid", h.handleUpdateVariant)
		api.GET("/orders", h.handleListOrders)
		api.GET("/orders/:id", h.handleGetOrder)
	}

	router.POST("/mock/reset", h.handleReset)
	router.POST("/mock/seed-products", h.handleSeedProducts)
	router.POST("/mock/seed-orders", h.handleSeedOrders)
}

func money(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
