package request

type SubscriptionPaymentAlertEvent struct {
	BusinessID    uint    `json:"business_id"`
	BusinessName  string  `json:"business_name"`
	PhoneNumber   string  `json:"phone_number"`
	DueDate       string  `json:"due_date"`
	CycleAmount   float64 `json:"cycle_amount"`
	WalletBalance float64 `json:"wallet_balance"`
}
