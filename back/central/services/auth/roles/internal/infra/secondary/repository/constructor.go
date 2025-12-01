package repository

import (
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Repository struct {
	database db.IDatabase
	logger   log.ILogger
}

func New(db db.IDatabase, logger log.ILogger) domain.IRoleRepository {
	return &Repository{
		database: db,
		logger:   logger,
	}
}
