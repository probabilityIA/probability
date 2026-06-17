package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateWalletTxConcept(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE transaction
    ADD COLUMN IF NOT EXISTS concept VARCHAR(50)
`).Error; err != nil {
		return fmt.Errorf("add concept column: %w", err)
	}

	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_transaction_concept ON transaction(concept)`).Error; err != nil {
		return fmt.Errorf("create concept index: %w", err)
	}

	steps := []struct {
		value string
		where string
	}{
		{"GUIDE", "reference ILIKE '%Guide generation%'"},
		{"SUBSCRIPTION", "reference ILIKE '%suscrip%' OR reference ILIKE '%mensual%' OR reference ILIKE '%plan%'"},
		{"REFUND", "reference ILIKE '%reembolso%' OR reference ILIKE '%devoluc%'"},
		{"ADJUSTMENT", "reference ILIKE 'ADMIN_ADJ%'"},
		{"RECHARGE", "type = 'RECHARGE' OR reference ILIKE 'WLT%' OR reference ILIKE 'BOLD%' OR reference ILIKE 'MANUAL%' OR reference ILIKE '%recarga%'"},
		{"OTHER", "TRUE"},
	}

	for _, s := range steps {
		sql := fmt.Sprintf(`UPDATE transaction SET concept = '%s' WHERE (concept IS NULL OR concept = '') AND (%s)`, s.value, s.where)
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("backfill concept %s: %w", s.value, err)
		}
	}

	return nil
}
