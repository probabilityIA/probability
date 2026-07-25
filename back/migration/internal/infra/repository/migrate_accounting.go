package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateAccounting(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(
		&models.AccountingConcept{},
		&models.AccountingTax{},
		&models.AccountingConceptTax{},
		&models.AccountingEntry{},
	); err != nil {
		return fmt.Errorf("automigrate accounting: %w", err)
	}

	if err := db.Exec(`
CREATE UNIQUE INDEX IF NOT EXISTS idx_acc_entry_source
ON accounting_entries(source_type, source_id)
WHERE source_id <> '' AND deleted_at IS NULL
`).Error; err != nil {
		return fmt.Errorf("create accounting source index: %w", err)
	}

	concepts := []struct {
		code, name, description, kind, sourceType string
		isRealIncome, isAutomatic                 bool
	}{
		{"GUIDE_MARGIN", "Ganancia por guias", "Margen cobrado sobre guias de envio (incluye margen COD)", "INCOME", "GUIDE_MARGIN", true, true},
		{"SUBSCRIPTION", "Membresias", "Pagos de membresias y planes de negocios", "INCOME", "SUBSCRIPTION", true, true},
		{"WALLET_RECHARGE", "Recargas de billetera", "Dinero que los negocios recargan a su billetera (anticipo, no ganancia)", "INCOME", "WALLET_RECHARGE", false, true},
		{"COD_PAYOUT", "Consignaciones COD a clientes", "Dinero recaudado contra entrega que se consigna al negocio (salida de caja)", "EXPENSE", "COD_PAYOUT", false, true},
		{"COLLABORATOR", "Pagos a colaboradores", "Pagos a personas de apoyo (no nomina formal)", "EXPENSE", "MANUAL", true, false},
	}
	for _, c := range concepts {
		if err := db.Exec(`
INSERT INTO accounting_concepts (code, name, description, kind, is_real_income, is_automatic, source_type, is_active, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, TRUE, NOW(), NOW())
ON CONFLICT (code) DO NOTHING
`, c.code, c.name, c.description, c.kind, c.isRealIncome, c.isAutomatic, c.sourceType).Error; err != nil {
			return fmt.Errorf("seed accounting concept %s: %w", c.code, err)
		}
	}

	taxes := []struct {
		code, name, description string
		rate                    float64
	}{
		{"GMF_4X1000", "4 x 1000 (GMF)", "Gravamen a los movimientos financieros sobre debitos bancarios", 0.4},
		{"RETEICA", "ReteICA", "Retencion de industria y comercio (tarifa segun municipio)", 0.966},
		{"RETEFUENTE", "Retencion en la fuente", "Retencion en la fuente por honorarios/servicios", 11},
		{"IVA", "IVA", "Impuesto al valor agregado", 19},
	}
	for _, t := range taxes {
		if err := db.Exec(`
INSERT INTO accounting_taxes (code, name, description, rate_percent, is_active, created_at, updated_at)
VALUES (?, ?, ?, ?, TRUE, NOW(), NOW())
ON CONFLICT (code) DO NOTHING
`, t.code, t.name, t.description, t.rate).Error; err != nil {
			return fmt.Errorf("seed accounting tax %s: %w", t.code, err)
		}
	}

	if err := db.Exec(`
INSERT INTO accounting_concept_taxes (concept_id, tax_id, is_active, created_at, updated_at)
SELECT c.id, t.id, TRUE, NOW(), NOW()
FROM accounting_concepts c, accounting_taxes t
WHERE c.code = 'COD_PAYOUT' AND t.code = 'GMF_4X1000'
AND NOT EXISTS (
    SELECT 1 FROM accounting_concept_taxes ct
    WHERE ct.concept_id = c.id AND ct.tax_id = t.id
)
`).Error; err != nil {
		return fmt.Errorf("seed accounting concept tax: %w", err)
	}

	return nil
}
