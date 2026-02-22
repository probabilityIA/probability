package siigo

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Bundle implementa core.IIntegrationContract para Siigo
// Este bundle NO tiene base de datos propia - es un cliente HTTP puro + async RabbitMQ consumer
type Bundle struct {
	client          ports.ISiigoClient
	coreIntegration core.IIntegrationCore
	log             log.ILogger
}

// New crea una nueva instancia del módulo Siigo
func New(
	config env.IConfig,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationCore,
) *Bundle {
	logger = logger.WithModule("siigo")

	// 1. Cliente HTTP de Siigo
	apiURL := config.Get("SIIGO_API_URL")
	if apiURL == "" {
		apiURL = "https://api.siigo.com" // Default
		logger.Warn(context.Background()).
			Msg("SIIGO_API_URL not configured, using default URL")
	} else {
		logger.Info(context.Background()).
			Str("api_url", apiURL).
			Str("env_var", "SIIGO_API_URL").
			Msg("🔍 Siigo API URL loaded from environment")
	}

	httpClient := client.New(apiURL, logger)
	logger.Info(context.Background()).
		Str("api_url", apiURL).
		Msg("✅ Siigo HTTP client initialized")

	// 2. Response Publisher (RabbitMQ)
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Siigo response publisher initialized")

	// 3. Invoice Request Consumer (escucha "invoicing.siigo.requests")
	var invoiceRequestConsumer *consumer.InvoiceRequestConsumer
	if rabbit != nil {
		invoiceRequestConsumer = consumer.NewInvoiceRequestConsumer(
			rabbit,
			coreIntegration,
			httpClient,
			responsePublisher,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ Siigo invoice request consumer initialized")

		// 4. Iniciar consumer en goroutine
		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting Siigo invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Siigo invoice request consumer failed to start or stopped with error")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de facturación (Siigo) deshabilitado")
	}

	logger.Info(context.Background()).Msg("✅ Siigo bundle initialized (HTTP client + RabbitMQ async consumer)")

	return &Bundle{
		client:          httpClient,
		coreIntegration: coreIntegration,
		log:             logger,
	}
}

// ═══════════════════════════════════════════════════════════════
// MÉTODOS DE IIntegrationContract (OBLIGATORIOS)
// ═══════════════════════════════════════════════════════════════

// RegisterRoutes registra las rutas HTTP del módulo Siigo
func (b *Bundle) RegisterRoutes(router *gin.RouterGroup) {
	b.log.Info(context.Background()).Msg("ℹ️ Siigo has no HTTP routes (uses IntegrationCore for CRUD)")
}

// TestConnection prueba la conexión con Siigo usando credenciales
func (b *Bundle) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	b.log.Info(ctx).
		Interface("config", config).
		Msg("🧪 Testing connection with Siigo API")

	username, okUsername := credentials["username"].(string)
	accessKey, okAccessKey := credentials["access_key"].(string)
	accountID, okAccountID := credentials["account_id"].(string)
	partnerID, okPartnerID := credentials["partner_id"].(string)
	apiURL, _ := credentials["api_url"].(string) // opcional

	b.log.Info(ctx).
		Bool("has_username", okUsername && username != "").
		Bool("has_access_key", okAccessKey && accessKey != "").
		Bool("has_account_id", okAccountID && accountID != "").
		Bool("has_partner_id", okPartnerID && partnerID != "").
		Bool("has_api_url", apiURL != "").
		Msg("📋 Siigo credentials validation")

	if !okUsername || username == "" {
		return fmt.Errorf("el campo username (email) es requerido en las credenciales")
	}
	if !okAccessKey || accessKey == "" {
		return fmt.Errorf("el campo access_key es requerido en las credenciales")
	}
	if !okAccountID || accountID == "" {
		return fmt.Errorf("el campo account_id (subscription key) es requerido en las credenciales")
	}
	if !okPartnerID || partnerID == "" {
		return fmt.Errorf("el campo partner_id es requerido en las credenciales")
	}

	b.log.Info(ctx).Msg("🔌 Calling Siigo client.TestAuthentication...")
	if err := b.client.TestAuthentication(ctx, username, accessKey, accountID, partnerID, apiURL); err != nil {
		b.log.Error(ctx).Err(err).Msg("❌ Siigo connection test failed")
		return err
	}

	b.log.Info(ctx).Msg("✅ Siigo connection test successful")
	return nil
}

// SyncOrdersByIntegrationID no aplica para integración de facturación
func (b *Bundle) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return fmt.Errorf("SyncOrdersByIntegrationID is not supported for invoicing integration (Siigo)")
}

// SyncOrdersByBusiness no aplica para integración de facturación
func (b *Bundle) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	return fmt.Errorf("SyncOrdersByBusiness is not supported for invoicing integration (Siigo)")
}

// GetWebhookURL no aplica para integración de facturación
func (b *Bundle) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*core.WebhookInfo, error) {
	return nil, fmt.Errorf("webhooks are not supported for invoicing integration (Siigo)")
}
