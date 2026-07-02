package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type SiigoClientMock struct {
	TestAuthenticationFn          func(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) error
	CreateInvoiceFn               func(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error)
	GetCustomerByIdentificationFn func(ctx context.Context, credentials dtos.Credentials, identification string) (*dtos.CustomerResult, error)
	CreateCustomerFn              func(ctx context.Context, credentials dtos.Credentials, req *dtos.CreateCustomerRequest) (*dtos.CustomerResult, error)
	ListInvoicesFn                func(ctx context.Context, credentials dtos.Credentials, params dtos.ListInvoicesParams) (*dtos.ListInvoicesResult, error)
	GetInvoiceByIDFn              func(ctx context.Context, credentials dtos.Credentials, invoiceID string) (*dtos.InvoiceDetail, error)
	GetStampErrorsFn              func(ctx context.Context, credentials dtos.Credentials, invoiceID string) ([]dtos.StampError, error)
	AnnulInvoiceFn                func(ctx context.Context, credentials dtos.Credentials, invoiceID string) (*dtos.AnnulInvoiceResult, error)
	ListProductsFn                func(ctx context.Context, credentials dtos.Credentials, page, pageSize int) ([]dtos.ProductItem, error)
	ListWarehousesFn              func(ctx context.Context, credentials dtos.Credentials) ([]dtos.WarehouseItem, error)
	ListPaymentTypesFn            func(ctx context.Context, credentials dtos.Credentials, documentType string) ([]dtos.PaymentTypeItem, error)
	CreateCashReceiptFn           func(ctx context.Context, req *dtos.CreateCashReceiptRequest) (*dtos.CreateCashReceiptResult, error)
	CreateCreditNoteFn            func(ctx context.Context, req *dtos.CreateCreditNoteRequest) (*dtos.CreateCreditNoteResult, error)
	CreateJournalFn               func(ctx context.Context, req *dtos.CreateJournalRequest) (*dtos.CreateJournalResult, error)
	ListWebhooksFn                func(ctx context.Context, credentials dtos.Credentials) ([]dtos.WebhookItem, error)
	CreateWebhookFn               func(ctx context.Context, credentials dtos.Credentials, input dtos.CreateWebhookInput) (*dtos.WebhookItem, error)
	DeleteWebhookFn               func(ctx context.Context, credentials dtos.Credentials, webhookID string) error
}

func (m *SiigoClientMock) ListWebhooks(ctx context.Context, credentials dtos.Credentials) ([]dtos.WebhookItem, error) {
	if m.ListWebhooksFn != nil {
		return m.ListWebhooksFn(ctx, credentials)
	}
	return nil, nil
}

func (m *SiigoClientMock) CreateWebhook(ctx context.Context, credentials dtos.Credentials, input dtos.CreateWebhookInput) (*dtos.WebhookItem, error) {
	if m.CreateWebhookFn != nil {
		return m.CreateWebhookFn(ctx, credentials, input)
	}
	return &dtos.WebhookItem{}, nil
}

func (m *SiigoClientMock) DeleteWebhook(ctx context.Context, credentials dtos.Credentials, webhookID string) error {
	if m.DeleteWebhookFn != nil {
		return m.DeleteWebhookFn(ctx, credentials, webhookID)
	}
	return nil
}

func (m *SiigoClientMock) CreateCreditNote(
	ctx context.Context,
	req *dtos.CreateCreditNoteRequest,
) (*dtos.CreateCreditNoteResult, error) {
	if m.CreateCreditNoteFn != nil {
		return m.CreateCreditNoteFn(ctx, req)
	}
	return &dtos.CreateCreditNoteResult{}, nil
}

func (m *SiigoClientMock) TestAuthentication(
	ctx context.Context,
	username, accessKey, accountID, partnerID, baseURL string,
) error {
	if m.TestAuthenticationFn != nil {
		return m.TestAuthenticationFn(ctx, username, accessKey, accountID, partnerID, baseURL)
	}
	return nil
}

