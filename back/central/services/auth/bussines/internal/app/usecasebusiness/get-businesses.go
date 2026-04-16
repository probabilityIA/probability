package usecasebusiness

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// GetBusinesses obtiene todos los negocios con paginación
func (uc *BusinessUseCase) GetBusinesses(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.BusinessResponse, int64, error) {
	uc.log.Info().Int("page", page).Int("per_page", perPage).Str("name", name).Msg("Obteniendo negocios")

	businesses, total, err := uc.repository.GetBusinesses(ctx, page, perPage, name, businessTypeID, isActive)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error al obtener negocios")
		return nil, 0, fmt.Errorf("error al obtener negocios: %w", err)
	}

	// Convertir entidades a DTOs
	response := make([]domain.BusinessResponse, len(businesses))
	for i, business := range businesses {
		fullLogoURL := business.LogoURL
		if fullLogoURL != "" && !strings.HasPrefix(fullLogoURL, "http") {
			base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
			if base != "" {
				fullLogoURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(fullLogoURL, "/"))
			}
		}
		fullNavbarImageURL := business.NavbarImageURL
		if fullNavbarImageURL != "" && !strings.HasPrefix(fullNavbarImageURL, "http") {
			base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
			if base != "" {
				fullNavbarImageURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(fullNavbarImageURL, "/"))
			}
		}

		// Mapear BusinessType
		businessType := domain.BusinessTypeResponse{
			ID: business.BusinessTypeID,
		}
		if business.BusinessType != nil {
			businessType = domain.BusinessTypeResponse{
				ID:          business.BusinessType.ID,
				Name:        business.BusinessType.Name,
				Code:        business.BusinessType.Code,
				Description: business.BusinessType.Description,
				Icon:        business.BusinessType.Icon,
				IsActive:    business.BusinessType.IsActive,
				CreatedAt:   business.BusinessType.CreatedAt,
				UpdatedAt:   business.BusinessType.UpdatedAt,
			}
		}

		response[i] = domain.BusinessResponse{
			ID:                 business.ID,
			Name:               business.Name,
			Code:               business.Code,
			BusinessType:       businessType,
			Timezone:           business.Timezone,
			Address:            business.Address,
			Description:        business.Description,
			LogoURL:            fullLogoURL,
			PrimaryColor:       business.PrimaryColor,
			SecondaryColor:     business.SecondaryColor,
			TertiaryColor:      business.TertiaryColor,
			QuaternaryColor:    business.QuaternaryColor,
			NavbarImageURL:     fullNavbarImageURL,
			CustomDomain:       business.CustomDomain,
			IsActive:           business.IsActive,
			EnableDelivery:     business.EnableDelivery,
			EnablePickup:       business.EnablePickup,
			EnableReservations: business.EnableReservations,
			CreatedAt:          business.CreatedAt,
			UpdatedAt:          business.UpdatedAt,
		}
	}

	uc.log.Info().Int("count", len(response)).Int64("total", total).Msg("Negocios obtenidos exitosamente")
	return response, total, nil
}
