package request

import "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"

type CheckoutItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type CheckoutAddressRequest struct {
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postal_code"`
	Instructions string `json:"instructions"`
}

type CreateCheckoutRequest struct {
	Items         []CheckoutItemRequest   `json:"items" binding:"required,min=1,dive"`
	CustomerName  string                  `json:"customer_name" binding:"required"`
	CustomerEmail string                  `json:"customer_email"`
	CustomerPhone string                  `json:"customer_phone"`
	CustomerDni   string                  `json:"customer_dni"`
	PaymentMethod string                  `json:"payment_method"`
	Address       *CheckoutAddressRequest `json:"address"`
}

func (r *CreateCheckoutRequest) ToDTO() *dtos.CreateCheckoutDTO {
	items := make([]dtos.CheckoutItemInput, 0, len(r.Items))
	for _, it := range r.Items {
		items = append(items, dtos.CheckoutItemInput{ProductID: it.ProductID, Quantity: it.Quantity})
	}
	var address *dtos.CheckoutAddressInput
	if r.Address != nil {
		address = &dtos.CheckoutAddressInput{
			Street:       r.Address.Street,
			City:         r.Address.City,
			State:        r.Address.State,
			Country:      r.Address.Country,
			PostalCode:   r.Address.PostalCode,
			Instructions: r.Address.Instructions,
		}
	}
	return &dtos.CreateCheckoutDTO{
		Items:         items,
		CustomerName:  r.CustomerName,
		CustomerEmail: r.CustomerEmail,
		CustomerPhone: r.CustomerPhone,
		CustomerDni:   r.CustomerDni,
		PaymentMethod: r.PaymentMethod,
		Address:       address,
	}
}
