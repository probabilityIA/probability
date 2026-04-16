package transport

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick"
	"github.com/secamc93/probability/back/central/services/integrations/transport/mipaquete"
	"github.com/secamc93/probability/back/central/services/integrations/transport/router"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New initializes all transport providers and the transport router.
func New(
	logger log.ILogger,
	rabbitMQ rabbitmq.IQueue,
	integrationCore core.IIntegrationCore,
) {
	// EnvioClick (type_id=12)
	// integrationCore satisfies ICredentialResolver (has DecryptCredential)
	envioclick.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeEnvioClick, nil)

	// Enviame (type_id=13)
	enviame.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeEnviame, nil)

	// MiPaquete (type_id=15)
	mipaquete.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeMiPaquete, nil)

	// Router: consumes transport.requests and routes to carrier queues.
	// Initialized last so carrier queues are already declared.
	router.New(logger, rabbitMQ)
}
