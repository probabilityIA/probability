package request

// RechargeWalletRequest cuerpo de petición para recargar billetera
type RechargeWalletRequest struct {
	Amount     float64 `json:"amount" binding:"required"`
	BusinessID *uint   `json:"business_id"` // Requerido para super admin
}

// ManualDebitRequest cuerpo de petición para débito manual (admin)
type ManualDebitRequest struct {
	BusinessID uint    `json:"business_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Reference  string  `json:"reference"`
}

// DebitForGuideRequest cuerpo de petición para débito por guía
type DebitForGuideRequest struct {
	Amount         float64 `json:"amount" binding:"required"`
	TrackingNumber string  `json:"tracking_number" binding:"required"`
}
