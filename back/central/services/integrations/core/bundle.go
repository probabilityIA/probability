package core

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/secondary/encryption"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo core de integraciones y retorna la interfaz pública
func New(
	router *gin.RouterGroup,
	db db.IDatabase,
	logger log.ILogger,
	config env.IConfig,
) IIntegrationCore {
	// 1. Inicializar Servicio de Encriptación
	encryptionService := encryption.New(config, logger)

	// 2. Inicializar Repositorio
	repo := repository.New(db, logger, encryptionService)

	// 3. Inicializar Casos de Uso
	useCase := app.New(repo, encryptionService, logger)

	// 4. Inicializar Handlers
	handler := handlers.New(useCase, logger)

	// 5. Registrar Rutas
	handler.RegisterRoutes(router, handler, logger)

	// 6. Crear y retornar interfaz pública
	return NewIntegrationCore(useCase)
}
