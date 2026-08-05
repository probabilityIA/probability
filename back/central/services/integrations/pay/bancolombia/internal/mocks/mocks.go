package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IntegrationRepositoryMock struct {
	GetConfigFn        func(ctx context.Context) (*ports.BancolombiaConfig, error)
	GetConfigForModeFn func(ctx context.Context, testMode bool) (*ports.BancolombiaConfig, error)

	ModosConsultados []bool
}

var _ ports.IIntegrationRepository = (*IntegrationRepositoryMock)(nil)

func configPorDefecto() *ports.BancolombiaConfig {
	return &ports.BancolombiaConfig{
		ClientID:      "client-id",
		ClientSecret:  "client-secret",
		WebhookSecret: "secreto",
		Environment:   "sandbox",
		MerchantID:    "M-1",
	}
}

func (m *IntegrationRepositoryMock) GetConfig(ctx context.Context) (*ports.BancolombiaConfig, error) {
	if m.GetConfigFn != nil {
		return m.GetConfigFn(ctx)
	}
	return configPorDefecto(), nil
}

func (m *IntegrationRepositoryMock) GetConfigForMode(ctx context.Context, testMode bool) (*ports.BancolombiaConfig, error) {
	m.ModosConsultados = append(m.ModosConsultados, testMode)
	if m.GetConfigForModeFn != nil {
		return m.GetConfigForModeFn(ctx, testMode)
	}
	return configPorDefecto(), nil
}

type ClientMock struct {
	GenerateQRFn func(ctx context.Context, config *ports.BancolombiaConfig, amount float64, currency, reference, description string) (*ports.QRGenerationResult, error)

	Llamadas []LlamadaQR
}

type LlamadaQR struct {
	Amount      float64
	Currency    string
	Reference   string
	Description string
	Config      *ports.BancolombiaConfig
}

var _ ports.IBancolombiaClient = (*ClientMock)(nil)

func (m *ClientMock) GenerateQR(ctx context.Context, config *ports.BancolombiaConfig, amount float64, currency, reference, description string) (*ports.QRGenerationResult, error) {
	m.Llamadas = append(m.Llamadas, LlamadaQR{
		Amount: amount, Currency: currency, Reference: reference, Description: description, Config: config,
	})
	if m.GenerateQRFn != nil {
		return m.GenerateQRFn(ctx, config, amount, currency, reference, description)
	}
	return &ports.QRGenerationResult{
		ExternalID:    "QR-1",
		QRValue:       "00020101021243650016",
		QRImageBase64: "iVBORw0KGgo=",
	}, nil
}

type ResponsePublisherMock struct {
	PublishPaymentResponseFn func(ctx context.Context, msg *ports.PaymentResponseMsg) error

	Publicados []ports.PaymentResponseMsg
}

var _ ports.IResponsePublisher = (*ResponsePublisherMock)(nil)

func (m *ResponsePublisherMock) PublishPaymentResponse(ctx context.Context, msg *ports.PaymentResponseMsg) error {
	m.Publicados = append(m.Publicados, *msg)
	if m.PublishPaymentResponseFn != nil {
		return m.PublishPaymentResponseFn(ctx, msg)
	}
	return nil
}

type WebhookPublisherMock struct {
	PublishWebhookEventFn func(ctx context.Context, msg *ports.BancolombiaWebhookMessage) error

	Publicados []ports.BancolombiaWebhookMessage
}

var _ ports.IWebhookPublisher = (*WebhookPublisherMock)(nil)

func (m *WebhookPublisherMock) PublishWebhookEvent(ctx context.Context, msg *ports.BancolombiaWebhookMessage) error {
	m.Publicados = append(m.Publicados, *msg)
	if m.PublishWebhookEventFn != nil {
		return m.PublishWebhookEventFn(ctx, msg)
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
