package softpymes

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Bundle implementa core.IIntegrationContract para Softpymes
// Este bundle NO tiene base de datos propia - es un cliente HTTP puro + async RabbitMQ consumer
type Bundle struct {
	client          ports.ISoftpymesClient
	coreIntegration core.IIntegrationCore
	log             log.ILogger
}

// New crea una nueva instancia del módulo Softpymes
// Este bundle funciona completamente sin base de datos:
// - Cliente HTTP para comunicarse con API de Softpymes
// - Consumer de RabbitMQ async para facturación
// - IntegrationCore para obtener credenciales y config
func New(
	config env.IConfig,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationCore,
) *Bundle {
	logger = logger.WithModule("softpymes")

	// 1. Cliente HTTP de Softpymes
	apiURL := config.Get("SOFTPYMES_API_URL")
	if apiURL == "" {
		apiURL = "https://api.softpymes.com" // Default
		logger.Warn(context.Background()).
			Msg("SOFTPYMES_API_URL not configured, using default URL")
	} else {
		logger.Info(context.Background()).
			Str("api_url", apiURL).
			Str("env_var", "SOFTPYMES_API_URL").
			Msg("🔍 DEBUG: Softpymes API URL loaded from environment")
	}
	// Create concrete HTTP client
	httpClient := client.New(apiURL, logger)
	logger.Info(context.Background()).
		Str("api_url", apiURL).
		Msg("✅ Softpymes HTTP client initialized")

	// 2. Response Publisher (RabbitMQ)
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Response publisher initialized")

	// 3. Invoice Request Consumer (ÚNICO consumer - procesa requests desde Invoicing Module)
	invoiceRequestConsumer := consumer.NewInvoiceRequestConsumer(
		rabbit,
		coreIntegration,
		httpClient, // Concrete *client.Client type
		responsePublisher,
		logger,
	)
	logger.Info(context.Background()).Msg("✅ Invoice request consumer initialized")

	// 4. Iniciar consumer en goroutine
	go func() {
		ctx := context.Background()
		logger.Info(ctx).Msg("🚀 Starting Softpymes invoice request consumer in background...")
		if err := invoiceRequestConsumer.Start(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("❌ Invoice request consumer failed to start or stopped with error")
		}
	}()

	logger.Info(context.Background()).Msg("✅ Softpymes bundle initialized (HTTP client + RabbitMQ async consumer)")

	return &Bundle{
		client:          httpClient, // Implicitly converts to ports.ISoftpymesClient
		coreIntegration: coreIntegration,
		log:             logger,
	}
}

// ═══════════════════════════════════════════════════════════════
// MÉTODOS DE IIntegrationContract (OBLIGATORIOS)
// ═══════════════════════════════════════════════════════════════

// RegisterRoutes registra las rutas HTTP del módulo Softpymes
// Como este módulo ya no tiene CRUD propio (se maneja desde IntegrationCore),
// este método es un stub vacío que cumple con la interfaz
func (b *Bundle) RegisterRoutes(router *gin.RouterGroup) {
	b.log.Info(context.Background()).Msg("ℹ️ Softpymes has no HTTP routes (uses IntegrationCore for CRUD)")
	// No hay rutas - el CRUD de integraciones se maneja desde IntegrationCore
}

