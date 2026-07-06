package domain

import "time"

type DemoRegisterRequest struct {
	FullName     string
	BusinessName string
	Email        string
	Password     string
	Phone        string
	Channel      string
}

type DemoRegisterResponse struct {
	Success bool
	Message string
}

type VerifyEmailRequest struct {
	Token string
}

type VerifyEmailResponse struct {
	Success bool
	Message string
}

type DemoVerifyOTPRequest struct {
	Email string
	Code  string
}

type DemoVerifyOTPResponse struct {
	Success bool
	Message string
}

type CreateDemoAccountParams struct {
	FullName     string
	BusinessName string
	BusinessCode string
	OrderPrefix  string
	Email        string
	Phone        string
	PasswordHash string
	RoleID       uint
	TokenHash    string
	ExpiresAt    time.Time
}

type DemoOTPEvent struct {
	Phone          string
	Code           string
	UserName       string
	ExpiresMinutes int
}

type EmailVerificationTokenInfo struct {
	ID        uint
	UserID    uint
	ExpiresAt time.Time
	UsedAt    *time.Time
}
