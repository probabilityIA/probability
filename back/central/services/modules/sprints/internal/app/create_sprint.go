package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/errors"
)

func (uc *UseCase) Create(ctx context.Context, dto dtos.CreateSprintDTO) (*entities.Sprint, error) {
	name := strings.TrimSpace(dto.Name)
	if name == "" {
		return nil, dom.ErrNameRequired
	}
	if dto.CreatedByID == 0 {
		return nil, dom.ErrCreatorRequired
	}

	start, ok := parseStartDate(dto.StartDate)
	if !ok {
		return nil, dom.ErrInvalidStartDate
	}
	end, ok := parseEndDate(dto.EndDate)
	if !ok {
		return nil, dom.ErrInvalidEndDate
	}
	if !end.After(start) {
		return nil, dom.ErrInvalidDateRange
	}

	status := strings.ToLower(strings.TrimSpace(dto.Status))
	if status == "" {
		status = StatusPlanned
	}
	if !validSprintStatuses[status] {
		return nil, dom.ErrInvalidStatus
	}

	sprint := &entities.Sprint{
		Name:        name,
		Goal:        strings.TrimSpace(dto.Goal),
		StartDate:   start,
		EndDate:     end,
		Status:      status,
		CreatedByID: dto.CreatedByID,
	}

	created, err := uc.repo.Create(ctx, sprint)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetByID(ctx, created.ID)
}
