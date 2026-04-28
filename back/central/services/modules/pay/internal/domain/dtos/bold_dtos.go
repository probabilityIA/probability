package dtos

type BoldSignatureResponse struct {
	OrderID        string  `json:"order_id"`
	Hash           string  `json:"hash"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	PublicKey      string  `json:"public_key"`
	RedirectionURL string  `json:"redirection_url,omitempty"`
	IsSandbox      bool    `json:"is_sandbox"`
}

type BoldStatusResponse struct {
	BoldOrderID   string  `json:"bold_order_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
}

type BoldCredentials struct {
	APIKey            string
	SecretKey         string
	Environment       string
	BaseURL           string
	IntegrationTypeID uint
}

type BoldBusinessIntegration struct {
	IntegrationID     uint
	IntegrationTypeID uint
	IsTesting         bool
}

type BoldSimulateDTO struct {
	BusinessID uint
	OrderID    string
	Amount     float64
}

type BoldSimulateResponse struct {
	Success       bool    `json:"success"`
	OrderID       string  `json:"order_id"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	NewBalance    float64 `json:"new_balance"`
	Status        string  `json:"status"`
}
