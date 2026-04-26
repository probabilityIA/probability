package app

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) Update(ctx context.Context, dto dtos.UpdateTicketDTO) (*entities.Ticket, error) {
	updates := map[string]any{}

	if dto.Title != nil {
		v := strings.TrimSpace(*dto.Title)
		if v == "" {
			return nil, dom.ErrTitleRequired
		}
		updates["title"] = v
	}
	if dto.Description != nil {
		v := strings.TrimSpace(*dto.Description)
		if v == "" {
			return nil, dom.ErrDescriptionRequired
		}
		updates["description"] = v
	}
	if dto.Type != nil {
		v := strings.ToLower(strings.TrimSpace(*dto.Type))
		if !validTypes[v] {
			return nil, dom.ErrInvalidType
		}
		updates["type"] = v
	}
	if dto.Category != nil {
		updates["category"] = strings.TrimSpace(*dto.Category)
	}
	if dto.Priority != nil {
		v := strings.ToLower(strings.TrimSpace(*dto.Priority))
		if !validPriorities[v] {
			return nil, dom.ErrInvalidPriority
		}
		updates["priority"] = v
	}
	if dto.Severity != nil {
		v := strings.ToLower(strings.TrimSpace(*dto.Severity))
		if v != "" && !validSeverities[v] {
			return nil, dom.ErrInvalidSeverity
		}
		updates["severity"] = v
	}
	if dto.AssignedToID != nil {
		exists, err := uc.repo.UserExists(ctx, *dto.AssignedToID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, dom.ErrAssigneeNotFound
		}
		updates["assigned_to_id"] = *dto.AssignedToID
	}
	if dto.ClearDueDate {
		updates["due_date"] = nil
	} else if dto.DueDate != nil && *dto.DueDate != "" {
		if parsed, err := time.Parse("2006-01-02", *dto.DueDate); err == nil {
			updates["due_date"] = parsed
		}
	}

	if len(updates) == 0 {
		return uc.repo.GetByID(ctx, dto.ID)
	}

	return uc.repo.Update(ctx, dto.ID, updates)
}
