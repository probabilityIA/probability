package request

type CreateClientRequest struct {
	Name    string  `json:"name" binding:"required,min=2,max=255"`
	Email   *string `json:"email" binding:"omitempty,email,max=255"`
	Phone   string  `json:"phone" binding:"omitempty,max=20"`
	Dni     *string `json:"dni" binding:"omitempty,max=30"`
	Address *string `json:"address" binding:"omitempty,max=255"`
	City    *string `json:"city" binding:"omitempty,max=120"`
	Notes   *string `json:"notes" binding:"omitempty,max=1000"`
}

type UpdateClientRequest struct {
	Name    string  `json:"name" binding:"required,min=2,max=255"`
	Email   *string `json:"email" binding:"omitempty,email,max=255"`
	Phone   string  `json:"phone" binding:"omitempty,max=20"`
	Dni     *string `json:"dni" binding:"omitempty,max=30"`
	Address *string `json:"address" binding:"omitempty,max=255"`
	City    *string `json:"city" binding:"omitempty,max=120"`
	Notes   *string `json:"notes" binding:"omitempty,max=1000"`
}
