package repository

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New crea una nueva instancia del repositorio de integraciones
func New(
	database db.IDatabase,
	logger log.ILogger,
	encryptionService domain.IEncryptionService,
) domain.IIntegrationRepository {
	return &integrationRepository{
		db:                database,
		log:               logger,
		encryptionService: encryptionService,
	}
}
