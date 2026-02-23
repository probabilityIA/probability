package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/ports"
)

// FactusClientMock mock de ports.IFactusClient para tests unitarios.
// Cada método de la interfaz tiene su correspondiente campo Fn que permite
// inyectar el comportamiento deseado en cada test.
type FactusClientMock struct {
	TestAuthenticationFn func(ctx context.Context, baseURL, clientID, clientSecret, username, password string) error
	CreateInvoiceFn      func(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error)
	ListBillsFn          func(ctx context.Context, credentials dtos.Credentials, params dtos.ListBillsParams) (*dtos.ListBillsResult, error)
	GetBillByNumberFn    func(ctx context.Context, credentials dtos.Credentials, number string) (*dtos.BillDetail, error)
}

// Verificar en tiempo de compilación que FactusClientMock implementa la interfaz.
var _ ports.IFactusClient = (*FactusClientMock)(nil)

func (m *FactusClientMock) TestAuthentication(ctx context.Context, baseURL, clientID, clientSecret, username, password string) error {
	if m.TestAuthenticationFn != nil {
		return m.TestAuthenticationFn(ctx, baseURL, clientID, clientSecret, username, password)
	}
	return nil
}

func (m *FactusClientMock) CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
	if m.CreateInvoiceFn != nil {
		return m.CreateInvoiceFn(ctx, req)
	}
	return &dtos.CreateInvoiceResult{}, nil
}

func (m *FactusClientMock) ListBills(ctx context.Context, credentials dtos.Credentials, params dtos.ListBillsParams) (*dtos.ListBillsResult, error) {
	if m.ListBillsFn != nil {
		return m.ListBillsFn(ctx, credentials, params)
	}
	return &dtos.ListBillsResult{}, nil
}

func (m *FactusClientMock) GetBillByNumber(ctx context.Context, credentials dtos.Credentials, number string) (*dtos.BillDetail, error) {
	if m.GetBillByNumberFn != nil {
		return m.GetBillByNumberFn(ctx, credentials, number)
	}
	return &dtos.BillDetail{}, nil
}
