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
	ID                   string  `json:"ID"`
	WalletID             string  `json:"WalletID"`
	Amount               float64 `json:"Amount"`
	Type                 string  `json:"Type"`
	Status               string  `json:"Status"`
	Reference            string  `json:"Reference"`
	QrCode               string  `json:"QrCode"`
	PaymentTransactionID *uint   `json:"PaymentTransactionID,omitempty"`
	CreatedAt            time.Time `json:"CreatedAt"`
}
