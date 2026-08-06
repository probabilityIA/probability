package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/pay/epayco/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IntegrationRepositoryMock struct {
	GetEPaycoConfigFn func(ctx context.Context) (*ports.EPaycoConfig, error)

	Consultas int
}

var _ ports.IIntegrationRepository = (*IntegrationRepositoryMock)(nil)

func (m *IntegrationRepositoryMock) GetEPaycoConfig(ctx context.Context) (*ports.EPaycoConfig, error) {
	m.Consultas++
	if m.GetEPaycoConfigFn != nil {
		return m.GetEPaycoConfigFn(ctx)
	}
	return &ports.EPaycoConfig{Environment: "sandbox"}, nil
}

type ClientMock struct {
	CreateCheckoutFn func(ctx context.Context, config *ports.EPaycoConfig, amount float64, currency, reference, description string) (string, string, error)

	Llamadas []Llamada
}

type Llamada struct {
	Amount      float64
	Currency    string
	Reference   string
	Description string
	Config      *ports.EPaycoConfig
}

var _ ports.IEPaycoClient = (*ClientMock)(nil)

func (m *ClientMock) CreateCheckout(ctx context.Context, config *ports.EPaycoConfig, amount float64, currency, reference, description string) (string, string, error) {
	m.Llamadas = append(m.Llamadas, Llamada{
		Amount: amount, Currency: currency, Reference: reference, Description: description, Config: config,
	})
	if m.CreateCheckoutFn != nil {
		return m.CreateCheckoutFn(ctx, config, amount, currency, reference, description)
	}
	return "CHK-1", "https://checkout.epayco.co/CHK-1", nil
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
