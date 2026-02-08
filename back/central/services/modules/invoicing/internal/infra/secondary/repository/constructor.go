package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Repository implementa LA interfaz IRepository con TODOS los métodos de persistencia del módulo
type Repository struct {
	db  db.IDatabase
	log log.ILogger
}

// New crea una nueva instancia del repositorio único
func New(database db.IDatabase, logger log.ILogger) ports.IRepository {
	return &Repository{
		db:  database,
		log: logger.WithModule("invoicing.repository")}
}
