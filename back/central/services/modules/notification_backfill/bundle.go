package notification_backfill

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/infra/secondary/businesses"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/infra/secondary/registry"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/infra/secondary/runner"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/infra/secondary/selectors"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type IBundle interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type bundle struct {
	handler handlers.IHandlers
}

func New(
	database db.IDatabase,
	rabbitMQ rabbitmq.IQueue,
	logger log.ILogger,
	environment env.IConfig,
	guideDispatcher selectors.GuideDispatcher,
	confirmationDispatcher selectors.ConfirmationDispatcher,
) IBundle {
	logger = logger.WithModule("notification_backfill")

	imageURLBase := ""
	if environment != nil {
		imageURLBase = environment.Get("URL_BASE_DOMAIN_S3")
	}

	reg := registry.New()
	if guideDispatcher != nil {
		reg.Register(selectors.NewGuideSelector(database, logger, guideDispatcher, imageURLBase))
	}
	if confirmationDispatcher != nil {
		reg.Register(selectors.NewConfirmationSelector(database, logger, confirmationDispatcher))
	}

	store := runner.NewStore()
	progress := runner.NewProgressPublisher(rabbitMQ, logger)
	businessResolver := businesses.NewResolver(database)

	uc := app.New(reg, store, progress, businessResolver, logger)
	handler := handlers.New(uc, logger)

	return &bundle{handler: handler}
}

func (b *bundle) RegisterRoutes(router *gin.RouterGroup) {
	b.handler.RegisterRoutes(router)
}
