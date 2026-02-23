package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
)

// ProcessOrderForInvoicing implementa ports.IInvoiceUseCase.
// El procesamiento real de órdenes se realiza directamente en el consumer (InvoiceRequestConsumer),
// que tiene acceso directo al cliente HTTP. Este stub satisface la interfaz.
func (uc *invoicingUseCase) ProcessOrderForInvoicing(_ context.Context, _ *ports.OrderEventMessage) error {
	return nil
}
