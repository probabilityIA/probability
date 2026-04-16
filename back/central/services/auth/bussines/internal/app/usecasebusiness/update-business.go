package usecasebusiness

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// UpdateBusiness actualiza un negocio existente
func (uc *BusinessUseCase) UpdateBusiness(ctx context.Context, id uint, request domain.UpdateBusinessRequest) (*domain.BusinessResponse, error) {
	uc.log.Info().Uint("id", id).Msg("Actualizando negocio")

	// Verificar que existe
	existing, err := uc.repository.GetBusinessByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener negocio para actualizar")
		return nil, fmt.Errorf("error al obtener negocio: %w", err)
	}
	if existing == nil {
		uc.log.Warn().Uint("id", id).Msg("Negocio no encontrado para actualizar")
		return nil, fmt.Errorf("negocio no encontrado")
	}

	// Validar cambios únicos
	if request.Code != nil && *request.Code != existing.Code {
		codeExists, err := uc.repository.GetBusinessByCode(ctx, *request.Code)
		if err != nil && err.Error() != "negocio no encontrado" {
			uc.log.Error().Err(err).Str("code", *request.Code).Msg("Error al verificar código existente")
			return nil, fmt.Errorf("error al verificar código existente: %w", err)
		}
		if codeExists != nil {
			uc.log.Warn().Str("code", *request.Code).Msg("Código de negocio ya existe")
			return nil, fmt.Errorf("el código '%s' ya existe", *request.Code)
		}
	}

	if request.CustomDomain != nil && *request.CustomDomain != existing.CustomDomain {
		domainExists, err := uc.repository.GetBusinessByCustomDomain(ctx, *request.CustomDomain)
		if err != nil && err.Error() != "negocio no encontrado" {
			uc.log.Error().Err(err).Str("domain", *request.CustomDomain).Msg("Error al verificar dominio existente")
			return nil, fmt.Errorf("error al verificar dominio existente: %w", err)
		}
		if domainExists != nil {
			uc.log.Warn().Str("domain", *request.CustomDomain).Msg("Dominio personalizado ya existe")
			return nil, fmt.Errorf("el dominio '%s' ya existe", *request.CustomDomain)
		}
	}

	// Valores base desde existente
	name := existing.Name
	code := existing.Code
	businessTypeID := existing.BusinessTypeID
	timezone := existing.Timezone
	address := existing.Address
	description := existing.Description
	logoURL := existing.LogoURL
	primaryColor := existing.PrimaryColor
	secondaryColor := existing.SecondaryColor
	tertiaryColor := existing.TertiaryColor
	quaternaryColor := existing.QuaternaryColor
	navbarImageURL := existing.NavbarImageURL
	customDomain := existing.CustomDomain
	isActive := existing.IsActive
	enableDelivery := existing.EnableDelivery
	enablePickup := existing.EnablePickup
	enableReservations := existing.EnableReservations

	// Merge de campos
	if request.Name != nil {
		name = *request.Name
	}
	if request.Code != nil {
		code = *request.Code
	}
	if request.BusinessTypeID != nil {
		businessTypeID = *request.BusinessTypeID
	}
	if request.Timezone != nil {
		timezone = *request.Timezone
	}
	if request.Address != nil {
		address = *request.Address
	}
	if request.Description != nil {
		description = *request.Description
	}
	if request.PrimaryColor != nil {
		primaryColor = *request.PrimaryColor
	}
	if request.SecondaryColor != nil {
		secondaryColor = *request.SecondaryColor
	}
	if request.TertiaryColor != nil {
		tertiaryColor = *request.TertiaryColor
	}
	if request.QuaternaryColor != nil {
		quaternaryColor = *request.QuaternaryColor
	}
	if request.NavbarImageFile != nil {
		uc.log.Info().Uint("business_id", id).Str("filename", request.NavbarImageFile.Filename).Msg("Subiendo nueva imagen de navbar a S3")
		path, err := uc.s3.UploadImage(ctx, request.NavbarImageFile, "navbar")
		if err != nil {
			uc.log.Error().Err(err).Uint("business_id", id).Msg("Error al subir imagen de navbar a S3")
			return nil, fmt.Errorf("error al subir imagen de navbar: %w", err)
		}
		if existing.NavbarImageURL != "" && existing.NavbarImageURL != path && !strings.HasPrefix(existing.NavbarImageURL, "http") {
			if err := uc.s3.DeleteImage(ctx, existing.NavbarImageURL); err != nil {
				uc.log.Warn().Err(err).Str("old_navbar_image", existing.NavbarImageURL).Msg("No se pudo eliminar imagen de navbar anterior (no crítico)")
			}
		}
		navbarImageURL = path
	}
	if request.CustomDomain != nil {
		customDomain = *request.CustomDomain
	}
	if request.IsActive != nil {
		isActive = *request.IsActive
	}
	if request.EnableDelivery != nil {
		enableDelivery = *request.EnableDelivery
	}
	if request.EnablePickup != nil {
		enablePickup = *request.EnablePickup
	}
	if request.EnableReservations != nil {
		enableReservations = *request.EnableReservations
	}

	// Logo: subir si viene archivo; si cambia, borrar el anterior relativo.
	if request.LogoFile != nil {
		uc.log.Info().Uint("business_id", id).Str("filename", request.LogoFile.Filename).Msg("Subiendo nuevo logo de negocio a S3")
		path, err := uc.s3.UploadImage(ctx, request.LogoFile, "businessLogo")
		if err != nil {
			uc.log.Error().Err(err).Uint("business_id", id).Msg("Error al subir nuevo logo a S3")
			return nil, fmt.Errorf("error al subir logo: %w", err)
		}
		if existing.LogoURL != "" && existing.LogoURL != path && !strings.HasPrefix(existing.LogoURL, "http") {
			if err := uc.s3.DeleteImage(ctx, existing.LogoURL); err != nil {
				uc.log.Warn().Err(err).Str("old_logo", existing.LogoURL).Msg("No se pudo eliminar logo anterior (no crítico)")
			}
		}
		logoURL = path
	}

	// Actualizar entidad
	business := domain.Business{
		Name:               name,
		Code:               code,
		BusinessTypeID:     businessTypeID,
		Timezone:           timezone,
		Address:            address,
		Description:        description,
		LogoURL:            logoURL,
		PrimaryColor:       primaryColor,
		SecondaryColor:     secondaryColor,
		TertiaryColor:      tertiaryColor,
		QuaternaryColor:    quaternaryColor,
		NavbarImageURL:     navbarImageURL,
		CustomDomain:       customDomain,
		IsActive:           isActive,
		EnableDelivery:     enableDelivery,
		EnablePickup:       enablePickup,
		EnableReservations: enableReservations,
	}

	// Guardar en repositorio
	_, err = uc.repository.UpdateBusiness(ctx, id, business)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al actualizar negocio")
		return nil, fmt.Errorf("error al actualizar negocio: %w", err)
	}

	// Obtener el negocio actualizado
	updated, err := uc.repository.GetBusinessByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener negocio actualizado")
		return nil, fmt.Errorf("error al obtener negocio actualizado: %w", err)
	}

	// Completar URL de logo si es path relativo
	fullLogoURL := updated.LogoURL
	if fullLogoURL != "" && !strings.HasPrefix(fullLogoURL, "http") {
		base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
		if base != "" {
			fullLogoURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(fullLogoURL, "/"))
		}
	}
	// Completar URL de imagen de navbar si es path relativo
	fullNavbarImageURL := updated.NavbarImageURL
	if fullNavbarImageURL != "" && !strings.HasPrefix(fullNavbarImageURL, "http") {
		base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
		if base != "" {
			fullNavbarImageURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(fullNavbarImageURL, "/"))
		}
	}

	// Mapear BusinessType
	businessType := domain.BusinessTypeResponse{
		ID: updated.BusinessTypeID,
	}
	if updated.BusinessType != nil {
		businessType = domain.BusinessTypeResponse{
			ID:          updated.BusinessType.ID,
			Name:        updated.BusinessType.Name,
			Code:        updated.BusinessType.Code,
			Description: updated.BusinessType.Description,
			Icon:        updated.BusinessType.Icon,
			IsActive:    updated.BusinessType.IsActive,
			CreatedAt:   updated.BusinessType.CreatedAt,
			UpdatedAt:   updated.BusinessType.UpdatedAt,
		}
	}

	response := &domain.BusinessResponse{
		ID:                 updated.ID,
		Name:               updated.Name,
		Code:               updated.Code,
		BusinessType:       businessType,
		Timezone:           updated.Timezone,
		Address:            updated.Address,
		Description:        updated.Description,
		LogoURL:            fullLogoURL,
		PrimaryColor:       updated.PrimaryColor,
		SecondaryColor:     updated.SecondaryColor,
		TertiaryColor:      updated.TertiaryColor,
		QuaternaryColor:    updated.QuaternaryColor,
		NavbarImageURL:     fullNavbarImageURL,
		CustomDomain:       updated.CustomDomain,
		IsActive:           updated.IsActive,
		EnableDelivery:     updated.EnableDelivery,
		EnablePickup:       updated.EnablePickup,
		EnableReservations: updated.EnableReservations,
		CreatedAt:          updated.CreatedAt,
		UpdatedAt:          updated.UpdatedAt,
	}

	uc.log.Info().Uint("id", id).Str("code", code).Msg("Negocio actualizado exitosamente")
	return response, nil
}
