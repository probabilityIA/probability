package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
)

// Repository implements ports.IRepository
type Repository struct {
	db db.IDatabase
}

// New creates a new storefront repository
func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}
