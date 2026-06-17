package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateWalletTxBusinessID(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE transaction
    ADD COLUMN IF NOT EXISTS business_id BIGINT
`).Error; err != nil {
		return fmt.Errorf("add business_id column: %w", err)
	}

	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_transaction_business_id ON transaction(business_id)`).Error; err != nil {
		return fmt.Errorf("create business_id index: %w", err)
	}

	if err := db.Exec(`
UPDATE transaction t
SET business_id = w.business_id
FROM wallet w
WHERE t.business_id IS NULL
  AND t.wallet_id = w.id
`).Error; err != nil {
		return fmt.Errorf("backfill business_id from wallet: %w", err)
	}

	return nil
}
