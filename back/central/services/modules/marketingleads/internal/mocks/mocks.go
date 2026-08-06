package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	CreateLeadFn           func(ctx context.Context, lead *entities.MarketingLead) error
	SetWhatsAppMessageIDFn func(ctx context.Context, leadID uint, messageID string) error
	ListLeadsFn            func(ctx context.Context, page, pageSize int) ([]entities.MarketingLead, int64, error)

	Created      []entities.MarketingLead
	MessageIDSet []MessageIDCall
	ListCalls    []ListCall
}

type MessageIDCall struct {
	LeadID    uint
	MessageID string
}

type ListCall struct {
	Page     int
	PageSize int
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) CreateLead(ctx context.Context, lead *entities.MarketingLead) error {
	if m.CreateLeadFn != nil {
		if err := m.CreateLeadFn(ctx, lead); err != nil {
			return err
		}
	}
	if lead.ID == 0 {
		lead.ID = uint(len(m.Created) + 1)
	}
	m.Created = append(m.Created, *lead)
	return nil
}

func (m *RepositoryMock) SetWhatsAppMessageID(ctx context.Context, leadID uint, messageID string) error {
	m.MessageIDSet = append(m.MessageIDSet, MessageIDCall{LeadID: leadID, MessageID: messageID})
	if m.SetWhatsAppMessageIDFn != nil {
		return m.SetWhatsAppMessageIDFn(ctx, leadID, messageID)
	}
	return nil
}

func (m *RepositoryMock) ListLeads(ctx context.Context, page, pageSize int) ([]entities.MarketingLead, int64, error) {
	m.ListCalls = append(m.ListCalls, ListCall{Page: page, PageSize: pageSize})
	if m.ListLeadsFn != nil {
		return m.ListLeadsFn(ctx, page, pageSize)
	}
	return []entities.MarketingLead{}, 0, nil
}

type WhatsAppSenderMock struct {
	SendFn func(ctx context.Context, phone, name, surveyTitle, level string) (string, error)

	Calls []WhatsAppCall
}

type WhatsAppCall struct {
	Phone       string
	Name        string
	SurveyTitle string
	Level       string
}

var _ ports.IWhatsAppSender = (*WhatsAppSenderMock)(nil)

func (m *WhatsAppSenderMock) SendDiagnosticResultTemplate(ctx context.Context, phone, name, surveyTitle, level string) (string, error) {
	m.Calls = append(m.Calls, WhatsAppCall{Phone: phone, Name: name, SurveyTitle: surveyTitle, Level: level})
	if m.SendFn != nil {
		return m.SendFn(ctx, phone, name, surveyTitle, level)
	}
	return "wamid.TEST", nil
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
