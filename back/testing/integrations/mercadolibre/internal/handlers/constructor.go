package handlers

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/shared/log"
)

const (
	MockSellerID     = 111222333
	MockSiteID       = "MCO"
	MockAccessToken  = "APP_USR-mock-access-token"
	MockRefreshToken = "TG-mock-refresh-token"
)

type itemAttribute struct {
	ID        string `json:"id"`
	ValueName string `json:"value_name"`
}

type itemVariation struct {
	ID                int64           `json:"id"`
	UserProductID     string          `json:"user_product_id"`
	AvailableQuantity int             `json:"available_quantity"`
	SellerCustomField string          `json:"seller_custom_field"`
	Attributes        []itemAttribute `json:"attributes"`
}

type item struct {
	ID                string          `json:"id"`
	SiteID            string          `json:"site_id"`
	Title             string          `json:"title"`
	Price             float64         `json:"price"`
	CurrencyID        string          `json:"currency_id"`
	AvailableQuantity int             `json:"available_quantity"`
	SoldQuantity      int             `json:"sold_quantity"`
	Status            string          `json:"status"`
	SubStatus         []string        `json:"sub_status"`
	Permalink         string          `json:"permalink"`
	SellerID          int64           `json:"seller_id"`
	SellerCustomField string          `json:"seller_custom_field"`
	UserProductID     string          `json:"user_product_id"`
	CatalogListing    bool            `json:"catalog_listing"`
	ListingTypeID     string          `json:"listing_type_id"`
	Shipping          itemShipping    `json:"shipping"`
	Attributes        []itemAttribute `json:"attributes"`
	Variations        []itemVariation `json:"variations"`
}

type itemShipping struct {
	LogisticType string   `json:"logistic_type"`
	Tags         []string `json:"tags"`
	Mode         string   `json:"mode"`
	FreeShipping bool     `json:"free_shipping"`
}

type userProductStock struct {
	UserProductID string
	Version       string
	Locations     []map[string]interface{}
}

type Handler struct {
	logger log.ILogger

	mu          sync.Mutex
	items       map[string]*item
	stocks      map[string]*userProductStock
	nextItemNum int64
	nextOrderID int64
	webhookURL  string
}

func New(logger log.ILogger, webhookURL string) *Handler {
	h := &Handler{
		logger:      logger,
		items:       make(map[string]*item),
		stocks:      make(map[string]*userProductStock),
		nextItemNum: 1,
		nextOrderID: 2000000001,
		webhookURL:  webhookURL,
	}
	h.seed("Camiseta Mock Meli", "MELI-MOCK-001", 75000, 10)
	h.seed("Gorra Mock Meli", "MELI-MOCK-002", 35000, 25)
	h.seed("Producto Meli sin SKU", "", 12000, 4)
	return h
}

func (h *Handler) seed(title, sku string, price float64, qty int) *item {
	id := h.newItemID()
	userProductID := "UP" + strconv.FormatInt(h.nextItemNum, 10)

	it := &item{
		ID:                id,
		SiteID:            MockSiteID,
		Title:             title,
		Price:             price,
		CurrencyID:        "COP",
		AvailableQuantity: qty,
		Status:            "active",
		SubStatus:         []string{},
		Permalink:         "http://meli-mock.local/" + id,
		SellerID:          MockSellerID,
		SellerCustomField: sku,
		UserProductID:     userProductID,
		ListingTypeID:     "gold_special",
		Shipping:          itemShipping{LogisticType: "drop_off", Tags: []string{}, Mode: "me2"},
		Attributes:        []itemAttribute{},
		Variations:        []itemVariation{},
	}
	if sku != "" {
		it.Attributes = append(it.Attributes, itemAttribute{ID: "SELLER_SKU", ValueName: sku})
	}

	h.items[id] = it
	h.stocks[userProductID] = &userProductStock{
		UserProductID: userProductID,
		Version:       "v1",
		Locations:     []map[string]interface{}{},
	}
	return it
}

func (h *Handler) newItemID() string {
	id := MockSiteID + strconv.FormatInt(1000000000+h.nextItemNum, 10)
	h.nextItemNum++
	return id
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)

	router.POST("/oauth/token", h.handleToken)
	router.GET("/users/me", h.handleUserMe)

	router.GET("/users/:user_id/items/search", h.handleSearchItems)

	router.GET("/items", h.handleMultigetItems)
	router.POST("/items", h.handleCreateItem)
	router.GET("/items/:id", h.handleGetItem)
	router.PUT("/items/:id", h.handleUpdateItem)

	router.GET("/user-products/:id/stock", h.handleGetUserProductStock)
	router.PUT("/user-products/:id/stock", h.handleUpdateUserProductStock)

	router.GET("/orders/:id", h.handleGetOrder)
	router.GET("/orders/search", h.handleSearchOrders)

	router.POST("/simulate/notification", h.handleSimulateNotification)
}
