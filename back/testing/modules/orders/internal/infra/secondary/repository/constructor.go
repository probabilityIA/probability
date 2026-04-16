package repository

import (
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/testing/shared/db"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}
