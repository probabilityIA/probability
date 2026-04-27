package app

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) Create(ctx context.Context, dto dtos.CreateTicketDTO) (*entities.Ticket, error) {
	title := strings.TrimSpace(dto.Title)
	desc := strings.TrimSpace(dto.Description)
	if title == "" {
		return nil, dom.ErrTitleRequired
	}
	if desc == "" {
		return nil, dom.ErrDescriptionRequired
	}

	t := strings.ToLower(strings.TrimSpace(dto.Type))
	if t == "" {
		t = "support"
	}
	if !validTypes[t] {
		return nil, dom.ErrInvalidType
	}

	pr := strings.ToLower(strings.TrimSpace(dto.Priority))
	if pr == "" {
		pr = "medium"
	}
	if !validPriorities[pr] {
		return nil, dom.ErrInvalidPriority
	}

	sv := strings.ToLower(strings.TrimSpace(dto.Severity))
	if sv != "" && !validSeverities[sv] {
		return nil, dom.ErrInvalidSeverity
	}

	src := strings.ToLower(strings.TrimSpace(dto.Source))
	if src == "" {
		src = "internal"
	}
	if !validSources[src] {
		src = "internal"
	}

	area := strings.ToLower(strings.TrimSpace(dto.Area))
	if area == "" {
		area = "soporte"
	}
	if !validAreas[area] {
		return nil, dom.ErrInvalidArea
	}

	if dto.AssignedToID != nil {
		exists, err := uc.repo.UserExists(ctx, *dto.AssignedToID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, dom.ErrAssigneeNotFound
		}
	}

	code, err := uc.repo.NextCode(ctx)
	if err != nil {
		return nil, err
	}

	var due *time.Time
	if dto.DueDate != nil && *dto.DueDate != "" {
		if parsed, err := time.Parse("2006-01-02", *dto.DueDate); err == nil {
			due = &parsed
		}
	}

	ticket := &entities.Ticket{
		Code:         code,
		BusinessID:   dto.BusinessID,
		CreatedByID:  dto.CreatedByID,
		AssignedToID: dto.AssignedToID,
		Title:        title,
		Description:  desc,
		Type:         t,
		Category:     strings.TrimSpace(dto.Category),
		Priority:     pr,
		Status:       "open",
		Source:       src,
		Severity:     sv,
		Area:         area,
		DueDate:      due,
	}

	created, err := uc.repo.Create(ctx, ticket)
	if err != nil {
		return nil, err
	}

	_ = uc.repo.AddHistory(ctx, created.ID, "", "open", dto.CreatedByID, "Ticket creado")

	return uc.repo.GetByID(ctx, created.ID)
}
