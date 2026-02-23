package router

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	// QueueInvoicingRequests es la cola unificada de entrada para todas las solicitudes de facturación.
	// modules/invoicing publica aquí; el router decide a qué facturador enrutar.
	QueueInvoicingRequests = "invoicing.requests"

	// Colas de cada proveedor (lectura de los consumers de cada módulo)
	QueueSoftpymesRequests  = "invoicing.softpymes.requests"
	QueueFactusRequests     = "invoicing.factus.requests"
	QueueSiigoRequests      = "invoicing.siigo.requests"
	QueueAlegraRequests     = "invoicing.alegra.requests"
	QueueWorldOfficeRequests = "invoicing.world_office.requests"
	QueueHelisaRequests     = "invoicing.helisa.requests"
)

// invoiceRequestHeader contiene solo los campos necesarios para enrutar el mensaje.
// El resto del payload se reenvía sin transformación.
type invoiceRequestHeader struct {
	InvoiceID     uint      `json:"invoice_id"`
	Provider      string    `json:"provider"`
	Operation     string    `json:"operation"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
}

// Bundle es el router centralizado de facturación electrónica.
// Su única responsabilidad: leer de invoicing.requests y reenviar
// al consumer del proveedor correcto.
type Bundle struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// New crea e inicia el router centralizado de facturación.
// Debe llamarse desde services/integrations/bundle.go después de inicializar
// los bundles de todos los proveedores.
func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
) *Bundle {
	logger = logger.WithModule("invoicing.router")

	b := &Bundle{
		rabbit: rabbit,
		log:    logger,
	}

	if rabbit == nil {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, router de facturación (invoicing.router) deshabilitado")
		return b
	}

	// Iniciar router en goroutine (mismo patrón que los bundles de proveedores)
	go func() {
		ctx := context.Background()
		logger.Info(ctx).Msg("🚀 Starting invoicing router in background...")
		if err := b.startRouter(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("❌ Invoicing router failed to start or stopped with error")
		}
	}()

	logger.Info(context.Background()).Msg("✅ Invoicing router initialized")

	return b
}

// startRouter declara las colas necesarias e inicia el consumo de invoicing.requests
func (b *Bundle) startRouter(ctx context.Context) error {
	if b.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	// Declarar la cola unificada de entrada
	if err := b.rabbit.DeclareQueue(QueueInvoicingRequests, true); err != nil {
		b.log.Error(ctx).Err(err).Str("queue", QueueInvoicingRequests).Msg("❌ Failed to declare invoicing.requests queue")
		return err
	}

	// Declarar colas de proveedores para garantizar que existen antes de publicar
	providerQueues := []string{
		QueueSoftpymesRequests,
		QueueFactusRequests,
		QueueSiigoRequests,
		QueueAlegraRequests,
		QueueWorldOfficeRequests,
		QueueHelisaRequests,
	}
	for _, q := range providerQueues {
		if err := b.rabbit.DeclareQueue(q, true); err != nil {
			b.log.Warn(ctx).Err(err).Str("queue", q).Msg("⚠️ Failed to declare provider queue")
			// No es fatal: el consumer del proveedor también las declara al iniciar
		}
	}

	b.log.Info(ctx).
		Str("queue", QueueInvoicingRequests).
		Msg("✅ Invoicing router listening")

	return b.rabbit.Consume(ctx, QueueInvoicingRequests, b.handleInvoiceRequest)
}

// handleInvoiceRequest enruta una solicitud de facturación a la cola del proveedor correspondiente.
// El mensaje se reenvía íntegro (sin transformación) para que el consumer del proveedor
// lo procese con su lógica propia.
func (b *Bundle) handleInvoiceRequest(message []byte) error {
	ctx := context.Background()

	// Decodificar solo el encabezado para obtener el proveedor
	var header invoiceRequestHeader
	if err := json.Unmarshal(message, &header); err != nil {
		b.log.Error(ctx).
			Err(err).
			Str("body", string(message)).
			Msg("❌ Failed to unmarshal invoice request header")
		return err
	}

	b.log.Info(ctx).
		Uint("invoice_id", header.InvoiceID).
		Str("provider", header.Provider).
		Str("operation", header.Operation).
		Str("correlation_id", header.CorrelationID).
		Msg("📨 Routing invoice request")

	// Determinar la cola destino según el proveedor
	targetQueue := b.getProviderQueue(header.Provider)
	if targetQueue == "" {
		b.log.Error(ctx).
			Str("provider", header.Provider).
			Uint("invoice_id", header.InvoiceID).
			Msg("❌ Unknown provider — cannot route invoice request (message discarded)")
		// No retornar error: el mensaje no debe re-encolar si el proveedor es desconocido
		return nil
	}

	// Reenviar el mensaje original al consumer del proveedor
	if err := b.rabbit.Publish(ctx, targetQueue, message); err != nil {
		b.log.Error(ctx).
			Err(err).
			Str("target_queue", targetQueue).
			Uint("invoice_id", header.InvoiceID).
			Msg("❌ Failed to forward invoice request to provider queue")
		return err // Sí retornar error para que RabbitMQ reintente
	}

	b.log.Info(ctx).
		Uint("invoice_id", header.InvoiceID).
		Str("provider", header.Provider).
		Str("target_queue", targetQueue).
		Msg("✅ Invoice request forwarded")

	return nil
}

// getProviderQueue retorna el nombre de la cola para el proveedor dado.
// Retorna "" si el proveedor no está registrado.
func (b *Bundle) getProviderQueue(provider string) string {
	switch provider {
	case "softpymes":
		return QueueSoftpymesRequests
	case "factus":
		return QueueFactusRequests
	case "siigo":
		return QueueSiigoRequests
	case "alegra":
		return QueueAlegraRequests
	case "world_office":
		return QueueWorldOfficeRequests
	case "helisa":
		return QueueHelisaRequests
	default:
		return ""
	}
}
