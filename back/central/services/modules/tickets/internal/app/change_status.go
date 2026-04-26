package app

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) ChangeStatus(ctx context.Context, dto dtos.ChangeStatusDTO) (*entities.Ticket, error) {
	st := strings.ToLower(strings.TrimSpace(dto.NewStatus))
	if !validStatuses[st] {
		return nil, dom.ErrInvalidStatus
	}

	current, err := uc.repo.GetByID(ctx, dto.TicketID)
	if err != nil {
		return nil, err
	}
	if current.Status == st {
		return current, nil
	}

	updates := map[string]any{"status": st}
	now := time.Now()
	if st == "resolved" {
		updates["resolved_at"] = now
	}
	if st == "closed" {
		updates["closed_at"] = now
		if current.ResolvedAt == nil {
			updates["resolved_at"] = now
		}
	}

	updated, err := uc.repo.Update(ctx, dto.TicketID, updates)
	if err != nil {
		return nil, err
	}
	_ = uc.repo.AddHistory(ctx, dto.TicketID, current.Status, st, dto.ChangedByID, dto.Note)
	return updated, nil
}
