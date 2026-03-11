package subscriptions

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/secondary/postgres"
	"github.com/secamc93/probability/back/central/shared/db"
)

type Dependencies struct {
	Handler *handlers.SubscriptionHandler
}

func Setup(database db.IDatabase) *Dependencies {
	repo := postgres.NewSubscriptionRepository(database)
	uc := usecases.NewSubscriptionUsecase(repo)
	handler := handlers.NewSubscriptionHandler(uc)

	return &Dependencies{
		Handler: handler,
	}
}

func (d *Dependencies) RegisterRoutes(r *gin.RouterGroup) {
	subs := r.Group("/subscriptions")
	{
		// Endpoints for business (require auth middleware later)
		subs.GET("/me", d.Handler.GetCurrentSubscription)

		// Endpoints for super admin (require admin middleware later)
		subs.POST("/register-payment", d.Handler.RegisterPayment)
		subs.POST("/disable", d.Handler.DisableSubscription)
	}
}
