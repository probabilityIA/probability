package repository

import (
	"context"
	"fmt"
)

const defaultProductMatchRulesJSON = `[{"probability":"sku","channel":"sku"}]`

type productMatchTypeDefaults struct {
	TypeCode string
	Options  string
	Rules    string
}

var productMatchDefaultsByType = []productMatchTypeDefaults{
	{
		TypeCode: "shopify",
		Options:  `{"probability":["sku","barcode","external_id"],"channel":["sku","external_id","variant_id"]}`,
		Rules:    defaultProductMatchRulesJSON,
	},
	{
		TypeCode: "mercado_libre",
		Options:  `{"probability":["sku","barcode","external_id"],"channel":["sku","external_id","variant_id"]}`,
		Rules:    defaultProductMatchRulesJSON,
	},
	{
		TypeCode: "woocommerce",
		Options:  `{"probability":["sku","barcode","external_id"],"channel":["sku","external_id","variant_id"]}`,
		Rules:    defaultProductMatchRulesJSON,
	},
	{
		TypeCode: "vtex",
		Options:  `{"probability":["sku","barcode","external_id"],"channel":["sku","barcode","external_id","variant_id"]}`,
		Rules:    defaultProductMatchRulesJSON,
	},
	{
		TypeCode: "jumpseller",
		Options:  `{"probability":["sku","barcode","external_id"],"channel":["sku","barcode","external_id","variant_id"]}`,
		Rules:    defaultProductMatchRulesJSON,
	},
	{
		TypeCode: "siigo",
		Options:  `{"probability":["sku","barcode","name"],"channel":["sku","name"]}`,
		Rules:    defaultProductMatchRulesJSON,
	},
}

func (r *Repository) migrateProductMatchRules(ctx context.Context) error {
	conn := r.db.Conn(ctx)

	if err := conn.Exec(
		`ALTER TABLE integration_types ADD COLUMN IF NOT EXISTS product_match_options jsonb`,
	).Error; err != nil {
		return fmt.Errorf("failed to add integration_types.product_match_options: %w", err)
	}
	if err := conn.Exec(
		`ALTER TABLE integration_types ADD COLUMN IF NOT EXISTS default_product_match_rules jsonb`,
	).Error; err != nil {
		return fmt.Errorf("failed to add integration_types.default_product_match_rules: %w", err)
	}
	if err := conn.Exec(
		`ALTER TABLE integrations ADD COLUMN IF NOT EXISTS product_match_rules jsonb`,
	).Error; err != nil {
		return fmt.Errorf("failed to add integrations.product_match_rules: %w", err)
	}

	for _, d := range productMatchDefaultsByType {
		if err := conn.Exec(
			`UPDATE integration_types
			 SET product_match_options = ?::jsonb,
			     default_product_match_rules = COALESCE(default_product_match_rules, ?::jsonb)
			 WHERE code = ?`,
			d.Options, d.Rules, d.TypeCode,
		).Error; err != nil {
			return fmt.Errorf("failed to seed product match defaults for %s: %w", d.TypeCode, err)
		}
	}

	return nil
}