func (m *SiigoClientMock) CreateInvoice(
	ctx context.Context,
	req *dtos.CreateInvoiceRequest,
) (*dtos.CreateInvoiceResult, error) {
	if m.CreateInvoiceFn != nil {
		return m.CreateInvoiceFn(ctx, req)
	}
	return &dtos.CreateInvoiceResult{}, nil
}

func (m *SiigoClientMock) GetCustomerByIdentification(
	ctx context.Context,
	credentials dtos.Credentials,
	identification string,
) (*dtos.CustomerResult, error) {
	if m.GetCustomerByIdentificationFn != nil {
		return m.GetCustomerByIdentificationFn(ctx, credentials, identification)
	}
	return nil, nil
}

func (m *SiigoClientMock) CreateCustomer(
	ctx context.Context,
	credentials dtos.Credentials,
	req *dtos.CreateCustomerRequest,
) (*dtos.CustomerResult, error) {
	if m.CreateCustomerFn != nil {
		return m.CreateCustomerFn(ctx, credentials, req)
	}
	return &dtos.CustomerResult{}, nil
}

func (m *SiigoClientMock) ListInvoices(
	ctx context.Context,
	credentials dtos.Credentials,
	params dtos.ListInvoicesParams,
) (*dtos.ListInvoicesResult, error) {
	if m.ListInvoicesFn != nil {
		return m.ListInvoicesFn(ctx, credentials, params)
	}
	return &dtos.ListInvoicesResult{}, nil
}

func (m *SiigoClientMock) GetInvoiceByID(
	ctx context.Context,
	credentials dtos.Credentials,
	invoiceID string,
) (*dtos.InvoiceDetail, error) {
	if m.GetInvoiceByIDFn != nil {
		return m.GetInvoiceByIDFn(ctx, credentials, invoiceID)
	}
	return &dtos.InvoiceDetail{}, nil
}

func (m *SiigoClientMock) GetStampErrors(
	ctx context.Context,
	credentials dtos.Credentials,
	invoiceID string,
) ([]dtos.StampError, error) {
	if m.GetStampErrorsFn != nil {
		return m.GetStampErrorsFn(ctx, credentials, invoiceID)
	}
	return nil, nil
}

func (m *SiigoClientMock) AnnulInvoice(
	ctx context.Context,
	credentials dtos.Credentials,
	invoiceID string,
) (*dtos.AnnulInvoiceResult, error) {
	if m.AnnulInvoiceFn != nil {
		return m.AnnulInvoiceFn(ctx, credentials, invoiceID)
	}
	return &dtos.AnnulInvoiceResult{}, nil
}

func (m *SiigoClientMock) ListProducts(
	ctx context.Context,
	credentials dtos.Credentials,
	page, pageSize int,
) ([]dtos.ProductItem, error) {
	if m.ListProductsFn != nil {
		return m.ListProductsFn(ctx, credentials, page, pageSize)
	}
	return nil, nil
}

func (m *SiigoClientMock) ListWarehouses(
	ctx context.Context,
	credentials dtos.Credentials,
) ([]dtos.WarehouseItem, error) {
	if m.ListWarehousesFn != nil {
		return m.ListWarehousesFn(ctx, credentials)
	}
	return nil, nil
}

func (m *SiigoClientMock) ListPaymentTypes(
	ctx context.Context,
	credentials dtos.Credentials,
	documentType string,
) ([]dtos.PaymentTypeItem, error) {
	if m.ListPaymentTypesFn != nil {
		return m.ListPaymentTypesFn(ctx, credentials, documentType)
	}
	return nil, nil
}

func (m *SiigoClientMock) CreateCashReceipt(
	ctx context.Context,
	req *dtos.CreateCashReceiptRequest,
) (*dtos.CreateCashReceiptResult, error) {
	if m.CreateCashReceiptFn != nil {
		return m.CreateCashReceiptFn(ctx, req)
	}
	return &dtos.CreateCashReceiptResult{}, nil
}

func (m *SiigoClientMock) CreateJournal(
	ctx context.Context,
	req *dtos.CreateJournalRequest,
) (*dtos.CreateJournalResult, error) {
	if m.CreateJournalFn != nil {
		return m.CreateJournalFn(ctx, req)
	}
	return &dtos.CreateJournalResult{}, nil
}
