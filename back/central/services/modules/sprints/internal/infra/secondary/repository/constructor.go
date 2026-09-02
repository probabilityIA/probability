package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
)

const (
	statusActive = "active"
	statusClosed = "closed"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}
