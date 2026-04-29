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

// Migrate ejecuta SOLO la migracion activa actual.
// Las migraciones pasadas ya estan aplicadas en produccion y no deben re-correr.
// Para una nueva migracion: dejar UNA sola llamada activa aqui, ejecutar, y volver a vaciar.
func (r *Repository) Migrate(ctx context.Context) error {
	return nil
}
