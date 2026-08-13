package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *ProductRepository) SaveChannelSnapshots(ctx context.Context, businessID, integrationID uint, entries []productmatch.SnapshotEntry) error {
	if len(entries) == 0 {
		return nil
	}

	ahora := time.Now()
	conn := r.db.Conn(ctx)

	for _, entry := range entries {
		if entry.ProbabilityID == "" || entry.Snapshot.IsEmpty() {
			continue
		}
		raw, err := json.Marshal(entry.Snapshot)
		if err != nil {
			continue
		}

		query := conn.Model(&models.ProductBusinessIntegration{}).
			Where("product_id = ? AND business_id = ? AND integration_id = ? AND deleted_at IS NULL",
				entry.ProbabilityID, businessID, integrationID)
		if entry.ChannelVariantID != "" {
			query = query.Where("external_variant_id = ?", entry.ChannelVariantID)
		}

		if err := query.Updates(map[string]interface{}{
			"channel_snapshot": raw,
			"snapshot_at":      ahora,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
