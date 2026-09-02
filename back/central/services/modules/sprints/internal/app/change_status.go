package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/errors"
)

func (uc *UseCase) ChangeStatus(ctx context.Context, dto dtos.ChangeSprintStatusDTO) (*entities.Sprint, error) {
	status := strings.ToLower(strings.TrimSpace(dto.Status))
	if !validSprintStatuses[status] {
		return nil, dom.ErrInvalidStatus
	}

	current, err := uc.repo.GetByID(ctx, dto.SprintID)
	if err != nil {
		return nil, err
	}
	if current.Status == status {
		return current, nil
	}

	return uc.repo.Update(ctx, dto.SprintID, map[string]any{"status": status})
}
