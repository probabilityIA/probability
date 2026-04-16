package usecaseoriginaddress

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type OriginAddressUseCase struct {
	repository domain.IRepository
}

func New(repo domain.IRepository) *OriginAddressUseCase {
	return &OriginAddressUseCase{
		repository: repo,
	}
}

func (uc *OriginAddressUseCase) Create(ctx context.Context, businessID uint, req domain.CreateOriginAddressRequest) (*domain.OriginAddress, error) {
	address := &domain.OriginAddress{
		BusinessID:   businessID,
		Alias:        req.Alias,
		Company:      req.Company,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		Street:       req.Street,
		Suburb:       req.Suburb,
		CityDaneCode: req.CityDaneCode,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		IsDefault:    req.IsDefault,
	}

	// If this is the first address, or is_default is true, handle default logic
	existing, _ := uc.repository.ListOriginAddressesByBusiness(ctx, businessID)
	if len(existing) == 0 {
		address.IsDefault = true
	}

	if err := uc.repository.CreateOriginAddress(ctx, address); err != nil {
		return nil, err
	}

	if address.IsDefault && len(existing) > 0 {
		if err := uc.repository.SetDefaultOriginAddress(ctx, businessID, address.ID); err != nil {
			return nil, err
		}
	}

	return address, nil
}

func (uc *OriginAddressUseCase) List(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
	return uc.repository.ListOriginAddressesByBusiness(ctx, businessID)
}

func (uc *OriginAddressUseCase) GetByID(ctx context.Context, id uint, businessID uint) (*domain.OriginAddress, error) {
	address, err := uc.repository.GetOriginAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if address.BusinessID != businessID {
		return nil, errors.New("no tienes permiso para ver esta dirección")
	}
	return address, nil
}

func (uc *OriginAddressUseCase) Update(ctx context.Context, id uint, businessID uint, req domain.UpdateOriginAddressRequest) (*domain.OriginAddress, error) {
	address, err := uc.repository.GetOriginAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if address.BusinessID != businessID {
		return nil, errors.New("no tienes permiso para actualizar esta dirección")
	}

	if req.Alias != nil {
		address.Alias = *req.Alias
	}
	if req.Company != nil {
		address.Company = *req.Company
	}
	if req.FirstName != nil {
		address.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		address.LastName = *req.LastName
	}
	if req.Email != nil {
		address.Email = *req.Email
	}
	if req.Phone != nil {
		address.Phone = *req.Phone
	}
	if req.Street != nil {
		address.Street = *req.Street
	}
	if req.Suburb != nil {
		address.Suburb = *req.Suburb
	}
	if req.CityDaneCode != nil {
		address.CityDaneCode = *req.CityDaneCode
	}
	if req.City != nil {
		address.City = *req.City
	}
	if req.State != nil {
		address.State = *req.State
	}
	if req.PostalCode != nil {
		address.PostalCode = *req.PostalCode
	}

	if err := uc.repository.UpdateOriginAddress(ctx, address); err != nil {
		return nil, err
	}

	if req.IsDefault != nil && *req.IsDefault {
		if err := uc.repository.SetDefaultOriginAddress(ctx, businessID, address.ID); err != nil {
			return nil, err
		}
		address.IsDefault = true
	}

	return address, nil
}

func (uc *OriginAddressUseCase) Delete(ctx context.Context, id uint, businessID uint) error {
	address, err := uc.repository.GetOriginAddressByID(ctx, id)
	if err != nil {
		return err
	}
	if address.BusinessID != businessID {
		return errors.New("no tienes permiso para eliminar esta dirección")
	}
	if address.IsDefault {
		return errors.New("no puedes eliminar la dirección por defecto")
	}
	return uc.repository.DeleteOriginAddress(ctx, id)
}
