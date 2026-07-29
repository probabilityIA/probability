package repository

import (
	"context"
	"fmt"
)

const defaultProductMatchRulesJSON = `[{"probability":"sku","channel":"sku"}]`

// Los codigos en integration_types no son homogeneos ("Shopify", "Mercado Libre"),
// asi que se matchea por id y por cualquier alias en minusculas.
type productMatchTypeDefaults struct {
	TypeID  uint
	Aliases []string
	Options string
	Rules   string
}

var productMatchDefaultsByType = []productMatchTypeDefaults{
	{
		TypeID:  1,
		Aliases: []string{"shopify"},
		Options: `{"probability":["sku","barcode","external_id"],"channel":["sku","barcode","external_id","variant_id"]}`,
		Rules:   defaultProductMatchRulesJSON,
	},
	{
		TypeID:  3,
		Aliases: []string{"mercado_libre", "mercado libre"},
		Options: `{"probability":["sku","barcode","external_id"],"channel":["sku","barcode","external_id"]}`,
		Rules:   defaultProductMatchRulesJSON,
	},
	{
		TypeID:  4,
		Aliases: []string{"woocommerce", "woocormerce"},
		Options: `{"probability":["sku","barcode","external_id"],"channel":["sku","barcode","external_id","variant_id"]}`,
		Rules:   defaultProductMatchRulesJSON,
	},
	{
		TypeID:  16,
		Aliases: []string{"vtex"},
		Options: `{"probability":["sku","barcode","external_id"],"channel":["sku","barcode","external_id","variant_id"]}`,
		Rules:   defaultProductMatchRulesJSON,
	},
	{
		TypeID:  33,
		Aliases: []string{"jumpseller"},
		Options: `{"probability":["sku","barcode","external_id"],"channel":["sku","barcode","external_id","variant_id"]}`,
		Rules:   defaultProductMatchRulesJSON,
	},
	{
		TypeID:  8,
		Aliases: []string{"siigo"},
		Options: `{"probability":["sku","barcode","name"],"channel":["sku","name"]}`,
		Rules:   defaultProductMatchRulesJSON,
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
			 WHERE id = ? OR lower(code) IN ?`,
			d.Options, d.Rules, d.TypeID, d.Aliases,
		).Error; err != nil {
			return fmt.Errorf("failed to seed product match defaults for type %d: %w", d.TypeID, err)
		}
	}

	return nil
}
