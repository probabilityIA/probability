package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
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
	if err := r.migrateAnnouncements(ctx); err != nil {
		return fmt.Errorf("failed to migrate announcements: %w", err)
	}

	if err := r.migrateCustomerAddressCoords(ctx); err != nil {
		return fmt.Errorf("failed to migrate customer address coords: %w", err)
	}

	if err := r.migrateWebhookLogs(ctx); err != nil {
		return fmt.Errorf("failed to migrate webhook logs: %w", err)
	}

	return nil
}
