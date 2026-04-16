package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Repository implementa ports.IRepository usando GORM
type Repository struct {
	db  db.IDatabase
	log log.ILogger
}

// New crea una nueva instancia del repositorio de pagos
func New(database db.IDatabase, logger log.ILogger) ports.IRepository {
	return &Repository{
		db:  database,
		log: logger.WithModule("pay.repository"),
	}
}
