package repository

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type Repository struct {
	db db.IDatabase
}

func New(db db.IDatabase) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Migrate(ctx context.Context) error {
	return r.db.Conn(ctx).AutoMigrate(
		&models.BusinessType{},
		&models.Scope{},
		&models.Business{},
		&models.BusinessResourceConfigured{},
		&models.Resource{},
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.BusinessStaff{},
		&models.Client{},
		&models.Action{},
		&models.APIKey{},
	)
}
