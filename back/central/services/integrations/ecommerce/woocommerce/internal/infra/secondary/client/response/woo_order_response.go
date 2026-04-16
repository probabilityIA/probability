package response

import (
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// WooOrderResponse es la respuesta JSON de la API REST de WooCommerce v3 para órdenes.
type WooOrderResponse struct {
	ID                 int64                       `json:"id"`
	ParentID           int64                       `json:"parent_id"`
	Number             string                      `json:"number"`
	OrderKey           string                      `json:"order_key"`
	CreatedVia         string                      `json:"created_via"`
	Version            string                      `json:"version"`
	Status             string                      `json:"status"`
	Currency           string                      `json:"currency"`
	DateCreated        string                      `json:"date_created"`
	DateModified       string                      `json:"date_modified"`
	DiscountTotal      string                      `json:"discount_total"`
	DiscountTax        string                      `json:"discount_tax"`
	ShippingTotal      string                      `json:"shipping_total"`
	ShippingTax        string                      `json:"shipping_tax"`
	CartTax            string                      `json:"cart_tax"`
	Total              string                      `json:"total"`
	TotalTax           string                      `json:"total_tax"`
	PricesIncludeTax   bool                        `json:"prices_include_tax"`
	CustomerID         int64                       `json:"customer_id"`
	CustomerNote       string                      `json:"customer_note"`
	Billing            WooBillingResponse          `json:"billing"`
	Shipping           WooShippingResponse         `json:"shipping"`
	PaymentMethod      string                      `json:"payment_method"`
	PaymentMethodTitle string                      `json:"payment_method_title"`
	TransactionID      string                      `json:"transaction_id"`
	DatePaid           *string                     `json:"date_paid"`
	DateCompleted      *string                     `json:"date_completed"`
	LineItems          []WooLineItemResponse       `json:"line_items"`
	ShippingLines      []WooShippingLineResponse   `json:"shipping_lines"`
	FeeLines           []WooFeeLineResponse        `json:"fee_lines"`
	CouponLines        []WooCouponLineResponse     `json:"coupon_lines"`
	MetaData           []WooMetaDataResponse       `json:"meta_data"`
}

type WooBillingResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type WooShippingResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
}

type WooLineItemResponse struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	ProductID   int64                 `json:"product_id"`
	VariationID int64                 `json:"variation_id"`
	Quantity    int                   `json:"quantity"`
	TaxClass    string                `json:"tax_class"`
	Subtotal    string                `json:"subtotal"`
	SubtotalTax string                `json:"subtotal_tax"`
	Total       string                `json:"total"`
	TotalTax    string                `json:"total_tax"`
	SKU         string                `json:"sku"`
	Price       float64               `json:"price"`
	Image       *WooLineItemImage     `json:"image"`
	MetaData    []WooMetaDataResponse `json:"meta_data"`
}

type WooLineItemImage struct {
	ID  int64  `json:"id"`
	Src string `json:"src"`
}

type WooShippingLineResponse struct {
	ID          int64                 `json:"id"`
	MethodTitle string                `json:"method_title"`
	MethodID    string                `json:"method_id"`
	Total       string                `json:"total"`
	TotalTax    string                `json:"total_tax"`
	MetaData    []WooMetaDataResponse `json:"meta_data"`
}

type WooFeeLineResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	TaxClass  string `json:"tax_class"`
	TaxStatus string `json:"tax_status"`
	Total     string `json:"total"`
	TotalTax  string `json:"total_tax"`
}

type WooCouponLineResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Discount    string `json:"discount"`
	DiscountTax string `json:"discount_tax"`
}

