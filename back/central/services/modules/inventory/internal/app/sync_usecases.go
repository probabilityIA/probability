package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (uc *useCase) InboundSync(ctx context.Context, dto request.InboundSyncDTO) (*response.InboundSyncResult, error) {
	hash := computePayloadHash(dto.Payload)

	existing, err := uc.repo.GetSyncLogByHash(ctx, dto.BusinessID, "in", hash)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return &response.InboundSyncResult{Log: existing, Duplicate: true}, nil
	}

	now := time.Now()
	log := &entities.InventorySyncLog{
		BusinessID:    dto.BusinessID,
		IntegrationID: &dto.IntegrationID,
		Direction:     "in",
		PayloadHash:   hash,
		Status:        "received",
		SyncedAt:      &now,
	}
	created, err := uc.repo.CreateSyncLog(ctx, log)
	if err != nil {
		return nil, err
	}
	return &response.InboundSyncResult{Log: created, Duplicate: false}, nil
}

func (uc *useCase) ListSyncLogs(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListSyncLogs(ctx, params)
}

func computePayloadHash(payload map[string]any) string {
	data, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
