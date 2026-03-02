package messaging

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	emailmod "github.com/secamc93/probability/back/central/services/integrations/messaging/email"
	whatsapp "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa todos los proveedores de mensajería y los registra en integrationCore.
// Debe llamarse después de inicializar integrationCore.
// Nota: Ningún proveedor de mensajería requiere DB directa.
// WhatsApp usa Redis cache + RabbitMQ. Email es stateless.
func New(
	config env.IConfig,
	logger log.ILogger,
	rabbitMQ rabbitmq.IQueue,
	redisClient redisclient.IRedis,
	integrationCore core.IIntegrationCore,
	emailService email.IEmailService,
) {
	// WhatsApp (type_id=2) — cache-first, DB-async via RabbitMQ
	whatsappBundle := whatsapp.New(config, logger, rabbitMQ, redisClient)
	integrationCore.RegisterIntegration(core.IntegrationTypeWhatsApp, whatsappBundle)

	// Email (type_id=29) — notificaciones por correo via SES (stateless, sin DB)
	emailBundle := emailmod.New(logger, rabbitMQ, emailService)
	integrationCore.RegisterIntegration(core.IntegrationTypeEmail, emailBundle)

	// SMS — pendiente de implementación
	// smsBundle := sms.New(config, logger, rabbitMQ)
	// integrationCore.RegisterIntegration(core.IntegrationTypeSMS, smsBundle)
}
