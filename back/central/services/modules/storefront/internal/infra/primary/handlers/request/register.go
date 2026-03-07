package request

// RegisterRequest payload para auto-registro de cliente final
type RegisterRequest struct {
	Name         string  `json:"name" binding:"required,min=2,max=255"`
	Email        string  `json:"email" binding:"required,email,max=255"`
	Password     string  `json:"password" binding:"required,min=6,max=100"`
	Phone        string  `json:"phone" binding:"omitempty,max=50"`
	Dni          *string `json:"dni" binding:"omitempty,max=50"`
	BusinessCode string  `json:"business_code" binding:"required"`
}
