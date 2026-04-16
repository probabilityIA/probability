package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// repository implementa ports.IRepository
type repository struct {
	db  db.IDatabase
	log log.ILogger
}

// New crea una nueva instancia del repositorio
// Retorna la interfaz del dominio (inversión de dependencias)
func New(database db.IDatabase, logger log.ILogger) ports.IRepository {
	return &repository{
		db:  database,
		log: logger.WithModule("orderstatus-repository"),
	}
}
