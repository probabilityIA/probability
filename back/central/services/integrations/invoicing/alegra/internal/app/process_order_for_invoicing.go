package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra/internal/domain/ports"
)

// ProcessOrderForInvoicing no aplica para Alegra en el flujo de eventos directos.
// La facturación con Alegra se maneja exclusivamente a través del consumer
// de la cola "invoicing.alegra.requests" (InvoiceRequestConsumer),
// que recibe las solicitudes enrutadas desde el módulo Invoicing principal.
func (uc *invoicingUseCase) ProcessOrderForInvoicing(ctx context.Context, event *ports.OrderEventMessage) error {
	uc.log.Warn(ctx).
		Str("order_id", event.OrderID).
		Msg("ProcessOrderForInvoicing not used - Alegra handles invoicing via invoicing.alegra.requests queue")
	return fmt.Errorf("alegra: use invoice request queue (invoicing.alegra.requests)")
}
