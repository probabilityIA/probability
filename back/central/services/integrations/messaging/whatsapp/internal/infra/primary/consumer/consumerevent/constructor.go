package consumerevent

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/consumer/consumerevent/request"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/log"
	rabbitmq "github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// IConsumer define la interfaz para el consumer de eventos de órdenes
type IConsumer interface {
	Start(ctx context.Context) error
}

// consumer consume eventos de órdenes desde Redis y valida configuraciones
type consumer struct {
	redisClient            redisclient.IRedis
	rabbitMQ               rabbitmq.IQueue
	notificationConfigRepo repository.INotificationConfigRepository // ← Usa repositorio
	integrationRepo        request.IntegrationRepository
	orderRepo              request.OrderRepository
	logger                 log.ILogger
	channel                string
}

// New crea una nueva instancia del consumidor de eventos de órdenes
func New(
	redisClient redisclient.IRedis,
	rabbitMQ rabbitmq.IQueue,
	notificationConfigRepo repository.INotificationConfigRepository, // ← Recibe repositorio
	integrationRepo request.IntegrationRepository,
	orderRepo request.OrderRepository,
	logger log.ILogger,
	channel string,
) IConsumer {
	return &consumer{
		redisClient:            redisClient,
		rabbitMQ:               rabbitMQ,
		notificationConfigRepo: notificationConfigRepo, // ← Guarda repositorio
		integrationRepo:        integrationRepo,
		orderRepo:              orderRepo,
		logger:                 logger.WithModule("whatsapp_order_event_consumer"),
		channel:                channel,
	}
}
