package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/ai"
	"github.com/secamc93/probability/back/central/services/modules/ai_sales"
	"github.com/secamc93/probability/back/central/services/modules/announcements"
	"github.com/secamc93/probability/back/central/services/modules/customers"
	"github.com/secamc93/probability/back/central/services/modules/dashboard"
	"github.com/secamc93/probability/back/central/services/modules/drivers"
	"github.com/secamc93/probability/back/central/services/modules/inventory"
	"github.com/secamc93/probability/back/central/services/modules/invoicing"
	"github.com/secamc93/probability/back/central/services/modules/monitoring"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill"
	"github.com/secamc93/probability/back/central/services/modules/notification_config"
	"github.com/secamc93/probability/back/central/services/modules/orders"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus"
	"github.com/secamc93/probability/back/central/services/modules/pay"
	"github.com/secamc93/probability/back/central/services/modules/payments"
	"github.com/secamc93/probability/back/central/services/modules/probability"
	"github.com/secamc93/probability/back/central/services/modules/products"
	"github.com/secamc93/probability/back/central/services/modules/publicsite"
	"github.com/secamc93/probability/back/central/services/modules/routes"
	"github.com/secamc93/probability/back/central/services/modules/shipments"
	"github.com/secamc93/probability/back/central/services/modules/storefront"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions"
	"github.com/secamc93/probability/back/central/services/modules/vehicles"
	"github.com/secamc93/probability/back/central/services/modules/warehouses"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig"
	"github.com/secamc93/probability/back/central/shared/bedrock"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/storage"
)

type ModuleBundles struct {
	router      *gin.RouterGroup
	database    db.IDatabase
	logger      log.ILogger
	environment env.IConfig
	rabbitMQ    rabbitmq.IQueue
	redisClient redis.IRedis
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue, redisClient redis.IRedis, s3 storage.IS3Service, bedrockClient bedrock.IBedrock) *ModuleBundles {
	announcements.New(router, database, logger, s3)
	payments.New(router, database, logger, environment)
	orderstatus.New(router, database, logger, environment)
	ordersBundle := orders.New(router, database, logger, environment, rabbitMQ)
	probability.New(database, logger, rabbitMQ)
	products.New(router, database, logger, environment, s3)
	customers.New(router, database, logger, rabbitMQ)
	shipments.New(router, database, logger, environment, rabbitMQ, redisClient)
	notification_config.New(router, database, redisClient, logger, rabbitMQ)
	notification_backfill.New(database, rabbitMQ, logger, ordersBundle.SendGuideNotificationUC, ordersBundle.RequestConfirmationUC).RegisterRoutes(router)
	ai.New(router, logger)
	dashboard.New(router, database, logger)
	pay.New(router, database, logger, environment, rabbitMQ, redisClient)
	invoicing.New(router, database, logger, environment, rabbitMQ, redisClient)
	warehouses.New(router, database)
	inventory.New(router, database, logger, environment, rabbitMQ, redisClient)
	drivers.New(router, database)
	vehicles.New(router, database)
	routes.New(router, database)
	storefront.New(router, database, logger, rabbitMQ, environment)
	publicsite.New(router, database, logger, environment)
	websiteconfig.New(router, database, logger)

	subModule := subscriptions.Setup(database)
	subModule.RegisterRoutes(router)

	if rabbitMQ != nil {
		monitoring.New(router, logger, environment, rabbitMQ)
	} else {
		logger.Warn().Msg("RabbitMQ no disponible, modulo de monitoreo no se inicializara")
	}

	if bedrockClient != nil && rabbitMQ != nil && redisClient != nil {
		ai_sales.New(database, logger, rabbitMQ, redisClient, bedrockClient)
	} else {
		logger.Warn().Msg("AI Sales: Bedrock, RabbitMQ o Redis no disponible, modulo no se inicializara")
	}

	return &ModuleBundles{
		router:      router,
		database:    database,
		logger:      logger,
		environment: environment,
		rabbitMQ:    rabbitMQ,
		redisClient: redisClient,
	}
}
