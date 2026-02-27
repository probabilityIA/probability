package dtos

// RechargeWalletDTO datos para solicitar una recarga
type RechargeWalletDTO struct {
	BusinessID uint
	Amount     float64
}

// ManualDebitDTO datos para débito manual (admin)
type ManualDebitDTO struct {
	BusinessID uint
	Amount     float64
	Reference  string
}

// DebitForGuideDTO datos para débito por generación de guía
type DebitForGuideDTO struct {
	BusinessID     uint
	Amount         float64
	TrackingNumber string
}
