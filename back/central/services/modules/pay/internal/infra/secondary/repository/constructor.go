package repository

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Repository struct {
	db              db.IDatabase
	log             log.ILogger
	integrationCore core.IIntegrationCore
}

func New(database db.IDatabase, logger log.ILogger, integrationCore core.IIntegrationCore) ports.IRepository {
	return &Repository{
		db:              database,
		log:             logger.WithModule("pay.repository"),
		integrationCore: integrationCore,
	}
}
