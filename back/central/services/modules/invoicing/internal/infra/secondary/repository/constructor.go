package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	invoicingRedis "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Repository implementa LA interfaz IRepository con TODOS los métodos de persistencia del módulo
type Repository struct {
	db          db.IDatabase
	configCache invoicingRedis.IConfigCache
	log         log.ILogger
}

// New crea una nueva instancia del repositorio único
func New(database db.IDatabase, configCache invoicingRedis.IConfigCache, logger log.ILogger) ports.IRepository {
	return &Repository{
		db:          database,
		configCache: configCache,
		log:         logger.WithModule("invoicing.repository"),
	}
}
