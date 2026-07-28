package entities

import "encoding/json"

// StorefrontCheckoutSnapshot es una copia de solo lectura de public_checkouts (tabla
// propiedad del modulo publicsite), replicada localmente siguiendo la regla de
// aislamiento de modulos (backend-conventions.md seccion 1) para poder confirmar el
// pago y publicar la orden canonica desde el webhook de Bold.
type StorefrontCheckoutSnapshot struct {
	ID                  string
	BusinessID          uint
	IntegrationID       uint
	Slug                string
	Reference           string
	Status              string
	Amount              float64
	Currency            string
	ItemsJSON           json.RawMessage
	CustomerName        string
	CustomerEmail       string
	CustomerPhone       string
	CustomerDni         string
	ShippingAddressJSON json.RawMessage
}

type StorefrontCheckoutItem struct {
	ProductID string  `json:"product_id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	ImageURL  string  `json:"image_url"`
}

type StorefrontCheckoutAddress struct {
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postal_code"`
	Instructions string `json:"instructions"`
}
