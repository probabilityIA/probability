package domain

import "time"

type DemoRegisterRequest struct {
	FullName     string
	BusinessName string
	Email        string
	Password     string
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

type CreateDemoAccountParams struct {
	FullName     string
	BusinessName string
	BusinessCode string
	OrderPrefix  string
	Email        string
	PasswordHash string
	RoleID       uint
	TokenHash    string
	ExpiresAt    time.Time
}

type EmailVerificationTokenInfo struct {
	ID        uint
	UserID    uint
	ExpiresAt time.Time
	UsedAt    *time.Time
}
