package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) GetSession(ctx context.Context, slug string, userID uint) (*entities.TiendaSession, error) {
	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, domainerrors.ErrBusinessNotFound
	}

	belongs, err := uc.repo.UserBelongsToBusiness(ctx, business.ID, userID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, domainerrors.ErrBusinessNotFound
	}

	avatar, _ := uc.repo.GetUserAvatar(ctx, userID)

	session, err := uc.repo.GetClientSession(ctx, business.ID, userID)
	if err != nil {
		return nil, err
	}
	if session != nil {
		session.AvatarURL = avatar
		return session, nil
	}

	session, err = uc.repo.GetUserBasic(ctx, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, domainerrors.ErrBusinessNotFound
	}
	session.AvatarURL = avatar
	return session, nil
}
