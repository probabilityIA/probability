package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/ports"
)

// ProcessOrderForInvoicing procesa un evento de orden para determinar si debe facturarse con Factus
// NOTA: Para Factus, el flujo de facturación se maneja a través del consumer de RabbitMQ
// (invoice_request_consumer.go) que escucha la cola "invoicing.factus.requests".
// Este método es un stub para cumplir con la interfaz IInvoiceUseCase.
func (uc *invoicingUseCase) ProcessOrderForInvoicing(ctx context.Context, event *ports.OrderEventMessage) error {
	uc.log.Warn(ctx).
		Str("event_id", event.EventID).
		Str("order_id", event.OrderID).
		Msg("ProcessOrderForInvoicing called on Factus use case - this should be handled by the invoice request consumer")

	return fmt.Errorf("factus: direct order processing not supported, use invoice request queue")
}
