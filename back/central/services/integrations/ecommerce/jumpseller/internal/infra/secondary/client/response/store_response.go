package response

import (
	"net/url"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func decodeIfEscaped(value string) string {
	decoded, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}

type StoreInfoEnvelope struct {
	Store StoreInfoResponse `json:"store"`
}

type StoreInfoResponse struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Country          string `json:"country"`
	Currency         string `json:"currency"`
	HooksToken       string `json:"hooks_token"`
	WeightUnit       string `json:"weight_unit"`
	SubscriptionPlan string `json:"subscription_plan"`
}

type LocationEnvelope struct {
	Location LocationResponse `json:"location"`
}

type LocationResponse struct {
	ID              int64                   `json:"id"`
	Name            string                  `json:"name"`
	Main            bool                    `json:"main"`
	IsStockOrigin   bool                    `json:"is_stock_origin"`
	PickupPoint     bool                    `json:"pickup_point"`
	LocationAddress LocationAddressResponse `json:"location_address"`
}

type LocationAddressResponse struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func (l LocationEnvelope) ToDomain() domain.Location {
	return domain.Location{
		ID:            l.Location.ID,
		Name:          decodeIfEscaped(l.Location.Name),
		Main:          l.Location.Main,
		IsStockOrigin: l.Location.IsStockOrigin,
		PickupPoint:   l.Location.PickupPoint,
		City:          decodeIfEscaped(l.Location.LocationAddress.City),
		Country:       decodeIfEscaped(l.Location.LocationAddress.Country),
	}
}

func (s StoreInfoEnvelope) ToDomain() domain.StoreInfo {
	return domain.StoreInfo{
		Code:             s.Store.Code,
		Name:             s.Store.Name,
		URL:              s.Store.URL,
		Country:          s.Store.Country,
		Currency:         s.Store.Currency,
		HooksToken:       s.Store.HooksToken,
		WeightUnit:       s.Store.WeightUnit,
		SubscriptionPlan: s.Store.SubscriptionPlan,
	}
}

type HookEnvelope struct {
	Hook HookResponse `json:"hook"`
}

type HookResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Event     string `json:"event"`
	URL       string `json:"url"`
	CreatedAt JSTime `json:"created_at"`
}

type CreateHookRequest struct {
	Hook CreateHookFields `json:"hook"`
}

type CreateHookFields struct {
	Event string `json:"event"`
	URL   string `json:"url"`
}

type UpdateOrderRequest struct {
	Order UpdateOrderFields `json:"order"`
}

type UpdateOrderFields struct {
	Status          string `json:"status,omitempty"`
	ShipmentStatus  string `json:"shipment_status,omitempty"`
	TrackingNumber  string `json:"tracking_number,omitempty"`
	TrackingCompany string `json:"tracking_company,omitempty"`
	AdditionalInfo  string `json:"additional_information,omitempty"`
}
