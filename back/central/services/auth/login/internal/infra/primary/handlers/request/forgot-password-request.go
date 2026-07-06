package request

type ForgotPasswordRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Channel string `json:"channel" binding:"omitempty,oneof=email whatsapp"`
}

type RecoveryChannelsRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6,numeric"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}
