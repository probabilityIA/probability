package repository

import (
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type repository struct {
	db     db.IDatabase
	logger log.ILogger
}

func New(database db.IDatabase, logger log.ILogger) domain.IRepository {
	return &repository{
		db:     database,
		logger: logger.WithModule("syncruns"),
	}
}