type WooMetaDataResponse struct {
	ID    int64       `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ToDomain convierte la respuesta de la API a la entidad de dominio.
func (r *WooOrderResponse) ToDomain() domain.WooCommerceOrder {
	order := domain.WooCommerceOrder{
		ID:                 r.ID,
		ParentID:           r.ParentID,
		Number:             r.Number,
		OrderKey:           r.OrderKey,
		CreatedVia:         r.CreatedVia,
		Version:            r.Version,
		Status:             r.Status,
		Currency:           r.Currency,
		DateCreated:        parseWooDate(r.DateCreated),
		DateModified:       parseWooDate(r.DateModified),
		DiscountTotal:      r.DiscountTotal,
		DiscountTax:        r.DiscountTax,
		ShippingTotal:      r.ShippingTotal,
		ShippingTax:        r.ShippingTax,
		CartTax:            r.CartTax,
		Total:              r.Total,
		TotalTax:           r.TotalTax,
		PricesIncludeTax:   r.PricesIncludeTax,
		CustomerID:         r.CustomerID,
		CustomerNote:       r.CustomerNote,
		PaymentMethod:      r.PaymentMethod,
		PaymentMethodTitle: r.PaymentMethodTitle,
		TransactionID:      r.TransactionID,
		DatePaid:           parseWooDatePtr(r.DatePaid),
		DateCompleted:      parseWooDatePtr(r.DateCompleted),
		Billing: domain.WooCommerceBilling{
			FirstName: r.Billing.FirstName,
			LastName:  r.Billing.LastName,
			Company:   r.Billing.Company,
			Address1:  r.Billing.Address1,
			Address2:  r.Billing.Address2,
			City:      r.Billing.City,
			State:     r.Billing.State,
			Postcode:  r.Billing.Postcode,
			Country:   r.Billing.Country,
			Email:     r.Billing.Email,
			Phone:     r.Billing.Phone,
		},
		Shipping: domain.WooCommerceShipping{
			FirstName: r.Shipping.FirstName,
			LastName:  r.Shipping.LastName,
			Company:   r.Shipping.Company,
			Address1:  r.Shipping.Address1,
			Address2:  r.Shipping.Address2,
			City:      r.Shipping.City,
			State:     r.Shipping.State,
			Postcode:  r.Shipping.Postcode,
			Country:   r.Shipping.Country,
			Phone:     r.Shipping.Phone,
		},
	}

	// Line items
	order.LineItems = make([]domain.WooCommerceLineItem, len(r.LineItems))
	for i, item := range r.LineItems {
		imageURL := ""
		if item.Image != nil {
			imageURL = item.Image.Src
		}
		order.LineItems[i] = domain.WooCommerceLineItem{
			ID:          item.ID,
			Name:        item.Name,
			ProductID:   item.ProductID,
			VariationID: item.VariationID,
			Quantity:    item.Quantity,
			TaxClass:    item.TaxClass,
			Subtotal:    item.Subtotal,
			SubtotalTax: item.SubtotalTax,
			Total:       item.Total,
			TotalTax:    item.TotalTax,
			SKU:         item.SKU,
			Price:       item.Price,
			ImageURL:    imageURL,
			MetaData:    mapMetaData(item.MetaData),
		}
	}

	// Shipping lines
	order.ShippingLines = make([]domain.WooCommerceShippingLine, len(r.ShippingLines))
	for i, sl := range r.ShippingLines {
		order.ShippingLines[i] = domain.WooCommerceShippingLine{
			ID:          sl.ID,
			MethodTitle: sl.MethodTitle,
			MethodID:    sl.MethodID,
			Total:       sl.Total,
			TotalTax:    sl.TotalTax,
			MetaData:    mapMetaData(sl.MetaData),
		}
	}

	// Fee lines
	order.FeeLines = make([]domain.WooCommerceFeeLine, len(r.FeeLines))
	for i, fl := range r.FeeLines {
		order.FeeLines[i] = domain.WooCommerceFeeLine{
			ID:        fl.ID,
			Name:      fl.Name,
			TaxClass:  fl.TaxClass,
			TaxStatus: fl.TaxStatus,
			Total:     fl.Total,
			TotalTax:  fl.TotalTax,
		}
	}

	// Coupon lines
	order.CouponLines = make([]domain.WooCommerceCouponLine, len(r.CouponLines))
	for i, cl := range r.CouponLines {
		order.CouponLines[i] = domain.WooCommerceCouponLine{
			ID:          cl.ID,
			Code:        cl.Code,
			Discount:    cl.Discount,
			DiscountTax: cl.DiscountTax,
		}
	}

	// Meta data
	order.MetaData = mapMetaData(r.MetaData)

	return order
}

func mapMetaData(mds []WooMetaDataResponse) []domain.WooCommerceMetaData {
	result := make([]domain.WooCommerceMetaData, len(mds))
	for i, md := range mds {
		result[i] = domain.WooCommerceMetaData{
			ID:    md.ID,
			Key:   md.Key,
			Value: md.Value,
		}
	}
	return result
}

// parseWooDate parsea un string de fecha WooCommerce.
// WooCommerce usa formato "2024-01-15T10:30:00" (sin timezone).
func parseWooDate(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	// Intentar con RFC3339 primero
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t
	}
	// WooCommerce a veces omite timezone
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		return t
	}
	return time.Time{}
}

func parseWooDatePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t := parseWooDate(*s)
	if t.IsZero() {
		return nil
	}
	return &t
}

// parseFloat64 convierte un string a float64, retornando 0 si no es válido.
func parseFloat64(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
