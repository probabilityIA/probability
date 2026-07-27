package handlers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type webhook struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Topic       string `json:"topic"`
	DeliveryURL string `json:"delivery_url"`
	Secret      string `json:"secret"`
	DateCreated string `json:"date_created"`
}

type wooImage struct {
	ID  int64  `json:"id"`
	Src string `json:"src"`
}

type wooProduct struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	SKU           string     `json:"sku"`
	Type          string     `json:"type"`
	Price         string     `json:"price"`
	StockQuantity *int       `json:"stock_quantity"`
	ManageStock   bool       `json:"manage_stock"`
	Status        string     `json:"status"`
	Images        []wooImage `json:"images"`
}

type wooVariation struct {
	ID            int64  `json:"id"`
	SKU           string `json:"sku"`
	Price         string `json:"price"`
	StockQuantity *int   `json:"stock_quantity"`
	ManageStock   bool   `json:"manage_stock"`
}

type Handler struct {
	logger log.ILogger

	mu         sync.Mutex
	webhooks   map[int64]*webhook
	products   map[int64]*wooProduct
	variations map[int64][]*wooVariation
	nextID     int64
	nextProdID int64
	nextOrder  int64
}

func New(logger log.ILogger) *Handler {
	h := &Handler{
		logger:     logger,
		webhooks:   make(map[int64]*webhook),
		products:   make(map[int64]*wooProduct),
		variations: make(map[int64][]*wooVariation),
		nextID:     1,
		nextProdID: 100,
		nextOrder:  5000,
	}
	h.seedCatalog()
	return h
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)

	api := router.Group("/wp-json/wc/v3")
	{
		api.GET("/system_status", h.handleSystemStatus)
		api.GET("/orders", h.handleListOrders)
		api.GET("/orders/:id", h.handleGetOrder)
		api.GET("/products", h.handleListProducts)
		api.GET("/products/:id", h.handleGetProduct)
		api.POST("/products", h.handleCreateProduct)
		api.PUT("/products/:id", h.handleUpdateProduct)
		api.GET("/products/:id/variations", h.handleListVariations)
		api.PUT("/products/:id/variations/:vid", h.handleUpdateVariation)
		api.GET("/webhooks", h.handleListWebhooks)
		api.POST("/webhooks", h.handleCreateWebhook)
		api.DELETE("/webhooks/:id", h.handleDeleteWebhook)
	}

	router.POST("/simulate/order", h.handleSimulateOrder)
	router.POST("/mock/seed-products", h.handleSeedProducts)
}
