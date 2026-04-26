package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
)

func (uc *UseCase) Get(ctx context.Context, id uint, requesterUserID uint, requesterBusinessID *uint, isSuperAdmin bool) (*entities.Ticket, error) {
	t, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !isSuperAdmin {
		if t.BusinessID == nil || requesterBusinessID == nil || *t.BusinessID != *requesterBusinessID {
			return nil, dom.ErrForbidden
		}
	}
	return t, nil
}
