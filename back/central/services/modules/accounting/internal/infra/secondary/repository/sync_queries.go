package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
)

type syncRow struct {
	SourceID    string
	BusinessID  *uint
	Amount      float64
	EntryDate   time.Time
	Description string
}

func (r *Repository) FindSyncCandidates(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
	var sql string
	switch sourceType {
	case entities.SourceGuideMargin:
		sql = `
SELECT s.id::text AS source_id,
       NULLIF(t.business_id, 0) AS business_id,
       COALESCE(s.applied_margin, 0) + COALESCE(s.cod_probability_margin, 0) AS amount,
       t.created_at AS entry_date,
       CONCAT('Margen guia ', COALESCE(s.tracking_number, '')) AS description
FROM transaction t
JOIN shipments s ON s.id = t.shipment_id
WHERE t.type = 'USAGE' AND t.status = 'COMPLETED' AND t.shipment_id IS NOT NULL
  AND s.deleted_at IS NULL AND s.is_test IS NOT TRUE
  AND COALESCE(s.applied_margin, 0) + COALESCE(s.cod_probability_margin, 0) > 0
  AND NOT EXISTS (
      SELECT 1 FROM accounting_entries ae
      WHERE ae.source_type = 'GUIDE_MARGIN' AND ae.source_id = s.id::text AND ae.deleted_at IS NULL
  )
ORDER BY t.created_at ASC
LIMIT ?`
	case entities.SourceSubscription:
		sql = `
SELECT bs.id::text AS source_id,
       NULLIF(bs.business_id, 0) AS business_id,
       bs.amount AS amount,
       COALESCE(bs.start_date, bs.created_at) AS entry_date,
       'Pago de membresia' AS description
FROM business_subscriptions bs
WHERE bs.status = 'paid' AND bs.deleted_at IS NULL AND bs.amount > 0
  AND NOT EXISTS (
      SELECT 1 FROM accounting_entries ae
      WHERE ae.source_type = 'SUBSCRIPTION' AND ae.source_id = bs.id::text AND ae.deleted_at IS NULL
  )
ORDER BY bs.created_at ASC
LIMIT ?`
	case entities.SourceWalletRecharge:
		sql = `
SELECT t.id::text AS source_id,
       NULLIF(t.business_id, 0) AS business_id,
       t.amount AS amount,
       t.created_at AS entry_date,
       CONCAT('Recarga billetera ', COALESCE(t.reference, '')) AS description
FROM transaction t
WHERE t.type = 'RECHARGE' AND t.status = 'COMPLETED' AND t.amount > 0
  AND NOT EXISTS (
      SELECT 1 FROM accounting_entries ae
      WHERE ae.source_type = 'WALLET_RECHARGE' AND ae.source_id = t.id::text AND ae.deleted_at IS NULL
  )
ORDER BY t.created_at ASC
LIMIT ?`
	case entities.SourceCodPayout:
		sql = `
SELECT c.id::text AS source_id,
       NULLIF(c.business_id, 0) AS business_id,
       c.total_net AS amount,
       COALESCE(c.confirmed_at, c.period_end) AS entry_date,
       CONCAT('Consignacion COD ', TO_CHAR(c.period_start, 'YYYY-MM-DD'), ' a ', TO_CHAR(c.period_end, 'YYYY-MM-DD')) AS description
FROM cod_payment_cut c
WHERE c.status = 'confirmed' AND c.deleted_at IS NULL AND c.total_net > 0
  AND NOT EXISTS (
      SELECT 1 FROM accounting_entries ae
      WHERE ae.source_type = 'COD_PAYOUT' AND ae.source_id = c.id::text AND ae.deleted_at IS NULL
  )
ORDER BY c.created_at ASC
LIMIT ?`
	default:
		return nil, nil
	}

	var rows []syncRow
	if err := r.db.Conn(ctx).Raw(sql, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]dtos.SyncCandidate, len(rows))
	for i, row := range rows {
		out[i] = dtos.SyncCandidate{
			SourceType:  sourceType,
			SourceID:    row.SourceID,
			BusinessID:  row.BusinessID,
			Amount:      row.Amount,
			EntryDate:   row.EntryDate,
			Description: row.Description,
		}
	}
	return out, nil
}
