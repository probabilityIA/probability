package mocks

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type AuthRepositoryMock struct {
	GetUserByEmailFn                     func(ctx context.Context, email string) (*domain.UserAuthInfo, error)
	HasPendingEmailVerificationFn        func(ctx context.Context, userID uint) (bool, error)
	GetUserByIDFn                        func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error)
	CreatePasswordResetTokenFn           func(ctx context.Context, userID uint, tokenHash string, channel string, expiresAt time.Time) error
	InvalidateUserPasswordResetTokensFn  func(ctx context.Context, userID uint) error
	GetValidPasswordResetTokenFn         func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error)
	GetActiveOTPTokenFn                  func(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error)
	IncrementPasswordResetTokenAttemptsFn func(ctx context.Context, tokenID uint) error
	MarkPasswordResetTokenUsedFn         func(ctx context.Context, tokenID uint) error
	GetUserRolesFn                       func(ctx context.Context, userID uint) ([]domain.Role, error)
	GetRolePermissionsFn                 func(ctx context.Context, roleID uint) ([]domain.Permission, error)
	UpdateLastLoginFn                    func(ctx context.Context, userID uint) error
	ChangePasswordFn                     func(ctx context.Context, userID uint, newPassword string) error
	GetUserBusinessesFn                  func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error)
	GetUserRoleByBusinessFn              func(ctx context.Context, userID uint, businessID uint) (*domain.Role, error)
	GetBusinessStaffRelationFn           func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error)
	GetBusinessConfiguredResourcesIDsFn  func(ctx context.Context, businessID uint) ([]uint, error)
	GetBusinessByIDFn                    func(ctx context.Context, businessID uint) (*domain.BusinessInfo, error)
	GetRoleByIDFn                        func(ctx context.Context, id uint) (*domain.Role, error)
}

var _ domain.IAuthRepository = (*AuthRepositoryMock)(nil)

