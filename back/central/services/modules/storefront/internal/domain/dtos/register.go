package dtos

// RegisterDTO contains data for client self-registration
type RegisterDTO struct {
	Name         string
	Email        string
	Password     string
	Phone        string
	Dni          *string
	BusinessCode string
}
