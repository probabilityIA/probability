package response

import "time"

// WalletResponse respuesta HTTP de una billetera
type WalletResponse struct {
	ID         string  `json:"ID"`
	BusinessID uint    `json:"BusinessID"`
	Balance    float64 `json:"Balance"`
	CreatedAt  time.Time `json:"CreatedAt"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
}

// WalletTransactionResponse respuesta HTTP de una transacción de billetera
type WalletTransactionResponse struct {
	ID                   string    `json:"ID"`
	WalletID             string    `json:"WalletID"`
	Amount               float64   `json:"Amount"`
	Type                 string    `json:"Type"`
	Status               string    `json:"Status"`
	Reference            string    `json:"Reference"`
	QrCode               string    `json:"QrCode"`
	PaymentTransactionID *uint     `json:"PaymentTransactionID,omitempty"`
	IntegrationTypeID    *uint     `json:"integration_type_id,omitempty"`
	IntegrationID        *uint     `json:"integration_id,omitempty"`
	IntegrationName      string    `json:"integration_name,omitempty"`
	IntegrationImageURL  string    `json:"integration_image_url,omitempty"`
	GatewayRequest       any       `json:"gateway_request,omitempty"`
	GatewayResponse      any       `json:"gateway_response,omitempty"`
	CreatedAt            time.Time `json:"CreatedAt"`
}
