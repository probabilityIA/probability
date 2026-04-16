package usecasebusiness

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// GetBusinessByID obtiene un negocio por ID
func (uc *BusinessUseCase) GetBusinessByID(ctx context.Context, id uint) (*domain.BusinessResponse, error) {
	uc.log.Info().Uint("id", id).Msg("Obteniendo negocio por ID")

	business, err := uc.repository.GetBusinessByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener negocio por ID")
		return nil, fmt.Errorf("error al obtener negocio: %w", err)
	}

	if business == nil {
		uc.log.Warn().Uint("id", id).Msg("Negocio no encontrado")
		return nil, fmt.Errorf("negocio no encontrado")
	}

	// Completar URL del logo si es relativo
	fullLogoURL := business.LogoURL
	if fullLogoURL != "" && !strings.HasPrefix(fullLogoURL, "http") {
		base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
		if base != "" {
			fullLogoURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(fullLogoURL, "/"))
		}
	}
	// Completar URL de imagen de navbar si es relativa
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

	response := &domain.BusinessResponse{
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

	uc.log.Info().Uint("id", id).Msg("Negocio obtenido exitosamente")
	return response, nil
}
