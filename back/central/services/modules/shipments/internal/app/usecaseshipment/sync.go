package usecaseshipment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type SyncShipmentsUseCase struct {
	repo         domain.IRepository
	transportPub domain.ITransportRequestPublisher
}

func NewSyncShipments(repo domain.IRepository, pub domain.ITransportRequestPublisher) *SyncShipmentsUseCase {
	return &SyncShipmentsUseCase{repo: repo, transportPub: pub}
}

func (uc *SyncShipmentsUseCase) SyncShipments(ctx context.Context, filter domain.SyncShipmentsFilter) (*domain.SyncShipmentsResult, error) {
	if filter.Provider == "" {
		filter.Provider = domain.SyncProviderEnvioclick
	}
	if filter.Provider != domain.SyncProviderEnvioclick {
		return nil, fmt.Errorf("provider %s not supported yet", filter.Provider)
	}
	if filter.BusinessID == 0 {
		return nil, fmt.Errorf("business_id is required")
	}
	if len(filter.Statuses) == 0 {
		filter.Statuses = []string{"pending", "in_transit", "picked_up", "out_for_delivery", "on_hold"}
	}

	integrationID, baseURL, err := uc.repo.GetBusinessActiveIntegration(ctx, filter.BusinessID, filter.Provider)
	if err != nil {
		return nil, err
	}

	rows, err := uc.repo.ListShipmentsForSync(ctx, filter)
	if err != nil {
		return nil, err
	}

	correlationID := "sync-" + uuid.New().String()
	result := &domain.SyncShipmentsResult{
		CorrelationID:  correlationID,
		TotalShipments: len(rows),
		BatchSize:      domain.SyncBatchSize,
	}

	if len(rows) == 0 {
		return result, nil
	}

	batches := chunkRows(rows, domain.SyncBatchSize)
	result.Batches = len(batches)
	result.EstimatedDurationSeconds = len(batches) * 3

	integrationTypeID := uint(12)

	go func() {
		bgCtx := context.Background()
		for i, batch := range batches {
			items := make([]map[string]any, 0, len(batch))
			for _, r := range batch {
				m := map[string]any{
					"shipment_id":     r.ShipmentID,
					"tracking_number": r.TrackingNumber,
					"carrier":         r.Carrier,
				}
				if r.EnvioclickIDOrder != nil {
					m["envioclick_id_order"] = *r.EnvioclickIDOrder
				}
				items = append(items, m)
			}

			req := &domain.TransportRequestMessage{
				Provider:          filter.Provider,
				IntegrationTypeID: integrationTypeID,
				Operation:         "sync_batch",
				CorrelationID:     fmt.Sprintf("%s-batch-%d", correlationID, i+1),
				BusinessID:        filter.BusinessID,
				IntegrationID:     integrationID,
				BaseURL:           baseURL,
				Timestamp:         time.Now(),
				Payload: map[string]any{
					"items": items,
				},
			}

			if err := uc.transportPub.PublishTransportRequest(bgCtx, req); err != nil {
				continue
			}
		}
	}()

	return result, nil
}

func chunkRows(rows []domain.SyncShipmentRow, size int) [][]domain.SyncShipmentRow {
	var out [][]domain.SyncShipmentRow
	for i := 0; i < len(rows); i += size {
		end := i + size
		if end > len(rows) {
			end = len(rows)
		}
		out = append(out, rows[i:end])
	}
	return out
}
