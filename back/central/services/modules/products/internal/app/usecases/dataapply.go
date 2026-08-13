package usecases

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func (uc *UseCases) ApplyChannelData(ctx context.Context, req domain.DataApplyRequest) (*domain.DataApplyResult, error) {
	if req.BusinessID == 0 || req.IntegrationID == 0 {
		return nil, domain.ErrInvalidProductData
	}
	req.Field = strings.TrimSpace(req.Field)
	if _, ok := productmatch.DataFieldByKey(req.Field); !ok {
		return nil, domain.ErrInvalidProductData
	}
	if req.Mode != domain.ModeOverwrite {
		req.Mode = domain.ModeFillEmpty
	}

	batchID := uuid.New().String()
	at := time.Now()

	uc.emitir(ctx, req.BusinessID, req.IntegrationID, EventChannelDataStarted, map[string]interface{}{
		"batch_id": batchID,
		"field":    req.Field,
	})

	go uc.aplicarEnSegundoPlano(req, batchID, at)

	return &domain.DataApplyResult{
		BatchID:   batchID,
		Field:     req.Field,
		AppliedAt: at,
	}, nil
}

func (uc *UseCases) aplicarEnSegundoPlano(req domain.DataApplyRequest, batchID string, at time.Time) {
	ctx := context.Background()

	avance := func(processed, total int) {
		uc.emitir(ctx, req.BusinessID, req.IntegrationID, EventChannelDataProgress, map[string]interface{}{
			"batch_id":  batchID,
			"field":     req.Field,
			"processed": processed,
			"total":     total,
		})
	}

	var applied int
	var err error
	if req.Field == domain.FieldVariant {
		applied, err = uc.repo.ApplyChannelVariant(ctx, req, batchID, at, avance)
	} else {
		applied, err = uc.repo.ApplyChannelField(ctx, req, batchID, at)
	}

	if err != nil {
		if uc.eventos != nil && uc.eventos.logger != nil {
			uc.eventos.logger.Error(ctx).Err(err).
				Uint("business_id", req.BusinessID).Str("field", req.Field).
				Msg("Error al aplicar los datos del canal")
		}
		uc.emitir(ctx, req.BusinessID, req.IntegrationID, EventChannelDataFailed, map[string]interface{}{
			"batch_id": batchID,
			"field":    req.Field,
			"applied":  applied,
			"error":    err.Error(),
		})
		return
	}

	uc.emitir(ctx, req.BusinessID, req.IntegrationID, EventChannelDataCompleted, map[string]interface{}{
		"batch_id": batchID,
		"field":    req.Field,
		"applied":  applied,
	})
}

func (uc *UseCases) UndoChannelData(ctx context.Context, businessID uint, batchID string, userID uint) (*domain.DataUndoResult, error) {
	batchID = strings.TrimSpace(batchID)
	if businessID == 0 || batchID == "" {
		return nil, domain.ErrInvalidProductData
	}

	reverted, err := uc.repo.UndoBatch(ctx, businessID, batchID, userID, time.Now())
	if err != nil {
		return nil, err
	}
	return &domain.DataUndoResult{BatchID: batchID, Reverted: reverted}, nil
}

func (uc *UseCases) ListDataBatches(ctx context.Context, businessID uint, limit int) ([]domain.DataBatch, error) {
	if businessID == 0 {
		return nil, domain.ErrInvalidProductData
	}
	return uc.repo.ListDataBatches(ctx, businessID, limit)
}