func (m *AuthRepositoryMock) GetUserByEmail(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
	if m.GetUserByEmailFn != nil {
		return m.GetUserByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) HasPendingEmailVerification(ctx context.Context, userID uint) (bool, error) {
	if m.HasPendingEmailVerificationFn != nil {
		return m.HasPendingEmailVerificationFn(ctx, userID)
	}
	return false, nil
}

func (m *AuthRepositoryMock) GetUserByID(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
	if m.GetUserByIDFn != nil {
		return m.GetUserByIDFn(ctx, userID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) CreatePasswordResetToken(ctx context.Context, userID uint, tokenHash string, channel string, expiresAt time.Time) error {
	if m.CreatePasswordResetTokenFn != nil {
		return m.CreatePasswordResetTokenFn(ctx, userID, tokenHash, channel, expiresAt)
	}
	return nil
}

func (m *AuthRepositoryMock) InvalidateUserPasswordResetTokens(ctx context.Context, userID uint) error {
	if m.InvalidateUserPasswordResetTokensFn != nil {
		return m.InvalidateUserPasswordResetTokensFn(ctx, userID)
	}
	return nil
}

func (m *AuthRepositoryMock) GetValidPasswordResetToken(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
	if m.GetValidPasswordResetTokenFn != nil {
		return m.GetValidPasswordResetTokenFn(ctx, tokenHash)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) GetActiveOTPToken(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error) {
	if m.GetActiveOTPTokenFn != nil {
		return m.GetActiveOTPTokenFn(ctx, userID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) IncrementPasswordResetTokenAttempts(ctx context.Context, tokenID uint) error {
	if m.IncrementPasswordResetTokenAttemptsFn != nil {
		return m.IncrementPasswordResetTokenAttemptsFn(ctx, tokenID)
	}
	return nil
}

func (m *AuthRepositoryMock) MarkPasswordResetTokenUsed(ctx context.Context, tokenID uint) error {
	if m.MarkPasswordResetTokenUsedFn != nil {
		return m.MarkPasswordResetTokenUsedFn(ctx, tokenID)
	}
	return nil
}

func (m *AuthRepositoryMock) GetUserRoles(ctx context.Context, userID uint) ([]domain.Role, error) {
	if m.GetUserRolesFn != nil {
		return m.GetUserRolesFn(ctx, userID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) GetRolePermissions(ctx context.Context, roleID uint) ([]domain.Permission, error) {
	if m.GetRolePermissionsFn != nil {
		return m.GetRolePermissionsFn(ctx, roleID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) UpdateLastLogin(ctx context.Context, userID uint) error {
	if m.UpdateLastLoginFn != nil {
		return m.UpdateLastLoginFn(ctx, userID)
	}
	return nil
}

func (m *AuthRepositoryMock) ChangePassword(ctx context.Context, userID uint, newPassword string) error {
	if m.ChangePasswordFn != nil {
		return m.ChangePasswordFn(ctx, userID, newPassword)
	}
	return nil
}

func (m *AuthRepositoryMock) GetUserBusinesses(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
	if m.GetUserBusinessesFn != nil {
		return m.GetUserBusinessesFn(ctx, userID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) GetUserRoleByBusiness(ctx context.Context, userID uint, businessID uint) (*domain.Role, error) {
	if m.GetUserRoleByBusinessFn != nil {
		return m.GetUserRoleByBusinessFn(ctx, userID, businessID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) GetBusinessStaffRelation(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
	if m.GetBusinessStaffRelationFn != nil {
		return m.GetBusinessStaffRelationFn(ctx, userID, businessID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) GetBusinessConfiguredResourcesIDs(ctx context.Context, businessID uint) ([]uint, error) {
	if m.GetBusinessConfiguredResourcesIDsFn != nil {
		return m.GetBusinessConfiguredResourcesIDsFn(ctx, businessID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) GetBusinessByID(ctx context.Context, businessID uint) (*domain.BusinessInfo, error) {
	if m.GetBusinessByIDFn != nil {
		return m.GetBusinessByIDFn(ctx, businessID)
	}
	return nil, nil
}

func (m *AuthRepositoryMock) GetRoleByID(ctx context.Context, id uint) (*domain.Role, error) {
	if m.GetRoleByIDFn != nil {
		return m.GetRoleByIDFn(ctx, id)
	}
	return nil, nil
}

type JWTServiceMock struct {
	GenerateTokenFn func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error)
	ValidateTokenFn func(tokenString string) (*domain.JWTClaims, error)
	RefreshTokenFn  func(tokenString string) (string, error)
}

var _ domain.IJWTService = (*JWTServiceMock)(nil)

func (m *JWTServiceMock) GenerateToken(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
	if m.GenerateTokenFn != nil {
		return m.GenerateTokenFn(userID, businessID, businessTypeID, roleID, subscriptionStatus)
	}
	return "token-de-prueba", nil
}

func (m *JWTServiceMock) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	if m.ValidateTokenFn != nil {
		return m.ValidateTokenFn(tokenString)
	}
	return nil, nil
}

func (m *JWTServiceMock) RefreshToken(tokenString string) (string, error) {
	if m.RefreshTokenFn != nil {
		return m.RefreshTokenFn(tokenString)
	}
	return "", nil
}

type EmailSenderMock struct {
	SendHTMLFn func(ctx context.Context, to, subject, html string) error
}

var _ domain.IEmailSender = (*EmailSenderMock)(nil)

func (m *EmailSenderMock) SendHTML(ctx context.Context, to, subject, html string) error {
	if m.SendHTMLFn != nil {
		return m.SendHTMLFn(ctx, to, subject, html)
	}
	return nil
}

type OTPEventPublisherMock struct {
	PublishPasswordResetOTPFn func(ctx context.Context, event domain.PasswordResetOTPEvent) error
}

var _ domain.IOTPEventPublisher = (*OTPEventPublisherMock)(nil)

func (m *OTPEventPublisherMock) PublishPasswordResetOTP(ctx context.Context, event domain.PasswordResetOTPEvent) error {
	if m.PublishPasswordResetOTPFn != nil {
		return m.PublishPasswordResetOTPFn(ctx, event)
	}
	return nil
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger {
	return &SilentLogger{}
}

func (l *SilentLogger) nop() zerolog.Logger {
	return zerolog.Nop()
}

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Info()
}

func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Error()
}

func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Warn()
}

func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Debug()
}

func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Fatal()
}

func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Panic()
}

func (l *SilentLogger) With() zerolog.Context {
	n := l.nop()
	return n.With()
}

func (l *SilentLogger) WithService(service string) log.ILogger {
	return l
}

func (l *SilentLogger) WithModule(module string) log.ILogger {
	return l
}

func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger {
	return l
}
