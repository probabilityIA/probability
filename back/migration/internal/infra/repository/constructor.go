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
	return r.backfillOrdersGeozone(ctx)
}
