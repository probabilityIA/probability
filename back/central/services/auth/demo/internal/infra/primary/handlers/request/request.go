package request

type DemoRegisterRequest struct {
	FullName     string `json:"full_name" binding:"required,min=2,max=120"`
	BusinessName string `json:"business_name" binding:"required,min=2,max=120"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6,max=100"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}
