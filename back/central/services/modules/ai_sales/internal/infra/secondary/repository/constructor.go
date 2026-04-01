package repository

import (
	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type repository struct {
	db  db.IDatabase
	log log.ILogger
}

// New crea un nuevo repositorio de productos para el modulo ai_sales
func New(database db.IDatabase, logger log.ILogger) domain.IProductRepository {
	return &repository{
		db:  database,
		log: logger,
	}
}
