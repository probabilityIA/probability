package domain

import "time"

type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
	IsTesting       bool
	BaseURL         string
	BaseURLTest     string
}

type Credential struct {
	AccessToken string
	BaseURL     string
}

type ShopInfo struct {
	ShopID               int64
	ShopName             string
	ShopRating           string
	ShopReviewCount      int64
	ItemSoldCount        int64
	ShopPerformanceValue int64
}

type ShopSnapshot struct {
	ID                   uint
	BusinessID           uint
	IntegrationID        uint
	ShopID               int64
	ShopName             string
	ShopRating           string
	ShopReviewCount      int64
	ItemSoldCount        int64
	ShopPerformanceValue int64
	CreatedAt            time.Time
}
