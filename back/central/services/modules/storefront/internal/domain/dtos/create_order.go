package dtos

// StorefrontCreateOrderDTO contains data to create a storefront order
type StorefrontCreateOrderDTO struct {
	Items    []StorefrontOrderItemDTO
	Notes    *string
	Address  *StorefrontAddressDTO
}

// StorefrontOrderItemDTO represents an item in the order
type StorefrontOrderItemDTO struct {
	ProductID string
	Quantity  int
}

// StorefrontAddressDTO represents a shipping address
type StorefrontAddressDTO struct {
	FirstName    string
	LastName     string
	Phone        string
	Street       string
	Street2      string
	City         string
	State        string
	Country      string
	PostalCode   string
	Instructions *string
}
