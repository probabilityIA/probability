package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
)

// ProcessOrderForInvoicing no aplica para Siigo en el flujo de eventos directos.
// La facturación con Siigo se maneja exclusivamente a través del consumer
// de la cola "invoicing.siigo.requests" (InvoiceRequestConsumer),
// que recibe las solicitudes enrutadas desde el módulo Invoicing principal.
func (uc *invoicingUseCase) ProcessOrderForInvoicing(ctx context.Context, event *ports.OrderEventMessage) error {
	uc.log.Warn(ctx).
		Str("order_id", event.OrderID).
		Msg("ProcessOrderForInvoicing not used - Siigo handles invoicing via invoicing.siigo.requests queue")
	return fmt.Errorf("siigo: use invoice request queue (invoicing.siigo.requests)")
}
