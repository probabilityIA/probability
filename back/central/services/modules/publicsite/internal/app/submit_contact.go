package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) SubmitContact(ctx context.Context, slug string, dto *dtos.ContactFormDTO) error {
	if dto.Name == "" || dto.Message == "" {
		return domainerrors.ErrInvalidContact
	}

	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if business == nil {
		return domainerrors.ErrBusinessNotFound
	}

	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, business.ID, tiendaWebIntegrationTypeID)
	if err != nil {
		return err
	}
	if !active {
		return domainerrors.ErrPublicSiteNotActive
	}

	return uc.repo.SaveContactSubmission(ctx, business.ID, dto)
}
