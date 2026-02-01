package softpymes

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Bundle implementa core.IIntegrationContract para Softpymes
type Bundle struct {
	useCase         app.IUseCase
	handler         handlers.IHandler
	client          ports.ISoftpymesClient
	coreIntegration core.IIntegrationCore
	log             log.ILogger
}

// New crea una nueva instancia del módulo Softpymes
// Este bundle se registra en integrationCore y proporciona funcionalidades
// de gestión de proveedores de facturación electrónica
func New(
	config env.IConfig,
	logger log.ILogger,
	database db.IDatabase,
	coreIntegration core.IIntegrationCore,
) *Bundle {
	logger = logger.WithModule("softpymes")

	// 1. Inicializar cliente HTTP de Softpymes
	apiURL := config.Get("SOFTPYMES_API_URL")
	if apiURL == "" {
		apiURL = "https://api.softpymes.com" // Default
		logger.Warn(context.Background()).
			Msg("SOFTPYMES_API_URL not configured, using default URL")
	}
	softpymesClient := client.New(apiURL, logger)

	// 2. Inicializar repositorios
	repos := repository.New(database, logger)

	// 3. Inicializar casos de uso
	useCase := app.New(
		repos.Provider,
		repos.ProviderType,
		softpymesClient,
		logger,
	)

	// 4. Inicializar handlers HTTP
	handler := handlers.New(useCase, logger)

	logger.Info(context.Background()).
		Str("api_url", apiURL).
		Msg("Softpymes integration module initialized")

	return &Bundle{
		useCase:        useCase,
		handler:        handler,
		client:         softpymesClient,
		coreIntegration: coreIntegration,
		log:            logger,
	}
}

// RegisterRoutes registra las rutas HTTP del módulo Softpymes
// Implementa método requerido por el patrón de integración
func (b *Bundle) RegisterRoutes(router *gin.RouterGroup) {
	b.handler.RegisterRoutes(router)
	b.log.Info(context.Background()).Msg("Softpymes routes registered")
}

// ═══════════════════════════════════════════════════════════════
// MÉTODOS DE IIntegrationContract (OBLIGATORIOS)
// ═══════════════════════════════════════════════════════════════

// TestConnection prueba la conexión con Softpymes usando credenciales
// Implementa core.IIntegrationContract
func (b *Bundle) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	b.log.Info(ctx).Msg("Testing connection with Softpymes API")

	// Extraer API key de las credenciales
	apiKey, ok := credentials["api_key"].(string)
	if !ok || apiKey == "" {
		return fmt.Errorf("api_key is required in credentials")
	}

	// Usar el cliente para probar la conexión
	if err := b.client.TestAuthentication(ctx, apiKey); err != nil {
		b.log.Error(ctx).Err(err).Msg("Softpymes connection test failed")
		return fmt.Errorf("failed to connect to Softpymes: %w", err)
	}

	b.log.Info(ctx).Msg("Softpymes connection test successful")
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

// GetProviderByID obtiene un proveedor por ID
// Método de conveniencia para acceso externo
func (b *Bundle) GetProviderByID(ctx context.Context, providerID uint) (interface{}, error) {
	return b.useCase.GetProvider(ctx, providerID)
}

// TestProviderConnection prueba la conexión de un proveedor específico
// Método de conveniencia para acceso externo
func (b *Bundle) TestProviderConnection(ctx context.Context, providerID uint) error {
	return b.useCase.TestProviderConnection(ctx, providerID)
}
