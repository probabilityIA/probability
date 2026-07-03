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

type Handler struct {
	logger log.ILogger

	mu         sync.Mutex
	webhooks   map[int64]*webhook
	nextID     int64
	nextProdID int64
	nextOrder  int64
}

func New(logger log.ILogger) *Handler {
	return &Handler{
		logger:     logger,
		webhooks:   make(map[int64]*webhook),
		nextID:     1,
		nextProdID: 100,
		nextOrder:  5000,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)

	api := router.Group("/wp-json/wc/v3")
	{
		api.GET("/system_status", h.handleSystemStatus)
		api.GET("/orders", h.handleListOrders)
		api.GET("/orders/:id", h.handleGetOrder)
		api.POST("/products", h.handleCreateProduct)
		api.PUT("/products/:id", h.handleUpdateProduct)
		api.GET("/webhooks", h.handleListWebhooks)
		api.POST("/webhooks", h.handleCreateWebhook)
		api.DELETE("/webhooks/:id", h.handleDeleteWebhook)
	}

	router.POST("/simulate/order", h.handleSimulateOrder)
}
