package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) AddComment(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error) {
	body := strings.TrimSpace(dto.Body)
	if body == "" {
		return nil, dom.ErrDescriptionRequired
	}
	if _, err := uc.repo.GetByID(ctx, dto.TicketID); err != nil {
		return nil, err
	}
	dto.Body = body
	return uc.repo.AddComment(ctx, dto)
}

func (uc *UseCase) ListComments(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error) {
	return uc.repo.ListComments(ctx, ticketID, includeInternal)
}
