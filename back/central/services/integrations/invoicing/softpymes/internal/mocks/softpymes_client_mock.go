package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
)

// SoftpymesClientMock implementa ports.ISoftpymesClient para tests unitarios.
// Cada método tiene un campo Fn inyectable para configurar el comportamiento.
type SoftpymesClientMock struct {
	TestAuthenticationFn  func(ctx context.Context, apiKey, apiSecret, referer string) error
	CreateInvoiceFn       func(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error)
	CreateCreditNoteFn    func(ctx context.Context, creditNoteData map[string]interface{}) error
	GetDocumentByNumberFn func(ctx context.Context, apiKey, apiSecret, referer, documentNumber string) (map[string]interface{}, error)
}

func (m *SoftpymesClientMock) TestAuthentication(ctx context.Context, apiKey, apiSecret, referer string) error {
	if m.TestAuthenticationFn != nil {
		return m.TestAuthenticationFn(ctx, apiKey, apiSecret, referer)
	}
	return nil
}

func (m *SoftpymesClientMock) CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
	if m.CreateInvoiceFn != nil {
		return m.CreateInvoiceFn(ctx, req)
	}
	return &dtos.CreateInvoiceResult{}, nil
}

func (m *SoftpymesClientMock) CreateCreditNote(ctx context.Context, creditNoteData map[string]interface{}) error {
	if m.CreateCreditNoteFn != nil {
		return m.CreateCreditNoteFn(ctx, creditNoteData)
	}
	return nil
}

func (m *SoftpymesClientMock) GetDocumentByNumber(ctx context.Context, apiKey, apiSecret, referer, documentNumber string) (map[string]interface{}, error) {
	if m.GetDocumentByNumberFn != nil {
		return m.GetDocumentByNumberFn(ctx, apiKey, apiSecret, referer, documentNumber)
	}
	return nil, nil
}
