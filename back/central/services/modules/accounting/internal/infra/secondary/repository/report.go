package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
)

func (r *Repository) Report(ctx context.Context, params dtos.ReportParams) ([]dtos.ReportConceptRow, []dtos.ReportTaxRow, error) {
	from := params.From.Format("2006-01-02")
	to := params.To.Format("2006-01-02")

	byConcept := []dtos.ReportConceptRow{}
	err := r.db.Conn(ctx).Raw(`
SELECT c.id AS concept_id, c.code, c.name, c.kind, c.is_real_income,
       COUNT(*) AS entries_count,
       COALESCE(SUM(e.amount), 0) AS amount,
       COALESCE(SUM(e.tax_total), 0) AS tax_total
FROM accounting_entries e
JOIN accounting_concepts c ON c.id = e.concept_id
WHERE e.deleted_at IS NULL AND e.entry_date >= ? AND e.entry_date <= ?
GROUP BY c.id, c.code, c.name, c.kind, c.is_real_income
ORDER BY c.kind ASC, amount DESC
`, from, to).Scan(&byConcept).Error
	if err != nil {
		return nil, nil, err
	}

	byTax := []dtos.ReportTaxRow{}
	err = r.db.Conn(ctx).Raw(`
SELECT d->>'tax_code' AS tax_code,
       MAX(d->>'tax_name') AS tax_name,
       COALESCE(SUM((d->>'amount')::numeric), 0) AS amount
FROM accounting_entries e, jsonb_array_elements(e.tax_detail) d
WHERE e.deleted_at IS NULL AND e.entry_date >= ? AND e.entry_date <= ?
GROUP BY d->>'tax_code'
ORDER BY amount DESC
`, from, to).Scan(&byTax).Error
	if err != nil {
		return nil, nil, err
	}

	return byConcept, byTax, nil
}