// TestConnection prueba la conexión con Softpymes usando credenciales
// Implementa core.IIntegrationContract
func (b *Bundle) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	b.log.Info(ctx).
		Interface("config", config).
		Msg("🧪 Testing connection with Softpymes API")

	// Extraer API key y API secret de las credenciales
	apiKey, okKey := credentials["api_key"].(string)
	apiSecret, okSecret := credentials["api_secret"].(string)

	// Extraer referer del config (identificación de la instancia del cliente)
	referer, okReferer := config["referer"].(string)

	b.log.Info(ctx).
		Bool("has_api_key", okKey && apiKey != "").
		Bool("has_api_secret", okSecret && apiSecret != "").
		Bool("has_referer", okReferer && referer != "").
		Int("api_key_length", len(apiKey)).
		Int("api_secret_length", len(apiSecret)).
		Msg("📋 Credentials and config validation")

	if !okKey || apiKey == "" {
		b.log.Error(ctx).Msg("❌ API key is missing or empty")
		return fmt.Errorf("api_key is required in credentials")
	}

	if !okSecret || apiSecret == "" {
		b.log.Error(ctx).Msg("❌ API secret is missing or empty")
		return fmt.Errorf("api_secret is required in credentials")
	}

	if !okReferer || referer == "" {
		b.log.Error(ctx).Msg("❌ Referer is missing or empty in config")
		return fmt.Errorf("referer is required in config (identificación de instancia del cliente)")
	}

	// Usar el cliente para probar la conexión
	b.log.Info(ctx).Msg("🔌 Calling client.TestAuthentication...")
	if err := b.client.TestAuthentication(ctx, apiKey, apiSecret, referer); err != nil {
		b.log.Error(ctx).
			Err(err).
			Msg("❌ Softpymes connection test failed")
		return fmt.Errorf("failed to connect to Softpymes: %w", err)
	}

	b.log.Info(ctx).Msg("✅ Softpymes connection test successful")
	return nil
}

// SyncOrdersByIntegrationID no aplica para integración de facturación
// Implementa core.IIntegrationContract (requerido pero no usado)
func (b *Bundle) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return fmt.Errorf("SyncOrdersByIntegrationID is not supported for invoicing integration (Softpymes)")
}

// SyncOrdersByBusiness no aplica para integración de facturación
// Implementa core.IIntegrationContract (requerido pero no usado)
func (b *Bundle) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	return fmt.Errorf("SyncOrdersByBusiness is not supported for invoicing integration (Softpymes)")
}

// GetWebhookURL no aplica para integración de facturación
// Implementa core.IIntegrationContract (requerido pero no usado)
func (b *Bundle) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*core.WebhookInfo, error) {
	return nil, fmt.Errorf("webhooks are not supported for invoicing integration (Softpymes)")
}

// ═══════════════════════════════════════════════════════════════
// MÉTODOS ESPECÍFICOS DE FACTURACIÓN (PÚBLICOS)
// Estos métodos son usados por modules/invoicing
// ═══════════════════════════════════════════════════════════════

// CreateInvoice crea una factura en Softpymes
// Este método es llamado por modules/invoicing cuando se necesita facturar
func (b *Bundle) CreateInvoice(ctx context.Context, invoiceData map[string]interface{}) error {
	b.log.Info(ctx).
		Interface("data", invoiceData).
		Msg("Creating invoice in Softpymes via bundle")

	return b.client.CreateInvoice(ctx, invoiceData)
}

// CreateCreditNote crea una nota crédito en Softpymes
// Este método es llamado por modules/invoicing cuando se necesita una nota de crédito
func (b *Bundle) CreateCreditNote(ctx context.Context, creditNoteData map[string]interface{}) error {
	b.log.Info(ctx).
		Interface("data", creditNoteData).
		Msg("Creating credit note in Softpymes via bundle")

	return b.client.CreateCreditNote(ctx, creditNoteData)
}

// GetDocumentByNumber consulta un documento completo por su número
// Usado para consulta posterior después de crear factura (esperar procesamiento DIAN ~3seg)
// Retorna el documento completo con todos sus detalles (items, totales, información de envío)
//
// Este método es llamado por modules/invoicing después de crear una factura exitosamente
// para enriquecer la factura con información completa del documento procesado
func (b *Bundle) GetDocumentByNumber(ctx context.Context, apiKey, apiSecret, referer, documentNumber string) (map[string]interface{}, error) {
	b.log.Info(ctx).
		Str("document_number", documentNumber).
		Msg("📄 Getting document by number via bundle")

	return b.client.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, documentNumber)
}
