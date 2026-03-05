package response

// OrderItemResponse representa un item de la orden en la respuesta HTTP
type OrderItemResponse struct {
	ID           uint     `json:"id"`
	ProductID    *string  `json:"product_id,omitempty"`
	ProductSKU   string   `json:"product_sku"`
	ProductName  string   `json:"product_name"`
	ProductTitle string   `json:"product_title,omitempty"`
	VariantID    *string  `json:"variant_id,omitempty"`
	Quantity     int      `json:"quantity"`
	UnitPrice    float64  `json:"unit_price"`
	TotalPrice   float64  `json:"total_price"`
	Currency     string   `json:"currency"`
	Discount     float64  `json:"discount"`
	Tax          float64  `json:"tax"`
	TaxRate      *float64 `json:"tax_rate,omitempty"`
	// Precios en moneda presentment
	UnitPricePresentment  float64 `json:"unit_price_presentment,omitempty"`
	TotalPricePresentment float64 `json:"total_price_presentment,omitempty"`
	DiscountPresentment   float64 `json:"discount_presentment,omitempty"`
	TaxPresentment        float64 `json:"tax_presentment,omitempty"`
	// Extras
	ImageURL    *string `json:"image_url,omitempty"`
	ProductURL  *string `json:"product_url,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`
}
