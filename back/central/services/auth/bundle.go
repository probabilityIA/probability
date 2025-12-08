package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/actions"
	business "github.com/secamc93/probability/back/central/services/auth/bussines"
	"github.com/secamc93/probability/back/central/services/auth/login"
	"github.com/secamc93/probability/back/central/services/auth/permissions"
	"github.com/secamc93/probability/back/central/services/auth/resources"
	"github.com/secamc93/probability/back/central/services/auth/roles"
	"github.com/secamc93/probability/back/central/services/auth/users"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

// New inicializa todos los módulos de autenticación y autorización
// Este bundle coordina la inicialización de todos los submódulos de auth
// (login, permissions, roles, users, business, actions, resources)
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, s3Service storage.IS3Service) {
	// Inicializar módulo de login
	login.New(router, database, logger, environment)

	// Inicializar módulo de permissions
	permissions.New(router, database, logger)

	// Inicializar módulo de roles
	roles.New(router, database, logger)

	// Inicializar módulo de users
	users.New(router, database, logger, environment, s3Service)

	// Inicializar módulo de business
	business.New(router, database, logger, environment, s3Service)

	// Inicializar módulo de actions
	actions.New(database, logger, router)

	// Inicializar módulo de resources
	resources.New(database, logger, router)
}
