package orderstatus

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// NewRepository crea y retorna un nuevo repositorio de order status
// Esta función pública permite que otros módulos accedan al repositorio
func NewRepository(database db.IDatabase, logger log.ILogger) IRepository {
	return repository.New(database, logger)
}

// New inicializa el módulo de order status mappings
// Sigue el patrón de arquitectura hexagonal:
//   1. Infraestructura secundaria (adaptadores de salida - repositorios)
//   2. Capa de aplicación (casos de uso)
//   3. Infraestructura primaria (adaptadores de entrada - handlers HTTP)
//   4. Registro de rutas
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig) {
	logger = logger.WithModule("orderstatus")

	// 1. Repositorio (infraestructura secundaria)
	repo := repository.New(database, logger)

	// 2. Casos de uso (capa de aplicación)
	useCase := app.New(repo, logger)

	// 3. Handlers (infraestructura primaria)
	handler := handlers.New(useCase, logger, environment)

	// 4. Rutas
	handler.RegisterRoutes(router)
}
