package repository

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Repository struct {
	db                db.IDatabase
	log               log.ILogger
	encryptionService domain.IEncryptionService
}

// New crea una nueva instancia del repositorio unificado
func New(database db.IDatabase, logger log.ILogger, encryptionService domain.IEncryptionService) domain.IRepository {
	return &Repository{
		db:                database,
		log:               logger,
		encryptionService: encryptionService,
	}
}
