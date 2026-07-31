package handlers

import (
	"sync"

	"github.com/secamc93/probability/back/testing/shared/log"
)

const (
	MockAccessToken = "clt.mock-token-tiktok"
	MockShopName    = "probability-mock-shop"
)

type shopData struct {
	ShopName             string `json:"shop_name"`
	ShopRating           string `json:"shop_rating"`
	ShopReviewCount      int64  `json:"shop_review_count"`
	ItemSoldCount        int64  `json:"item_sold_count"`
	ShopID               int64  `json:"shop_id"`
	ShopPerformanceValue int64  `json:"shop_performance_value"`
}

type Handler struct {
	logger log.ILogger

	mu    sync.Mutex
	shops map[string]*shopData
}

func New(logger log.ILogger) *Handler {
	h := &Handler{
		logger: logger,
		shops:  make(map[string]*shopData),
	}
	h.seed()
	return h
}

func (h *Handler) seed() {
	h.shops[MockShopName] = &shopData{
		ShopName:             MockShopName,
		ShopRating:           "4.7",
		ShopReviewCount:      1289,
		ItemSoldCount:        5430,
		ShopID:               7495623810245678901,
		ShopPerformanceValue: 87,
	}
	h.shops["tienda-demo"] = &shopData{
		ShopName:             "tienda-demo",
		ShopRating:           "4.2",
		ShopReviewCount:      342,
		ItemSoldCount:        980,
		ShopID:               7495623810245678902,
		ShopPerformanceValue: 65,
	}
}
