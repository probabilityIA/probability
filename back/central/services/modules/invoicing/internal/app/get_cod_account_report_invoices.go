package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

func (uc *useCase) GetCODAccountReportInvoices(ctx context.Context, businessID uint, startDate, endDate, accountNumber string, isCOD bool, page, pageSize int) ([]*entities.Invoice, int64, error) {
	uc.log.Info(ctx).
		Uint("business_id", businessID).
		Str("account_number", accountNumber).
		Bool("is_cod", isCOD).
		Int("page", page).
		Msg("Getting COD account report invoices")

	invoices, total, err := uc.repo.GetCODAccountReportInvoices(ctx, businessID, startDate, endDate, accountNumber, isCOD, page, pageSize)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get COD account report invoices")
		return nil, 0, err
	}

	return invoices, total, nil
}
