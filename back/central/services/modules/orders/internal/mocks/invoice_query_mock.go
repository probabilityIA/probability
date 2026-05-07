package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/stretchr/testify/mock"
)

type InvoiceQueryMock struct {
	mock.Mock
}

func (m *InvoiceQueryMock) GetInvoiceByOrderID(ctx context.Context, orderID string) (*dtos.InvoiceData, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.InvoiceData), args.Error(1)
}
