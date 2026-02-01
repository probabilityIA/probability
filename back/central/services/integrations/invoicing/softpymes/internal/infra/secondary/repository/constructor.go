package repository

import (
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Repository agrupa configuración base de los repositorios
type Repository struct {
	db  db.IDatabase
	log log.ILogger
}

// Repositories contiene todos los repositorios del módulo Softpymes
type Repositories struct {
	Provider     ports.IProviderRepository
	ProviderType ports.IProviderTypeRepository
}

// New crea una nueva instancia de todos los repositorios
func New(database db.IDatabase, logger log.ILogger) *Repositories {
	baseRepo := &Repository{
		db:  database,
		log: logger.WithModule("softpymes.repository"),
	}

	return &Repositories{
		Provider:     NewProviderRepository(baseRepo),
		ProviderType: NewProviderTypeRepository(baseRepo),
	}
}
