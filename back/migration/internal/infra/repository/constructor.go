package repository

import (
	"context"

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
	if err := r.migrateCatalogPricing(ctx); err != nil {
		return err
	}
	if err := r.migrateCodReport(ctx); err != nil {
		return err
	}
	if err := r.migrateShippingMarginCOD(ctx); err != nil {
		return err
	}
	if err := r.migrateShipmentCodMargin(ctx); err != nil {
		return err
	}
	if err := r.migrateShipmentCodRefactor(ctx); err != nil {
		return err
	}
	if err := r.migrateShipmentProbabilityGuide(ctx); err != nil {
		return err
	}
	if err := r.migrateGuideFormats(ctx); err != nil {
		return err
	}
	if err := r.migrateShippingQuotes(ctx); err != nil {
		return err
	}
	if err := r.migrateInvoicePartialUniqueIndex(ctx); err != nil {
		return err
	}
	if err := r.backfillGeocodePendingOrders(ctx); err != nil {
		return err
	}
	if err := r.backfillOrdersGeozoneByPoint(ctx); err != nil {
		return err
	}
	return r.backfillOrdersGeozone(ctx)
}
