package orders

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/app"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/testing/shared/db"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// IWebhookSimulator is the public interface for webhook simulators.
// Re-exported from internal ports so external packages (cmd/main.go) can use it.
type IWebhookSimulator = ports.IWebhookSimulator

func New(router *gin.RouterGroup, database db.IDatabase, centralAPIURL string, logger log.ILogger, simulators map[string]IWebhookSimulator) {
	repo := repository.New(database)
	centralClient := client.New(centralAPIURL)
	useCase := app.New(repo, centralClient, logger, simulators)
	handler := handlers.New(useCase, logger)

	handler.RegisterRoutes(router)

	logger.Info().Msg("Orders testing module initialized")
}
