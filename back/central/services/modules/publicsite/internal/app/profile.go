package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) resolveClient(ctx context.Context, slug string, userID uint) (uint, uint, error) {
	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return 0, 0, err
	}
	if business == nil {
		return 0, 0, domainerrors.ErrBusinessNotFound
	}
	session, err := uc.repo.GetClientSession(ctx, business.ID, userID)
	if err != nil {
		return 0, 0, err
	}
	if session == nil || session.CustomerID == 0 {
		return 0, 0, domainerrors.ErrBusinessNotFound
	}
	return business.ID, session.CustomerID, nil
}

func (uc *UseCase) UpdateProfile(ctx context.Context, slug string, userID uint, name, phone, dni string) error {
	businessID, _, err := uc.resolveClient(ctx, slug, userID)
	if err != nil {
		return err
	}
	return uc.repo.UpdateClientProfile(ctx, businessID, userID, name, phone, dni)
}

func (uc *UseCase) SaveAvatar(ctx context.Context, slug string, userID uint, avatarURL string) error {
	if _, _, err := uc.resolveClient(ctx, slug, userID); err != nil {
		return err
	}
	return uc.repo.UpdateUserAvatar(ctx, userID, avatarURL)
}

func (uc *UseCase) ListAddresses(ctx context.Context, slug string, userID uint) ([]entities.CustomerAddress, error) {
	businessID, customerID, err := uc.resolveClient(ctx, slug, userID)
	if err != nil {
		return nil, err
	}
	return uc.repo.ListCustomerAddresses(ctx, businessID, customerID)
}

func (uc *UseCase) AddAddress(ctx context.Context, slug string, userID uint, addr *entities.CustomerAddress) error {
	businessID, customerID, err := uc.resolveClient(ctx, slug, userID)
	if err != nil {
		return err
	}
	existing, err := uc.repo.ListCustomerAddresses(ctx, businessID, customerID)
	if err != nil {
		return err
	}
	if len(existing) == 0 {
		addr.IsPrimary = true
	}
	return uc.repo.CreateCustomerAddress(ctx, businessID, customerID, addr)
}

func (uc *UseCase) SetPrimaryAddress(ctx context.Context, slug string, userID uint, addressID uint) error {
	businessID, customerID, err := uc.resolveClient(ctx, slug, userID)
	if err != nil {
		return err
	}
	return uc.repo.SetPrimaryAddress(ctx, businessID, customerID, addressID)
}

func (uc *UseCase) DeleteAddress(ctx context.Context, slug string, userID uint, addressID uint) error {
	businessID, customerID, err := uc.resolveClient(ctx, slug, userID)
	if err != nil {
		return err
	}
	return uc.repo.DeleteCustomerAddress(ctx, businessID, customerID, addressID)
}
