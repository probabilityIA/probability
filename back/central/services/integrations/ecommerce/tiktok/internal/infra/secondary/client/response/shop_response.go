package response

import "github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type ShopData struct {
	ShopName             string `json:"shop_name"`
	ShopRating           string `json:"shop_rating"`
	ShopReviewCount      int64  `json:"shop_review_count"`
	ItemSoldCount        int64  `json:"item_sold_count"`
	ShopID               int64  `json:"shop_id"`
	ShopPerformanceValue int64  `json:"shop_performance_value"`
}

type ShopQueryResponse struct {
	Data struct {
		ShopData []ShopData `json:"shop_data"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (r *ShopQueryResponse) ToDomain() []domain.ShopInfo {
	shops := make([]domain.ShopInfo, 0, len(r.Data.ShopData))
	for _, s := range r.Data.ShopData {
		shops = append(shops, domain.ShopInfo{
			ShopID:               s.ShopID,
			ShopName:             s.ShopName,
			ShopRating:           s.ShopRating,
			ShopReviewCount:      s.ShopReviewCount,
			ItemSoldCount:        s.ItemSoldCount,
			ShopPerformanceValue: s.ShopPerformanceValue,
		})
	}
	return shops
}
