package dtos

// CheckoutItemInput is a cart line item as sent by the client (price is never trusted).
type CheckoutItemInput struct {
	ProductID string
	Quantity  int
}

// CheckoutAddressInput is the shipping address as sent by the client.
type CheckoutAddressInput struct {
	Street       string
	City         string
	State        string
	Country      string
	PostalCode   string
	Instructions string
}

// CreateCheckoutDTO is the input to start a checkout session (before payment).
type CreateCheckoutDTO struct {
	Items         []CheckoutItemInput
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	CustomerDni   string
	Address       *CheckoutAddressInput
}

// CheckoutSessionDTO is what the client needs to open the Bold checkout widget.
type CheckoutSessionDTO struct {
	Reference      string  `json:"reference"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Hash           string  `json:"hash"`
	PublicKey      string  `json:"public_key"`
	RedirectionURL string  `json:"redirection_url"`
	IsSandbox      bool    `json:"is_sandbox"`
	PollingEnabled bool    `json:"polling_enabled"`
}
