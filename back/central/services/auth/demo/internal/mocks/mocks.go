package mocks

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	GetDemoUserByEmailFn             func(ctx context.Context, email string) (*domain.PendingDemoUser, error)
	InvalidateEmailVerificationFn    func(ctx context.Context, userID uint) error
	CreateEmailVerificationTokenFn   func(ctx context.Context, userID uint, tokenHash string, expiresAt time.Time) error
	UpdateUserPhoneFn                func(ctx context.Context, userID uint, phone string) error
	BusinessCodeExistsFn             func(ctx context.Context, code string) (bool, error)
	GetDemoRoleIDFn                  func(ctx context.Context) (uint, error)
	CreateDemoAccountFn              func(ctx context.Context, params domain.CreateDemoAccountParams) (uint, error)
	GetValidEmailVerificationTokenFn func(ctx context.Context, tokenHash string) (*domain.EmailVerificationTokenInfo, error)
	ActivateUserAndConsumeTokenFn    func(ctx context.Context, tokenID, userID uint) error
	GetBusinessIDByUserIDFn          func(ctx context.Context, userID uint) (uint, error)
	ProvisionDemoIntegrationsFn      func(ctx context.Context, businessID, userID uint) error

	CuentasCreadas     []domain.CreateDemoAccountParams
	CodigosVerificados []string
	TokensCreados      []TokenCreado
	TokensInvalidados  []uint
	Activaciones       []Activacion
	Aprovisionamientos []Aprovisionamiento
	TelefonosCambiados []TelefonoCambiado
}

type TokenCreado struct {
	UserID    uint
	TokenHash string
	ExpiresAt time.Time
}

type Activacion struct {
	TokenID uint
	UserID  uint
}

type Aprovisionamiento struct {
	BusinessID uint
	UserID     uint
}

type TelefonoCambiado struct {
	UserID uint
	Phone  string
}

var _ domain.IDemoRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) GetDemoUserByEmail(ctx context.Context, email string) (*domain.PendingDemoUser, error) {
	if m.GetDemoUserByEmailFn != nil {
		return m.GetDemoUserByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *RepositoryMock) InvalidateEmailVerificationTokens(ctx context.Context, userID uint) error {
	m.TokensInvalidados = append(m.TokensInvalidados, userID)
	if m.InvalidateEmailVerificationFn != nil {
		return m.InvalidateEmailVerificationFn(ctx, userID)
	}
	return nil
}

func (m *RepositoryMock) CreateEmailVerificationToken(ctx context.Context, userID uint, tokenHash string, expiresAt time.Time) error {
	m.TokensCreados = append(m.TokensCreados, TokenCreado{UserID: userID, TokenHash: tokenHash, ExpiresAt: expiresAt})
	if m.CreateEmailVerificationTokenFn != nil {
		return m.CreateEmailVerificationTokenFn(ctx, userID, tokenHash, expiresAt)
	}
	return nil
}

func (m *RepositoryMock) UpdateUserPhone(ctx context.Context, userID uint, phone string) error {
	m.TelefonosCambiados = append(m.TelefonosCambiados, TelefonoCambiado{UserID: userID, Phone: phone})
	if m.UpdateUserPhoneFn != nil {
		return m.UpdateUserPhoneFn(ctx, userID, phone)
	}
	return nil
}

func (m *RepositoryMock) BusinessCodeExists(ctx context.Context, code string) (bool, error) {
	m.CodigosVerificados = append(m.CodigosVerificados, code)
	if m.BusinessCodeExistsFn != nil {
		return m.BusinessCodeExistsFn(ctx, code)
	}
	return false, nil
}

func (m *RepositoryMock) GetDemoRoleID(ctx context.Context) (uint, error) {
	if m.GetDemoRoleIDFn != nil {
		return m.GetDemoRoleIDFn(ctx)
	}
	return 9, nil
}

func (m *RepositoryMock) CreateDemoAccount(ctx context.Context, params domain.CreateDemoAccountParams) (uint, error) {
	m.CuentasCreadas = append(m.CuentasCreadas, params)
	if m.CreateDemoAccountFn != nil {
		return m.CreateDemoAccountFn(ctx, params)
	}
	return uint(len(m.CuentasCreadas)), nil
}

func (m *RepositoryMock) GetValidEmailVerificationToken(ctx context.Context, tokenHash string) (*domain.EmailVerificationTokenInfo, error) {
	if m.GetValidEmailVerificationTokenFn != nil {
		return m.GetValidEmailVerificationTokenFn(ctx, tokenHash)
	}
	return nil, nil
}

func (m *RepositoryMock) ActivateUserAndConsumeToken(ctx context.Context, tokenID, userID uint) error {
	m.Activaciones = append(m.Activaciones, Activacion{TokenID: tokenID, UserID: userID})
	if m.ActivateUserAndConsumeTokenFn != nil {
		return m.ActivateUserAndConsumeTokenFn(ctx, tokenID, userID)
	}
	return nil
}

func (m *RepositoryMock) GetBusinessIDByUserID(ctx context.Context, userID uint) (uint, error) {
	if m.GetBusinessIDByUserIDFn != nil {
		return m.GetBusinessIDByUserIDFn(ctx, userID)
	}
	return 100, nil
}

func (m *RepositoryMock) ProvisionDemoIntegrations(ctx context.Context, businessID, userID uint) error {
	m.Aprovisionamientos = append(m.Aprovisionamientos, Aprovisionamiento{BusinessID: businessID, UserID: userID})
	if m.ProvisionDemoIntegrationsFn != nil {
		return m.ProvisionDemoIntegrationsFn(ctx, businessID, userID)
	}
	return nil
}

type EmailSenderMock struct {
	SendHTMLFn func(ctx context.Context, to, subject, html string) error

	Enviados []CorreoEnviado
}

type CorreoEnviado struct {
	To      string
	Subject string
	HTML    string
}

var _ domain.IEmailSender = (*EmailSenderMock)(nil)

func (m *EmailSenderMock) SendHTML(ctx context.Context, to, subject, html string) error {
	m.Enviados = append(m.Enviados, CorreoEnviado{To: to, Subject: subject, HTML: html})
	if m.SendHTMLFn != nil {
		return m.SendHTMLFn(ctx, to, subject, html)
	}
	return nil
}

type OTPPublisherMock struct {
	PublishDemoOTPFn func(ctx context.Context, event domain.DemoOTPEvent) error

	Publicados []domain.DemoOTPEvent
}

var _ domain.IDemoOTPPublisher = (*OTPPublisherMock)(nil)

func (m *OTPPublisherMock) PublishDemoOTP(ctx context.Context, event domain.DemoOTPEvent) error {
	m.Publicados = append(m.Publicados, event)
	if m.PublishDemoOTPFn != nil {
		return m.PublishDemoOTPFn(ctx, event)
	}
	return nil
}

type ConfigMock struct {
	Values map[string]string
}

var _ env.IConfig = (*ConfigMock)(nil)

func (m *ConfigMock) Get(key string) string {
	return m.Values[key]
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
