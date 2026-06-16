package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateWalletKPISelection(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
CREATE TABLE IF NOT EXISTS wallet_kpi_selection (
    id SERIAL PRIMARY KEY,
    selected_business_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
`).Error; err != nil {
		return fmt.Errorf("create wallet_kpi_selection table: %w", err)
	}

	if err := db.Exec(`INSERT INTO wallet_kpi_selection (selected_business_ids) SELECT '[]'::jsonb WHERE NOT EXISTS (SELECT 1 FROM wallet_kpi_selection)`).Error; err != nil {
		return fmt.Errorf("seed wallet_kpi_selection: %w", err)
	}

	return nil
}
