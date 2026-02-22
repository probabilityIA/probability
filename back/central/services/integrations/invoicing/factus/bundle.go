package factus

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Bundle implementa core.IIntegrationContract para Factus
// Este bundle NO tiene base de datos propia - es un cliente HTTP puro + async RabbitMQ consumer
type Bundle struct {
	client          ports.IFactusClient
	coreIntegration core.IIntegrationCore
	log             log.ILogger
}

// New crea una nueva instancia del módulo Factus
func New(
	config env.IConfig,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationCore,
) *Bundle {
	logger = logger.WithModule("factus")

	// 1. Cliente HTTP de Factus
	apiURL := config.Get("FACTUS_API_URL")
	if apiURL == "" {
		apiURL = "https://api.factus.com.co" // Default
		logger.Warn(context.Background()).
			Msg("FACTUS_API_URL not configured, using default URL")
	} else {
		logger.Info(context.Background()).
			Str("api_url", apiURL).
			Str("env_var", "FACTUS_API_URL").
			Msg("🔍 Factus API URL loaded from environment")
	}

	httpClient := client.New(apiURL, logger)
	logger.Info(context.Background()).
		Str("api_url", apiURL).
		Msg("✅ Factus HTTP client initialized")

	// 2. Response Publisher (RabbitMQ)
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Factus response publisher initialized")

	// 3. Invoice Request Consumer (escucha "invoicing.factus.requests")
	var invoiceRequestConsumer *consumer.InvoiceRequestConsumer
	if rabbit != nil {
		invoiceRequestConsumer = consumer.NewInvoiceRequestConsumer(
			rabbit,
			coreIntegration,
			httpClient,
			responsePublisher,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ Factus invoice request consumer initialized")

		// 4. Iniciar consumer en goroutine
		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting Factus invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Factus invoice request consumer failed to start or stopped with error")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de facturación (Factus) deshabilitado")
	}

	logger.Info(context.Background()).Msg("✅ Factus bundle initialized (HTTP client + RabbitMQ async consumer)")

	return &Bundle{
		client:          httpClient,
		coreIntegration: coreIntegration,
		log:             logger,
	}
}

// ═══════════════════════════════════════════════════════════════
// MÉTODOS DE IIntegrationContract (OBLIGATORIOS)
// ═══════════════════════════════════════════════════════════════

// RegisterRoutes registra las rutas HTTP del módulo Factus
func (b *Bundle) RegisterRoutes(router *gin.RouterGroup) {
	b.log.Info(context.Background()).Msg("ℹ️ Factus has no HTTP routes (uses IntegrationCore for CRUD)")
}

// TestConnection prueba la conexión con Factus usando credenciales
func (b *Bundle) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	b.log.Info(ctx).
		Interface("config", config).
		Msg("🧪 Testing connection with Factus API")

	clientID, okClientID := credentials["client_id"].(string)
	clientSecret, okClientSecret := credentials["client_secret"].(string)
	username, okUsername := credentials["username"].(string)
	password, okPassword := credentials["password"].(string)
	apiURL, _ := credentials["api_url"].(string) // opcional

	b.log.Info(ctx).
		Bool("has_client_id", okClientID && clientID != "").
		Bool("has_client_secret", okClientSecret && clientSecret != "").
		Bool("has_username", okUsername && username != "").
		Bool("has_password", okPassword && password != "").
		Bool("has_api_url", apiURL != "").
		Msg("📋 Factus credentials validation")

	if !okClientID || clientID == "" {
		return fmt.Errorf("el campo client_id es requerido en las credenciales")
	}
	if !okClientSecret || clientSecret == "" {
		return fmt.Errorf("el campo client_secret es requerido en las credenciales")
	}
	if !okUsername || username == "" {
		return fmt.Errorf("el campo username (email) es requerido en las credenciales")
	}
	if !okPassword || password == "" {
		return fmt.Errorf("el campo password es requerido en las credenciales")
	}

	b.log.Info(ctx).Msg("🔌 Calling Factus client.TestAuthentication...")
	if err := b.client.TestAuthentication(ctx, apiURL, clientID, clientSecret, username, password); err != nil {
		b.log.Error(ctx).Err(err).Msg("❌ Factus connection test failed")
		return err
	}

	b.log.Info(ctx).Msg("✅ Factus connection test successful")
	return nil
}

// SyncOrdersByIntegrationID no aplica para integración de facturación
func (b *Bundle) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return fmt.Errorf("SyncOrdersByIntegrationID is not supported for invoicing integration (Factus)")
}

// SyncOrdersByBusiness no aplica para integración de facturación
func (b *Bundle) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	return fmt.Errorf("SyncOrdersByBusiness is not supported for invoicing integration (Factus)")
}

// GetWebhookURL no aplica para integración de facturación
func (b *Bundle) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*core.WebhookInfo, error) {
	return nil, fmt.Errorf("webhooks are not supported for invoicing integration (Factus)")
}
