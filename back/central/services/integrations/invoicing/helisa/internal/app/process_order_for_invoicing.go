package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/domain/ports"
)

// ProcessOrderForInvoicing no aplica para Helisa en el flujo de eventos directos.
// La facturación con Helisa se maneja exclusivamente a través del consumer
// de la cola "invoicing.helisa.requests" (InvoiceRequestConsumer),
// que recibe las solicitudes enrutadas desde el módulo Invoicing principal.
func (uc *invoicingUseCase) ProcessOrderForInvoicing(ctx context.Context, event *ports.OrderEventMessage) error {
	uc.log.Warn(ctx).
		Str("order_id", event.OrderID).
		Msg("ProcessOrderForInvoicing not used - Helisa handles invoicing via invoicing.helisa.requests queue")
	return fmt.Errorf("helisa: use invoice request queue (invoicing.helisa.requests)")
}
