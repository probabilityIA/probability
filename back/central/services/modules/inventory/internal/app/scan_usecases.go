package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (uc *useCase) Scan(ctx context.Context, dto request.ScanDTO) (*response.ScanResponse, error) {
	resolution, err := uc.repo.ResolveScanCode(ctx, dto.BusinessID, dto.Code)
	if err != nil {
		return nil, err
	}

	codeType := "unknown"
	if resolution != nil {
		codeType = resolution.CodeType
	}

	event := &entities.ScanEvent{
		BusinessID:  dto.BusinessID,
		UserID:      dto.UserID,
		DeviceID:    dto.DeviceID,
		ScannedCode: dto.Code,
		CodeType:    codeType,
		Action:      dto.Action,
		ScannedAt:   time.Now(),
	}
	saved, err := uc.repo.RecordScanEvent(ctx, event)
	if err != nil {
		saved = event
	}

	return &response.ScanResponse{
		Resolved:   resolution != nil,
		Resolution: resolution,
		Event:      saved,
	}, nil
}
