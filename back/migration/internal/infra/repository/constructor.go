package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type Repository struct {
	db  db.IDatabase
	cfg env.IConfig
}

func New(db db.IDatabase, cfg env.IConfig) *Repository {
	return &Repository{
		db:  db,
		cfg: cfg,
	}
}

func (r *Repository) Migrate(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.CustomerSummary{},
		&models.CustomerAddress{},
		&models.CustomerProductHistory{},
		&models.CustomerOrderItem{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate customer history tables: %w", err)
	}

	if err := r.seedCustomerHistory(ctx); err != nil {
		return fmt.Errorf("failed to seed customer history: %w", err)
	}

	return nil
}
