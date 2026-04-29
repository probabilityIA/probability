package repository

import (
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) domain.IWebhookLogRepository {
	return &Repository{db: database}
}

func NewSyncLog(database db.IDatabase) domain.ISyncLogRepository {
	return &Repository{db: database}
}
