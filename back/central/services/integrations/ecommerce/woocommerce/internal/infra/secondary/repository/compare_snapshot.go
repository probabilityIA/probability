package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/shared/inventorycompare"
)

func (r *ProductRepository) SaveCompareSnapshot(
	ctx context.Context,
	businessID, integrationID uint,
	rows []inventorycompare.Row,
	checkedAt time.Time,
) error {
	return inventorycompare.Save(ctx, r.db.Conn(ctx), businessID, integrationID, rows, checkedAt)
}

func (r *ProductRepository) LoadCompareSnapshot(
	ctx context.Context,
	businessID, integrationID uint,
	opts inventorycompare.LoadOptions,
) (*inventorycompare.Page, error) {
	return inventorycompare.Load(ctx, r.db.Conn(ctx), businessID, integrationID, opts)
}
