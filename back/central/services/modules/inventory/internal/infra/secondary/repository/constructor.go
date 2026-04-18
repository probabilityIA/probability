package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
)

type Repository struct {
	db    db.IDatabase
	cache IInventoryCache
}

func New(database db.IDatabase, cache IInventoryCache) ports.IRepository {
	return &Repository{db: database, cache: cache}
}
