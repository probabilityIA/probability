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

// Migrate ejecuta SOLO la migracion activa actual.
// Las migraciones pasadas ya estan aplicadas en produccion y no deben re-correr.
// Para una nueva migracion: dejar UNA sola llamada activa aqui, ejecutar, y volver a vaciar.
func (r *Repository) Migrate(ctx context.Context) error {
	if err := r.migrateTicketArea(ctx); err != nil {
		return fmt.Errorf("failed to migrate ticket area columns: %w", err)
	}
	return nil
}
