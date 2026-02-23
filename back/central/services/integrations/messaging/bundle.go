package messaging

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	whatsapp "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp"
	"github.com/secamc93/probability/back/central/services/modules"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa todos los proveedores de mensajería y los registra en integrationCore.
// Debe llamarse después de inicializar integrationCore.
func New(
	config env.IConfig,
	logger log.ILogger,
	database db.IDatabase,
	rabbitMQ rabbitmq.IQueue,
	redisClient redisclient.IRedis,
	moduleBundles *modules.ModuleBundles,
	integrationCore core.IIntegrationCore,
) {
	// WhatsApp (type_id=2)
	whatsappBundle := whatsapp.New(config, logger, database, rabbitMQ, redisClient, moduleBundles)
	integrationCore.RegisterIntegration(core.IntegrationTypeWhatsApp, whatsappBundle)

	// SMS — pendiente de implementación
	// smsBundle := sms.New(config, logger, rabbitMQ)
	// integrationCore.RegisterIntegration(core.IntegrationTypeSMS, smsBundle)
}
