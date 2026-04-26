package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
)

func (uc *UseCase) Escalate(ctx context.Context, dto dtos.EscalateTicketDTO) (*entities.Ticket, error) {
	current, err := uc.repo.GetByID(ctx, dto.TicketID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	updates := map[string]any{
		"escalated_to_dev": true,
		"escalated_at":     now,
	}
	if current.Status == "open" {
		updates["status"] = "in_review"
	}

	updated, err := uc.repo.Update(ctx, dto.TicketID, updates)
	if err != nil {
		return nil, err
	}

	note := dto.Note
	if note == "" {
		note = "Escalado a desarrollo"
	}
	_ = uc.repo.AddHistory(ctx, dto.TicketID, current.Status, updated.Status, dto.ChangedByID, note)
	return updated, nil
}
