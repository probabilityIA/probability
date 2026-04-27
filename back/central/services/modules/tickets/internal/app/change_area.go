package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) ChangeArea(ctx context.Context, dto dtos.ChangeAreaDTO) (*entities.Ticket, error) {
	newArea := strings.ToLower(strings.TrimSpace(dto.NewArea))
	if !validAreas[newArea] {
		return nil, dom.ErrInvalidArea
	}

	current, err := uc.repo.GetByID(ctx, dto.TicketID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, dom.ErrTicketNotFound
	}

	if current.Area == newArea {
		return current, nil
	}

	updated, err := uc.repo.Update(ctx, dto.TicketID, map[string]any{"area": newArea})
	if err != nil {
		return nil, err
	}

	if err := uc.repo.AddAreaHistory(ctx, dto.TicketID, current.Area, newArea, dto.ChangedByID, dto.Note); err != nil {
		uc.log.Warn().Err(err).Uint("ticket_id", dto.TicketID).Msg("change_area: failed to record history")
	}

	return updated, nil
}
