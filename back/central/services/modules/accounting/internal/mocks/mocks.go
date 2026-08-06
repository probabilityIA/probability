package mocks

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type RepositoryMock struct {
	ListConceptsFn      func(ctx context.Context) ([]entities.Concept, error)
	GetConceptByIDFn    func(ctx context.Context, id uint) (*entities.Concept, error)
	ConceptCodeExistsFn func(ctx context.Context, code string, excludeID *uint) (bool, error)
	CreateConceptFn     func(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error)
	UpdateConceptFn     func(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error)

	ListTaxesFn     func(ctx context.Context) ([]entities.Tax, error)
	GetTaxByIDFn    func(ctx context.Context, id uint) (*entities.Tax, error)
	TaxCodeExistsFn func(ctx context.Context, code string, excludeID *uint) (bool, error)
	CreateTaxFn     func(ctx context.Context, dto dtos.CreateTaxDTO) (*entities.Tax, error)
	UpdateTaxFn     func(ctx context.Context, dto dtos.UpdateTaxDTO) (*entities.Tax, error)

	SetConceptTaxFn func(ctx context.Context, dto dtos.SetConceptTaxDTO) error

	ListEntriesFn  func(ctx context.Context, params dtos.ListEntriesParams) ([]entities.Entry, int64, error)
	GetEntryByIDFn func(ctx context.Context, id uint) (*entities.Entry, error)
	CreateEntryFn  func(ctx context.Context, e *entities.Entry) (bool, error)
	DeleteEntryFn  func(ctx context.Context, id uint) error

	ReportFn func(ctx context.Context, params dtos.ReportParams) ([]dtos.ReportConceptRow, []dtos.ReportTaxRow, error)

	FindSyncCandidatesFn func(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error)

	CreateInvoiceFn       func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error)
	UpdateInvoiceFn       func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error)
	GetInvoiceByIDFn      func(ctx context.Context, id uint) (*entities.Invoice, error)
	ListInvoicesFn        func(ctx context.Context, params dtos.ListInvoicesParams) ([]entities.Invoice, int64, error)
	DeleteInvoiceFn       func(ctx context.Context, id uint) error
	SetInvoiceStatusFn    func(ctx context.Context, id uint, status string, sentAt, paidAt *time.Time, emailTo string) error
	SetInvoiceCustomerFn  func(ctx context.Context, id uint, document, phone, address string) error
	SetInvoiceDianFn      func(ctx context.Context, id uint, result *entities.DianEmitResult) error
	DeleteEntryBySourceFn func(ctx context.Context, sourceType, sourceID string) error
	GetTaxesByIDsFn       func(ctx context.Context, ids []uint) ([]entities.Tax, error)
	GetBusinessNameFn     func(ctx context.Context, id uint) (string, error)
	GetBusinessDefEmailFn func(ctx context.Context, id uint) (string, error)
	GetClientProfileFn    func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error)

	ListServicesFn      func(ctx context.Context) ([]entities.Service, error)
	GetServiceByIDFn    func(ctx context.Context, id uint) (*entities.Service, error)
	ServiceCodeExistsFn func(ctx context.Context, code string, excludeID *uint) (bool, error)
	CreateServiceFn     func(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error)
	UpdateServiceFn     func(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error)
	DeleteServiceFn     func(ctx context.Context, id uint) error

	CreatedEntries      []*entities.Entry
	DeletedEntrySources []string
	StatusCalls         []StatusCall
}

type StatusCall struct {
	ID      uint
	Status  string
	SentAt  *time.Time
	PaidAt  *time.Time
	EmailTo string
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) ListConcepts(ctx context.Context) ([]entities.Concept, error) {
	if m.ListConceptsFn != nil {
		return m.ListConceptsFn(ctx)
	}
	return nil, nil
}

func (m *RepositoryMock) GetConceptByID(ctx context.Context, id uint) (*entities.Concept, error) {
	if m.GetConceptByIDFn != nil {
		return m.GetConceptByIDFn(ctx, id)
	}
	return &entities.Concept{ID: id, Kind: entities.KindIncome, IsActive: true}, nil
}

