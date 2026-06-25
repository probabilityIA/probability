package domain

import "context"

type IDemoRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	BusinessCodeExists(ctx context.Context, code string) (bool, error)
	GetDemoRoleID(ctx context.Context) (uint, error)
	CreateDemoAccount(ctx context.Context, params CreateDemoAccountParams) (uint, error)
	GetValidEmailVerificationToken(ctx context.Context, tokenHash string) (*EmailVerificationTokenInfo, error)
	ActivateUserAndConsumeToken(ctx context.Context, tokenID, userID uint) error
	GetBusinessIDByUserID(ctx context.Context, userID uint) (uint, error)
	ProvisionDemoIntegrations(ctx context.Context, businessID, userID uint) error
}

type IEmailSender interface {
	SendHTML(ctx context.Context, to, subject, html string) error
}
