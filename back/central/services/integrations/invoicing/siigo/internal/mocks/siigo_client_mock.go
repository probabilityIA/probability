package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

// SiigoClientMock es el mock de ports.ISiigoClient.
// Cada campo *Fn permite configurar el comportamiento por test.
type SiigoClientMock struct {
	TestAuthenticationFn          func(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) error
	CreateInvoiceFn               func(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error)
	GetCustomerByIdentificationFn func(ctx context.Context, credentials dtos.Credentials, identification string) (*dtos.CustomerResult, error)
	CreateCustomerFn              func(ctx context.Context, credentials dtos.Credentials, req *dtos.CreateCustomerRequest) (*dtos.CustomerResult, error)
	ListInvoicesFn                func(ctx context.Context, credentials dtos.Credentials, params dtos.ListInvoicesParams) (*dtos.ListInvoicesResult, error)
}

// TestAuthentication implementa ports.ISiigoClient.
func (m *SiigoClientMock) TestAuthentication(
	ctx context.Context,
	username, accessKey, accountID, partnerID, baseURL string,
) error {
	if m.TestAuthenticationFn != nil {
		return m.TestAuthenticationFn(ctx, username, accessKey, accountID, partnerID, baseURL)
	}
	return nil
}

// CreateInvoice implementa ports.ISiigoClient.
func (m *SiigoClientMock) CreateInvoice(
	ctx context.Context,
	req *dtos.CreateInvoiceRequest,
) (*dtos.CreateInvoiceResult, error) {
	if m.CreateInvoiceFn != nil {
		return m.CreateInvoiceFn(ctx, req)
	}
	return &dtos.CreateInvoiceResult{}, nil
}

// GetCustomerByIdentification implementa ports.ISiigoClient.
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

// CreateCustomer implementa ports.ISiigoClient.
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

// ListInvoices implementa ports.ISiigoClient.
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
