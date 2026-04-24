package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) CreateCrossDockLink(ctx context.Context, dto request.CreateCrossDockLinkDTO) (*entities.CrossDockLink, error) {
	if dto.Quantity <= 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}
	link := &entities.CrossDockLink{
		BusinessID:        dto.BusinessID,
		InboundShipmentID: dto.InboundShipmentID,
		OutboundOrderID:   dto.OutboundOrderID,
		ProductID:         dto.ProductID,
		Quantity:          dto.Quantity,
		Status:            "pending",
	}
	return uc.repo.CreateCrossDockLink(ctx, link)
}

func (uc *useCase) ListCrossDockLinks(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListCrossDockLinks(ctx, params)
}

func (uc *useCase) ExecuteCrossDock(ctx context.Context, dto request.ExecuteCrossDockDTO) (*entities.CrossDockLink, error) {
	existing, err := uc.repo.GetCrossDockLinkByID(ctx, dto.BusinessID, dto.LinkID)
	if err != nil {
		return nil, err
	}
	if existing.Status == "executed" {
		return nil, domainerrors.ErrCrossDockExecuted
	}
	now := time.Now()
	existing.Status = "executed"
	existing.ExecutedAt = &now
	return uc.repo.UpdateCrossDockLink(ctx, existing)
}
