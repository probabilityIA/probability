package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/integration_cache"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// NewInvoicingUseCaseForBundle crea el use case de facturación automática
// Este constructor se usa desde bundle.go para inicializar el consumer
// NO requiere repositorios - todo funciona con RabbitMQ + Redis + IntegrationCore
func NewInvoicingUseCaseForBundle(
	softpymesClient ports.ISoftpymesClient,
	configCache cache.IConfigCache,
	redisClient redis.IRedis,
	integrationCore core.IIntegrationCore,
	integrationCache integration_cache.IIntegrationCacheClient, // ✅ NUEVO
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return NewInvoicingUseCase(
		softpymesClient,
		configCache,
		redisClient,
		integrationCore,
		integrationCache, // ✅ NUEVO
		logger,
	)
}
