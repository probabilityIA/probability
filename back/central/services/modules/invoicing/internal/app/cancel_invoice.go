package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// CancelInvoice cancela una factura emitida
// NOT IMPLEMENTED: Pendiente de re-implementar usando integrationCore + softpymes bundle
func (uc *useCase) CancelInvoice(ctx context.Context, dto *dtos.CancelInvoiceDTO) error {
	return errors.ErrCancelNotImplemented
}
