package dtos

// InvoiceCustomerData datos del cliente para facturación
type InvoiceCustomerData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	DNI     string `json:"dni"`
	Address string `json:"address,omitempty"`
}

// InvoiceItemData datos de un item para facturación
type InvoiceItemData struct {
	ProductID   *string  `json:"product_id"`
	SKU         string   `json:"sku"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Quantity    int      `json:"quantity"`
	UnitPrice   float64  `json:"unit_price"`
	UnitPriceBase float64 `json:"unit_price_base"` // Precio sin impuestos
	TotalPrice  float64  `json:"total_price"`
	Tax         float64  `json:"tax"`
	TaxRate     *float64 `json:"tax_rate"`
	Discount        float64  `json:"discount"`
	DiscountPercent float64  `json:"discount_percent"`
	// Precios en moneda presentment (moneda local, ej: COP)
	UnitPricePresentment      float64 `json:"unit_price_presentment"`
	UnitPriceBasePresentment  float64 `json:"unit_price_base_presentment"`
	TotalPricePresentment     float64 `json:"total_price_presentment"`
	DiscountPresentment       float64 `json:"discount_presentment"`
	TaxPresentment            float64 `json:"tax_presentment"`
}

// ShippingCostBase calcula el costo de envío sin impuestos
func ShippingCostBase(shippingCost float64, taxRate float64) float64 {
	if taxRate > 0 {
		return shippingCost / (1 + taxRate)
	}
	return shippingCost
}

// InvoiceData datos completos para crear factura (viaja por RabbitMQ)
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
	ShippingCostBase float64                `json:"shipping_cost_base"` // Envío sin impuestos
	Currency         string                 `json:"currency"`
	OrderID          string                 `json:"order_id"`
	OrderNumber      string                 `json:"order_number,omitempty"`
	Config           map[string]interface{} `json:"config"`
}
