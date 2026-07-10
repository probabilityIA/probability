package request

type DemoRegisterRequest struct {
	FullName     string `json:"full_name" binding:"required,min=2,max=120"`
	BusinessName string `json:"business_name" binding:"required,min=2,max=120"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6,max=100"`
	Phone        string `json:"phone" binding:"omitempty,min=7,max=20"`
	Channel      string `json:"channel" binding:"omitempty,oneof=email whatsapp"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type DemoResendRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Channel string `json:"channel" binding:"omitempty,oneof=email whatsapp"`
	Phone   string `json:"phone" binding:"required_if=Channel whatsapp,omitempty,min=7,max=20"`
}

type DemoVerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6,numeric"`
}
