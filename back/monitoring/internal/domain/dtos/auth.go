package dtos

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
	Name  string
	Email string
}
