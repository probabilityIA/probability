package request

type CreateClientRequest struct {
	Name     string  `json:"name" binding:"required,min=2"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"omitempty,min=6"`
	Phone    string  `json:"phone"`
	Dni      *string `json:"dni"`
}
