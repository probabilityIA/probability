package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) ChangeSprint(ctx context.Context, dto dtos.ChangeSprintDTO) (*entities.Ticket, error) {
	current, err := uc.repo.GetByID(ctx, dto.TicketID)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{}
	note := ""
	if dto.SprintID == nil {
		updates["sprint_id"] = nil
		note = "Ticket retirado del sprint"
	} else {
		name, found, err := uc.repo.FindSprintName(ctx, *dto.SprintID)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, dom.ErrSprintNotFound
		}
		updates["sprint_id"] = *dto.SprintID
		note = fmt.Sprintf("Ticket movido al sprint %s", name)
	}

	updated, err := uc.repo.Update(ctx, dto.TicketID, updates)
	if err != nil {
		return nil, err
	}
	_ = uc.repo.AddHistory(ctx, dto.TicketID, current.Status, current.Status, dto.ChangedByID, note)
	return updated, nil
}
