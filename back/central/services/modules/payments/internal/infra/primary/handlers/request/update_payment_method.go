package request

// UpdatePaymentMethod representa la solicitud HTTP para actualizar un método de pago
type UpdatePaymentMethod struct {
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"required,oneof=card digital_wallet bank_transfer cash"`
	Provider    string `json:"provider" binding:"max=64"`
	Icon        string `json:"icon" binding:"max=255"`
	Color       string `json:"color" binding:"max=32"`
}
