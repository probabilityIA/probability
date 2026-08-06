package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IntegrationRepositoryMock struct {
	GetNequiConfigFn func(ctx context.Context) (*ports.NequiConfig, error)

	Consultas int
}

var _ ports.IIntegrationRepository = (*IntegrationRepositoryMock)(nil)

func (m *IntegrationRepositoryMock) GetNequiConfig(ctx context.Context) (*ports.NequiConfig, error) {
	m.Consultas++
	if m.GetNequiConfigFn != nil {
		return m.GetNequiConfigFn(ctx)
	}
	return &ports.NequiConfig{APIKey: "api-key", Environment: "sandbox", PhoneCode: "NIT_1"}, nil
}

type ClientMock struct {
	GenerateQRFn func(ctx context.Context, config *ports.NequiConfig, amount float64, reference string) (string, string, error)

	Llamadas []Llamada
}

type Llamada struct {
	Amount    float64
	Reference string
	Config    *ports.NequiConfig
}

var _ ports.INequiClient = (*ClientMock)(nil)

func (m *ClientMock) GenerateQR(ctx context.Context, config *ports.NequiConfig, amount float64, reference string) (string, string, error) {
	m.Llamadas = append(m.Llamadas, Llamada{Amount: amount, Reference: reference, Config: config})
	if m.GenerateQRFn != nil {
		return m.GenerateQRFn(ctx, config, amount, reference)
	}
	return "00020101021243650016", "TX-1", nil
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
