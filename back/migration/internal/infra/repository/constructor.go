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

	if err := r.migrateWarehouseHierarchy(ctx); err != nil {
		return fmt.Errorf("failed to migrate warehouse hierarchy: %w", err)
	}

	if err := r.migrateInventoryTraceability(ctx); err != nil {
		return fmt.Errorf("failed to migrate inventory traceability: %w", err)
	}

	if err := r.migrateInventoryOperations(ctx); err != nil {
		return fmt.Errorf("failed to migrate inventory operations: %w", err)
	}

	if err := r.migrateInventoryAudit(ctx); err != nil {
		return fmt.Errorf("failed to migrate inventory audit: %w", err)
	}

	if err := r.migrateInventoryCapture(ctx); err != nil {
		return fmt.Errorf("failed to migrate inventory capture: %w", err)
	}

	if err := r.migrateProductVariants(ctx); err != nil {
		return fmt.Errorf("failed to migrate product variants: %w", err)
	}

	if err := r.migrateBusinessOrderPrefix(ctx); err != nil {
		return fmt.Errorf("failed to migrate business order_prefix: %w", err)
	}

	return nil
}
