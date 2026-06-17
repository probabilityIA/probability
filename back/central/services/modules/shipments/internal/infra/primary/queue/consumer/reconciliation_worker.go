package consumer

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	reconciliationInterval = 10 * time.Minute
	reconciliationLookback = 24 * time.Hour
	reconciliationGrace    = 5 * time.Minute
	reconciliationBatch    = 200
)

type WalletReconciliationWorker struct {
	repo domain.IRepository
	log  log.ILogger
}

func NewWalletReconciliationWorker(repo domain.IRepository, logger log.ILogger) *WalletReconciliationWorker {
	return &WalletReconciliationWorker{
		repo: repo,
		log:  logger.WithModule("shipments.wallet_reconciliation"),
	}
}

func (w *WalletReconciliationWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(reconciliationInterval)
	defer ticker.Stop()

	w.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *WalletReconciliationWorker) runOnce(ctx context.Context) {
	now := time.Now()
	createdAfter := now.Add(-reconciliationLookback)
	createdBefore := now.Add(-reconciliationGrace)

	guides, err := w.repo.FindUnchargedGuides(ctx, createdAfter, createdBefore, reconciliationBatch)
	if err != nil {
		w.log.Error(ctx).Err(err).Msg("Wallet reconciliation: failed to query uncharged guides")
		return
	}
	if len(guides) == 0 {
		return
	}

	w.log.Warn(ctx).Int("count", len(guides)).Msg("Wallet reconciliation: uncharged guides found, charging now")

	recovered := 0
	for _, g := range guides {
		shipmentID := g.ShipmentID
		if err := w.repo.DebitWalletForGuide(ctx, g.BusinessID, g.TotalCost, g.TrackingNumber, &shipmentID); err != nil {
			w.log.Error(ctx).Err(err).
				Uint("business_id", g.BusinessID).
				Uint("shipment_id", g.ShipmentID).
				Float64("amount", g.TotalCost).
				Str("tracking_number", g.TrackingNumber).
				Msg("Wallet reconciliation: failed to debit guide")
			continue
		}
		recovered++
		w.log.Info(ctx).
			Uint("business_id", g.BusinessID).
			Uint("shipment_id", g.ShipmentID).
			Float64("amount", g.TotalCost).
			Str("tracking_number", g.TrackingNumber).
			Msg("Wallet reconciliation: guide charged")
	}

	w.log.Warn(ctx).
		Int("recovered", recovered).
		Int("total", len(guides)).
		Msg("Wallet reconciliation: cycle complete")
}
