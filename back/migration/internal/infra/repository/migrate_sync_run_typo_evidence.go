package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// migrateSyncRunTypoEvidence agrega la evidencia que sostiene cada sugerencia de
// correccion de SKU: el SKU del otro sistema, su nombre, las existencias de cada
// lado, el patron detectado y en cual de los dos conviene corregir. Sin esto el
// comparativo solo puede afirmar que algo se parece, no proponer que hacer.
func (r *Repository) migrateSyncRunTypoEvidence(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.IntegrationSyncRun{}); err != nil {
		return fmt.Errorf("failed to auto-migrate integration_sync_runs: %w", err)
	}
	if err := r.db.Conn(ctx).AutoMigrate(&models.IntegrationSyncRunItem{}); err != nil {
		return fmt.Errorf("failed to auto-migrate integration_sync_run_items: %w", err)
	}
	return nil
}
