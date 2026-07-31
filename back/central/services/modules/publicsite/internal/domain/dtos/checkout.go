package dtos

type CheckoutItemInput struct {
	ProductID string
	Quantity  int
}

type CheckoutAddressInput struct {
	Street       string
	City         string
	State        string
	Country      string
	PostalCode   string
	Instructions string
}

type CreateCheckoutDTO struct {
	Items         []CheckoutItemInput
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	CustomerDni   string
	PaymentMethod string
	Address       *CheckoutAddressInput
}

type CheckoutSessionDTO struct {
	PaymentMethod  string  `json:"payment_method"`
	Reference      string  `json:"reference"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Hash           string  `json:"hash"`
	PublicKey      string  `json:"public_key"`
	RedirectionURL string  `json:"redirection_url"`
	IsSandbox      bool    `json:"is_sandbox"`
	PollingEnabled bool    `json:"polling_enabled"`
}
