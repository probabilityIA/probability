package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)


type EmailClientMock struct {
	SendHTMLFn func(ctx context.Context, to, subject, html string) error
	Calls      []EmailClientCall
}

type EmailClientCall struct {
	To      string
	Subject string
	HTML    string
}

func (m *EmailClientMock) SendHTML(ctx context.Context, to, subject, html string) error {
	m.Calls = append(m.Calls, EmailClientCall{To: to, Subject: subject, HTML: html})
	if m.SendHTMLFn != nil {
		return m.SendHTMLFn(ctx, to, subject, html)
	}
	return nil
}


type ResultPublisherMock struct {
	PublishResultFn func(ctx context.Context, result *entities.DeliveryResult) error
	Results         []*entities.DeliveryResult
}

func (m *ResultPublisherMock) PublishResult(ctx context.Context, result *entities.DeliveryResult) error {
	m.Results = append(m.Results, result)
	if m.PublishResultFn != nil {
		return m.PublishResultFn(ctx, result)
	}
	return nil
}


type UseCaseMock struct {
	SendNotificationEmailFn func(ctx context.Context, dto dtos.SendEmailDTO) error
	Calls                   []dtos.SendEmailDTO
}

func (m *UseCaseMock) SendNotificationEmail(ctx context.Context, dto dtos.SendEmailDTO) error {
	m.Calls = append(m.Calls, dto)
	if m.SendNotificationEmailFn != nil {
		return m.SendNotificationEmailFn(ctx, dto)
	}
	return nil
}


type RabbitMQMock struct {
	PublishFn func(ctx context.Context, queue string, body []byte) error
	ConsumeFn func(ctx context.Context, queue string, handler func([]byte) error) error
	Published []PublishedMessage
}

type PublishedMessage struct {
	Queue string
	Body  []byte
}

func (m *RabbitMQMock) Publish(ctx context.Context, queue string, body []byte) error {
	m.Published = append(m.Published, PublishedMessage{Queue: queue, Body: body})
	if m.PublishFn != nil {
		return m.PublishFn(ctx, queue, body)
	}
	return nil
}

func (m *RabbitMQMock) Consume(ctx context.Context, queue string, handler func([]byte) error) error {
	if m.ConsumeFn != nil {
		return m.ConsumeFn(ctx, queue, handler)
	}
	return nil
}

func (m *RabbitMQMock) DeclareExchange(name, kind string, durable bool) error { return nil }
func (m *RabbitMQMock) DeclareQueue(name string, durable bool) error           { return nil }
func (m *RabbitMQMock) BindQueue(queue, exchange, routingKey string) error      { return nil }
func (m *RabbitMQMock) PublishToExchange(ctx context.Context, exchange, routingKey string, body []byte) error {
	return nil
}
func (m *RabbitMQMock) Close() error { return nil }
func (m *RabbitMQMock) Ping() error  { return nil }


type LoggerMock struct {
	nop zerolog.Logger
}

func NewLoggerMock() log.ILogger {
	nop := zerolog.Nop()
	return &LoggerMock{nop: nop}
}

func (m *LoggerMock) Info(ctx ...context.Context) *zerolog.Event  { return m.nop.Info() }
func (m *LoggerMock) Error(ctx ...context.Context) *zerolog.Event { return m.nop.Error() }
func (m *LoggerMock) Warn(ctx ...context.Context) *zerolog.Event  { return m.nop.Warn() }
func (m *LoggerMock) Debug(ctx ...context.Context) *zerolog.Event { return m.nop.Debug() }
func (m *LoggerMock) Fatal(ctx ...context.Context) *zerolog.Event { return m.nop.Fatal() }
func (m *LoggerMock) Panic(ctx ...context.Context) *zerolog.Event { return m.nop.Panic() }
func (m *LoggerMock) With() zerolog.Context                       { return m.nop.With() }
func (m *LoggerMock) WithService(service string) log.ILogger      { return m }
func (m *LoggerMock) WithModule(module string) log.ILogger        { return m }
func (m *LoggerMock) WithBusinessID(businessID uint) log.ILogger  { return m }
