package dtos

type GenerateOrdersDTO struct {
	Count            int    `json:"count" binding:"required,min=1,max=20"`
	IntegrationID    uint   `json:"integration_id"`
	RandomProducts   bool   `json:"random_products"`
	MaxItemsPerOrder int    `json:"max_items_per_order"`
	Topic            string `json:"topic"`
	// Override opcionales para datos del cliente (si vacíos, se generan random)
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
}

func (d *GenerateOrdersDTO) ApplyDefaults() {
	if d.MaxItemsPerOrder <= 0 || d.MaxItemsPerOrder > 5 {
		d.MaxItemsPerOrder = 3
	}
	if d.Count <= 0 {
		d.Count = 1
	}
	if d.Count > 20 {
		d.Count = 20
	}
}