func (m *RepositoryMock) ConceptCodeExists(ctx context.Context, code string, excludeID *uint) (bool, error) {
	if m.ConceptCodeExistsFn != nil {
		return m.ConceptCodeExistsFn(ctx, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) CreateConcept(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error) {
	if m.CreateConceptFn != nil {
		return m.CreateConceptFn(ctx, dto)
	}
	return &entities.Concept{ID: 1, Code: dto.Code, Name: dto.Name, Kind: dto.Kind}, nil
}

func (m *RepositoryMock) UpdateConcept(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error) {
	if m.UpdateConceptFn != nil {
		return m.UpdateConceptFn(ctx, dto)
	}
	return &entities.Concept{ID: dto.ID, Name: dto.Name}, nil
}

func (m *RepositoryMock) ListTaxes(ctx context.Context) ([]entities.Tax, error) {
	if m.ListTaxesFn != nil {
		return m.ListTaxesFn(ctx)
	}
	return nil, nil
}

func (m *RepositoryMock) GetTaxByID(ctx context.Context, id uint) (*entities.Tax, error) {
	if m.GetTaxByIDFn != nil {
		return m.GetTaxByIDFn(ctx, id)
	}
	return &entities.Tax{ID: id, IsActive: true}, nil
}

func (m *RepositoryMock) TaxCodeExists(ctx context.Context, code string, excludeID *uint) (bool, error) {
	if m.TaxCodeExistsFn != nil {
		return m.TaxCodeExistsFn(ctx, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) CreateTax(ctx context.Context, dto dtos.CreateTaxDTO) (*entities.Tax, error) {
	if m.CreateTaxFn != nil {
		return m.CreateTaxFn(ctx, dto)
	}
	return &entities.Tax{ID: 1, Code: dto.Code, Name: dto.Name, RatePercent: dto.RatePercent}, nil
}

func (m *RepositoryMock) UpdateTax(ctx context.Context, dto dtos.UpdateTaxDTO) (*entities.Tax, error) {
	if m.UpdateTaxFn != nil {
		return m.UpdateTaxFn(ctx, dto)
	}
	return &entities.Tax{ID: dto.ID, Name: dto.Name, RatePercent: dto.RatePercent}, nil
}

func (m *RepositoryMock) SetConceptTax(ctx context.Context, dto dtos.SetConceptTaxDTO) error {
	if m.SetConceptTaxFn != nil {
		return m.SetConceptTaxFn(ctx, dto)
	}
	return nil
}

func (m *RepositoryMock) ListEntries(ctx context.Context, params dtos.ListEntriesParams) ([]entities.Entry, int64, error) {
	if m.ListEntriesFn != nil {
		return m.ListEntriesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetEntryByID(ctx context.Context, id uint) (*entities.Entry, error) {
	if m.GetEntryByIDFn != nil {
		return m.GetEntryByIDFn(ctx, id)
	}
	return &entities.Entry{ID: id, SourceType: entities.SourceManual}, nil
}

func (m *RepositoryMock) CreateEntry(ctx context.Context, e *entities.Entry) (bool, error) {
	m.CreatedEntries = append(m.CreatedEntries, e)
	if m.CreateEntryFn != nil {
		return m.CreateEntryFn(ctx, e)
	}
	return true, nil
}

func (m *RepositoryMock) DeleteEntry(ctx context.Context, id uint) error {
	if m.DeleteEntryFn != nil {
		return m.DeleteEntryFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) Report(ctx context.Context, params dtos.ReportParams) ([]dtos.ReportConceptRow, []dtos.ReportTaxRow, error) {
	if m.ReportFn != nil {
		return m.ReportFn(ctx, params)
	}
	return nil, nil, nil
}

func (m *RepositoryMock) FindSyncCandidates(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
	if m.FindSyncCandidatesFn != nil {
		return m.FindSyncCandidatesFn(ctx, sourceType, limit)
	}
	return nil, nil
}

func (m *RepositoryMock) CreateInvoice(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
	if m.CreateInvoiceFn != nil {
		return m.CreateInvoiceFn(ctx, inv)
	}
	inv.ID = 1
	return inv, nil
}

func (m *RepositoryMock) UpdateInvoice(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
	if m.UpdateInvoiceFn != nil {
		return m.UpdateInvoiceFn(ctx, inv)
	}
	return inv, nil
}

func (m *RepositoryMock) GetInvoiceByID(ctx context.Context, id uint) (*entities.Invoice, error) {
	if m.GetInvoiceByIDFn != nil {
		return m.GetInvoiceByIDFn(ctx, id)
	}
	return &entities.Invoice{ID: id, Status: entities.InvoiceStatusDraft}, nil
}

func (m *RepositoryMock) ListInvoices(ctx context.Context, params dtos.ListInvoicesParams) ([]entities.Invoice, int64, error) {
	if m.ListInvoicesFn != nil {
		return m.ListInvoicesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) DeleteInvoice(ctx context.Context, id uint) error {
	if m.DeleteInvoiceFn != nil {
		return m.DeleteInvoiceFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) SetInvoiceStatus(ctx context.Context, id uint, status string, sentAt, paidAt *time.Time, emailTo string) error {
	m.StatusCalls = append(m.StatusCalls, StatusCall{ID: id, Status: status, SentAt: sentAt, PaidAt: paidAt, EmailTo: emailTo})
	if m.SetInvoiceStatusFn != nil {
		return m.SetInvoiceStatusFn(ctx, id, status, sentAt, paidAt, emailTo)
	}
	return nil
}

func (m *RepositoryMock) SetInvoiceCustomer(ctx context.Context, id uint, document, phone, address string) error {
	if m.SetInvoiceCustomerFn != nil {
		return m.SetInvoiceCustomerFn(ctx, id, document, phone, address)
	}
	return nil
}

func (m *RepositoryMock) SetInvoiceDian(ctx context.Context, id uint, result *entities.DianEmitResult) error {
	if m.SetInvoiceDianFn != nil {
		return m.SetInvoiceDianFn(ctx, id, result)
	}
	return nil
}

func (m *RepositoryMock) DeleteEntryBySource(ctx context.Context, sourceType, sourceID string) error {
	m.DeletedEntrySources = append(m.DeletedEntrySources, sourceType+":"+sourceID)
	if m.DeleteEntryBySourceFn != nil {
		return m.DeleteEntryBySourceFn(ctx, sourceType, sourceID)
	}
	return nil
}

func (m *RepositoryMock) GetTaxesByIDs(ctx context.Context, ids []uint) ([]entities.Tax, error) {
	if m.GetTaxesByIDsFn != nil {
		return m.GetTaxesByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (m *RepositoryMock) GetBusinessName(ctx context.Context, id uint) (string, error) {
	if m.GetBusinessNameFn != nil {
		return m.GetBusinessNameFn(ctx, id)
	}
	return "Demo", nil
}

func (m *RepositoryMock) GetBusinessDefaultEmail(ctx context.Context, id uint) (string, error) {
	if m.GetBusinessDefEmailFn != nil {
		return m.GetBusinessDefEmailFn(ctx, id)
	}
	return "", nil
}

func (m *RepositoryMock) GetClientProfile(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
	if m.GetClientProfileFn != nil {
		return m.GetClientProfileFn(ctx, businessID)
	}
	return &dtos.ClientProfileDTO{}, nil
}

func (m *RepositoryMock) ListServices(ctx context.Context) ([]entities.Service, error) {
	if m.ListServicesFn != nil {
		return m.ListServicesFn(ctx)
	}
	return nil, nil
}

func (m *RepositoryMock) GetServiceByID(ctx context.Context, id uint) (*entities.Service, error) {
	if m.GetServiceByIDFn != nil {
		return m.GetServiceByIDFn(ctx, id)
	}
	return &entities.Service{ID: id}, nil
}

func (m *RepositoryMock) ServiceCodeExists(ctx context.Context, code string, excludeID *uint) (bool, error) {
	if m.ServiceCodeExistsFn != nil {
		return m.ServiceCodeExistsFn(ctx, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) CreateService(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error) {
	if m.CreateServiceFn != nil {
		return m.CreateServiceFn(ctx, dto)
	}
	return &entities.Service{ID: 1, Code: dto.Code, Name: dto.Name, IsActive: dto.IsActive}, nil
}

func (m *RepositoryMock) UpdateService(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error) {
	if m.UpdateServiceFn != nil {
		return m.UpdateServiceFn(ctx, dto)
	}
	return &entities.Service{ID: dto.ID, Code: dto.Code, Name: dto.Name}, nil
}

func (m *RepositoryMock) DeleteService(ctx context.Context, id uint) error {
	if m.DeleteServiceFn != nil {
		return m.DeleteServiceFn(ctx, id)
	}
	return nil
}

type EmailCall struct {
	To      string
	Subject string
	HTML    string
}

type EmailSenderMock struct {
	SendHTMLFn func(ctx context.Context, to, subject, html string) error
	Sent       []EmailCall
}

var _ ports.IEmailSender = (*EmailSenderMock)(nil)

func (m *EmailSenderMock) SendHTML(ctx context.Context, to, subject, html string) error {
	m.Sent = append(m.Sent, EmailCall{To: to, Subject: subject, HTML: html})
	if m.SendHTMLFn != nil {
		return m.SendHTMLFn(ctx, to, subject, html)
	}
	return nil
}

type DianEmitterMock struct {
	EmitFn      func(ctx context.Context, integrationID uint, inv *entities.Invoice, config map[string]interface{}) (*entities.DianEmitResult, error)
	EmitConfigs []map[string]interface{}
}

var _ ports.IDianEmitter = (*DianEmitterMock)(nil)

func (m *DianEmitterMock) Emit(ctx context.Context, integrationID uint, inv *entities.Invoice, config map[string]interface{}) (*entities.DianEmitResult, error) {
	m.EmitConfigs = append(m.EmitConfigs, config)
	if m.EmitFn != nil {
		return m.EmitFn(ctx, integrationID, inv, config)
	}
	return &entities.DianEmitResult{}, nil
}

type DianIntegrationMock struct {
	GetGlobalIntegrationIDFn func(ctx context.Context) (uint, error)
	GetConfigFn              func(ctx context.Context, integrationID uint) (map[string]interface{}, error)
	EnsureFn                 func(ctx context.Context, name string, config, credentials map[string]interface{}) (uint, error)
	TestConnectionFn         func(ctx context.Context, config, credentials map[string]interface{}) error

	EnsuredConfigs     []map[string]interface{}
	EnsuredCredentials []map[string]interface{}
	TestedCredentials  []map[string]interface{}
}

var _ ports.IDianIntegration = (*DianIntegrationMock)(nil)

func (m *DianIntegrationMock) GetGlobalIntegrationID(ctx context.Context) (uint, error) {
	if m.GetGlobalIntegrationIDFn != nil {
		return m.GetGlobalIntegrationIDFn(ctx)
	}
	return 1, nil
}

func (m *DianIntegrationMock) GetConfig(ctx context.Context, integrationID uint) (map[string]interface{}, error) {
	if m.GetConfigFn != nil {
		return m.GetConfigFn(ctx, integrationID)
	}
	return map[string]interface{}{}, nil
}

func (m *DianIntegrationMock) Ensure(ctx context.Context, name string, config, credentials map[string]interface{}) (uint, error) {
	m.EnsuredConfigs = append(m.EnsuredConfigs, config)
	m.EnsuredCredentials = append(m.EnsuredCredentials, credentials)
	if m.EnsureFn != nil {
		return m.EnsureFn(ctx, name, config, credentials)
	}
	return 1, nil
}

func (m *DianIntegrationMock) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	m.TestedCredentials = append(m.TestedCredentials, credentials)
	if m.TestConnectionFn != nil {
		return m.TestConnectionFn(ctx, config, credentials)
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
