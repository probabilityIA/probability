package app

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/errors"
)

func (uc *UseCase) Update(ctx context.Context, dto dtos.UpdateSprintDTO) (*entities.Sprint, error) {
	current, err := uc.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{}

	if dto.Name != nil {
		name := strings.TrimSpace(*dto.Name)
		if name == "" {
			return nil, dom.ErrNameRequired
		}
		updates["name"] = name
	}
	if dto.Goal != nil {
		updates["goal"] = strings.TrimSpace(*dto.Goal)
	}

	start := current.StartDate
	end := current.EndDate
	if dto.StartDate != nil {
		parsed, ok := parseStartDate(*dto.StartDate)
		if !ok {
			return nil, dom.ErrInvalidStartDate
		}
		start = parsed
		updates["start_date"] = parsed
	}
	if dto.EndDate != nil {
		parsed, ok := parseEndDate(*dto.EndDate)
		if !ok {
			return nil, dom.ErrInvalidEndDate
		}
		end = parsed
		updates["end_date"] = parsed
	}
	if err := validateRange(start, end); err != nil {
		return nil, err
	}

	if dto.Status != nil {
		status := strings.ToLower(strings.TrimSpace(*dto.Status))
		if !validSprintStatuses[status] {
			return nil, dom.ErrInvalidStatus
		}
		updates["status"] = status
	}

	if len(updates) == 0 {
		return current, nil
	}
	return uc.repo.Update(ctx, dto.ID, updates)
}

func validateRange(start, end time.Time) error {
	if !end.After(start) {
		return dom.ErrInvalidDateRange
	}
	return nil
}
