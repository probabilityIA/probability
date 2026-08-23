package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type repository struct {
	db     db.IDatabase
	logger log.ILogger
}

func New(database db.IDatabase, logger log.ILogger) ports.IRepository {
	return &repository{
		db:     database,
		logger: logger.WithModule("orderscompare.repository"),
	}
}
