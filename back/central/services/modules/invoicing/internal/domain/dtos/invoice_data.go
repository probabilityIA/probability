package dtos

type InvoiceCustomerData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	DNI     string `json:"dni"`
	Address string `json:"address,omitempty"`
}
type InvoiceItemData struct {
	ProductID                *string  `json:"product_id"`
	SKU                      string   `json:"sku"`
	Name                     string   `json:"name"`
	Description              *string  `json:"description"`
	Quantity                 int      `json:"quantity"`
	UnitPrice                float64  `json:"unit_price"`
	UnitPriceBase            float64  `json:"unit_price_base"`
	TotalPrice               float64  `json:"total_price"`
	Tax                      float64  `json:"tax"`
	TaxRate                  *float64 `json:"tax_rate"`
	Discount                 float64  `json:"discount"`
	DiscountPercent          float64  `json:"discount_percent"`
	UnitPricePresentment     float64  `json:"unit_price_presentment"`
	UnitPriceBasePresentment float64  `json:"unit_price_base_presentment"`
	TotalPricePresentment    float64  `json:"total_price_presentment"`
	DiscountPresentment      float64  `json:"discount_presentment"`
	TaxPresentment           float64  `json:"tax_presentment"`
}

func InvoiceShippingCost(order *OrderData) float64 {
	if order == nil {
		return 0
	}
	if !order.IsCOD || order.CodCarrierFee <= 0 {
		return order.ShippingCost
	}
	return order.ShippingCost + order.CodCarrierFee
}

func InvoiceTotalAmount(order *OrderData) float64 {
	if order == nil {
		return 0
	}
	if !order.IsCOD || order.CodCarrierFee <= 0 {
		return order.TotalAmount
	}
	return order.TotalAmount + order.CodCarrierFee
}

func ShippingCostBase(shippingCost float64, taxRate float64) float64 {
	if taxRate > 0 {
		return shippingCost / (1 + taxRate)
	}
	return shippingCost
}

type InvoiceData struct {
	IntegrationID    uint                   `json:"integration_id"`
	Customer         InvoiceCustomerData    `json:"customer"`
	Items            []InvoiceItemData      `json:"items"`
	Total            float64                `json:"total"`
	Subtotal         float64                `json:"subtotal"`
	Tax              float64                `json:"tax"`
	Discount         float64                `json:"discount"`
	ShippingCost     float64                `json:"shipping_cost"`
	ShippingDiscount float64                `json:"shipping_discount"`
	FreeShipping     bool                   `json:"free_shipping"`
	ShippingCostBase float64                `json:"shipping_cost_base"`
	Currency         string                 `json:"currency"`
	OrderID          string                 `json:"order_id"`
	OrderNumber      string                 `json:"order_number,omitempty"`
	Config           map[string]interface{} `json:"config"`
}
