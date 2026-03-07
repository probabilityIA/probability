package request

// CreateOrderRequest payload para crear una orden desde el storefront
type CreateOrderRequest struct {
	Items   []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
	Notes   *string            `json:"notes"`
	Address *AddressRequest    `json:"address"`
}

// OrderItemRequest item de la orden
type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// AddressRequest direccion de envio
type AddressRequest struct {
	FirstName    string  `json:"first_name" binding:"required"`
	LastName     string  `json:"last_name"`
	Phone        string  `json:"phone"`
	Street       string  `json:"street" binding:"required"`
	Street2      string  `json:"street2"`
	City         string  `json:"city" binding:"required"`
	State        string  `json:"state"`
	Country      string  `json:"country"`
	PostalCode   string  `json:"postal_code"`
	Instructions *string `json:"instructions"`
}
