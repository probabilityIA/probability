package request

type WalletBalanceAlertEvent struct {
	BusinessID   uint    `json:"business_id"`
	BusinessName string  `json:"business_name"`
	PhoneNumber  string  `json:"phone_number"`
	Balance      float64 `json:"balance"`
}
