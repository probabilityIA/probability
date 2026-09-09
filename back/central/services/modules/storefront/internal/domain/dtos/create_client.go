package dtos

type CreateClientDTO struct {
	Name     string
	Email    string
	Password string
	Phone    string
	Dni      *string
}
