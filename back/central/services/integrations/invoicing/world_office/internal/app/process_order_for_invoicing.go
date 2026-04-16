package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/domain/ports"
)

// ProcessOrderForInvoicing no aplica para World Office en el flujo de eventos directos.
// La facturación con World Office se maneja exclusivamente a través del consumer
// de la cola "invoicing.world_office.requests" (InvoiceRequestConsumer),
// que recibe las solicitudes enrutadas desde el módulo Invoicing principal.
func (uc *invoicingUseCase) ProcessOrderForInvoicing(ctx context.Context, event *ports.OrderEventMessage) error {
	uc.log.Warn(ctx).
		Str("order_id", event.OrderID).
		Msg("ProcessOrderForInvoicing not used - World Office handles invoicing via invoicing.world_office.requests queue")
	return fmt.Errorf("world_office: use invoice request queue (invoicing.world_office.requests)")
}
