package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New crea una nueva instancia del repositorio
func New(database db.IDatabase, logger log.ILogger) ports.IRepository {
	return &repository{
		db:     database,
		logger: logger.WithModule("notification_config_repository"),
	}
}
