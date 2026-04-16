package business

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/app/usecasebusiness"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/app/usecasebusinesstype"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo de business
func New(router *gin.RouterGroup, db db.IDatabase, logger log.ILogger, cfg env.IConfig, s3Service domain.IS3Service) {
	// 1. Inicializar Repositorio
	repo := repository.New(db, logger)

	// 2. Inicializar Casos de Uso
	businessUC := usecasebusiness.New(repo, logger, s3Service, cfg)
	businessTypeUC := usecasebusinesstype.New(repo, logger)

	// 3. Inicializar Handlers
	businessH := businesshandler.New(businessUC, logger)
	businessTypeH := businesstypehandler.New(businessTypeUC, logger)

	// 4. Registrar Rutas
	businessH.RegisterRoutes(router, businessH)
	businesstypehandler.RegisterRoutes(router, businessTypeH)
}
