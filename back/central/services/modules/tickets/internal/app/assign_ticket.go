package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) Assign(ctx context.Context, dto dtos.AssignTicketDTO) (*entities.Ticket, error) {
	current, err := uc.repo.GetByID(ctx, dto.TicketID)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{}
	note := ""
	if dto.AssignedToID == nil {
		updates["assigned_to_id"] = nil
		note = "Ticket sin asignar"
	} else {
		exists, err := uc.repo.UserExists(ctx, *dto.AssignedToID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, dom.ErrAssigneeNotFound
		}
		updates["assigned_to_id"] = *dto.AssignedToID
		note = fmt.Sprintf("Asignado a usuario %d", *dto.AssignedToID)
	}

	updated, err := uc.repo.Update(ctx, dto.TicketID, updates)
	if err != nil {
		return nil, err
	}
	_ = uc.repo.AddHistory(ctx, dto.TicketID, current.Status, current.Status, dto.ChangedByID, note)
	return updated, nil
}
