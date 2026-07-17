package response

import (
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

type OrderEnvelope struct {
	Order JumpsellerOrderResponse `json:"order"`
}

type JumpsellerOrderResponse struct {
	ID                int64                   `json:"id"`
	CreatedAt         JSTime                  `json:"created_at"`
	Status            string                  `json:"status"`
	Currency          string                  `json:"currency"`
	Subtotal          float64                 `json:"subtotal"`
	Tax               float64                 `json:"tax"`
	ShippingTax       float64                 `json:"shipping_tax"`
	Shipping          float64                 `json:"shipping"`
	ShippingRequired  bool                    `json:"shipping_required"`
	Total             float64                 `json:"total"`
	Discount          float64                 `json:"discount"`
	ShippingDiscount  float64                 `json:"shipping_discount"`
	FulfillmentStatus *string                 `json:"fulfillment_status"`
	ShipmentStatus    string                  `json:"shipment_status"`
	ShippingMethodID  int64                   `json:"shipping_method_id"`
	ShippingMethod    string                  `json:"shipping_method_name"`
	PaymentMethodName string                  `json:"payment_method_name"`
	PaymentMethodType string                  `json:"payment_method_type"`
	PaymentInfo       string                  `json:"payment_information"`
	AdditionalInfo    string                  `json:"additional_information"`
	TrackingNumber    *string                 `json:"tracking_number"`
	TrackingCompany   *string                 `json:"tracking_company"`
	TrackingURL       *string                 `json:"tracking_url"`
	ShippingOption    string                  `json:"shipping_option"`
	Customer          OrderCustomerResponse   `json:"customer"`
	ShippingAddress   AddressResponse         `json:"shipping_address"`
	BillingAddress    AddressResponse         `json:"billing_address"`
	Products          []OrderProductResponse  `json:"products"`
	AdditionalFields  []OrderAddFieldResponse `json:"additional_fields"`
}

type OrderCustomerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	IP    string `json:"ip"`
}

type AddressResponse struct {
	Name         string   `json:"name"`
	Surname      string   `json:"surname"`
	TaxID        *string  `json:"taxid"`
	Address      string   `json:"address"`
	City         string   `json:"city"`
	Postal       string   `json:"postal"`
	Region       string   `json:"region"`
	Country      string   `json:"country"`
	CountryCode  string   `json:"country_code"`
	RegionCode   string   `json:"region_code"`
	StreetNumber *string  `json:"street_number"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type OrderProductResponse struct {
	ID        int64   `json:"id"`
	VariantID int64   `json:"variant_id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
	Tax       float64 `json:"tax"`
	Discount  float64 `json:"discount"`
	Weight    float64 `json:"weight"`
}

type OrderAddFieldResponse struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
	Area  string `json:"area"`
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (a AddressResponse) toDomain() domain.Address {
	return domain.Address{
		Name:         a.Name,
		Surname:      a.Surname,
		TaxID:        deref(a.TaxID),
		Address:      a.Address,
		StreetNumber: deref(a.StreetNumber),
		City:         a.City,
		Postal:       a.Postal,
		Region:       a.Region,
		Country:      a.Country,
		CountryCode:  a.CountryCode,
		RegionCode:   a.RegionCode,
		Latitude:     a.Latitude,
		Longitude:    a.Longitude,
	}
}

func (o JumpsellerOrderResponse) ToDomain() domain.JumpsellerOrder {
	products := make([]domain.OrderProduct, 0, len(o.Products))
	for _, p := range o.Products {
		products = append(products, domain.OrderProduct{
			ID:        p.ID,
			VariantID: p.VariantID,
			SKU:       p.SKU,
			Name:      p.Name,
			Qty:       p.Qty,
			Price:     p.Price,
			Tax:       p.Tax,
			Discount:  p.Discount,
			Weight:    p.Weight,
		})
	}

	fields := make([]domain.OrderAdditionalField, 0, len(o.AdditionalFields))
	for _, f := range o.AdditionalFields {
		fields = append(fields, domain.OrderAdditionalField{
			ID:    f.ID,
			Label: f.Label,
			Value: f.Value,
			Area:  f.Area,
		})
	}

	return domain.JumpsellerOrder{
		ID:                o.ID,
		CreatedAt:         o.CreatedAt.Time,
		Status:            o.Status,
		Currency:          o.Currency,
		Subtotal:          o.Subtotal,
		Tax:               o.Tax,
		ShippingTax:       o.ShippingTax,
		Shipping:          o.Shipping,
		ShippingRequired:  o.ShippingRequired,
		Total:             o.Total,
		Discount:          o.Discount,
		ShippingDiscount:  o.ShippingDiscount,
		FulfillmentStatus: deref(o.FulfillmentStatus),
		ShipmentStatus:    o.ShipmentStatus,
		ShippingMethodID:  o.ShippingMethodID,
		ShippingMethod:    o.ShippingMethod,
		PaymentMethodName: o.PaymentMethodName,
		PaymentMethodType: o.PaymentMethodType,
		PaymentInfo:       o.PaymentInfo,
		AdditionalInfo:    o.AdditionalInfo,
		TrackingNumber:    deref(o.TrackingNumber),
		TrackingCompany:   deref(o.TrackingCompany),
		TrackingURL:       deref(o.TrackingURL),
		ShippingOption:    o.ShippingOption,
		Customer: domain.OrderCustomer{
			ID:    o.Customer.ID,
			Name:  o.Customer.Name,
			Email: o.Customer.Email,
			Phone: o.Customer.Phone,
			IP:    o.Customer.IP,
		},
		ShippingAddress:  o.ShippingAddress.toDomain(),
		BillingAddress:   o.BillingAddress.toDomain(),
		Products:         products,
		AdditionalFields: fields,
	}
}
