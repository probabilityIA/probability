package usecasebusiness

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

// generateCodeFromName genera un código único basado en el nombre del negocio
func generateCodeFromName(name string) string {
	// Normalizar el nombre: convertir a minúsculas, eliminar espacios y caracteres especiales
	normalized := strings.ToLower(name)
	var codeBuilder strings.Builder

	for _, char := range normalized {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			codeBuilder.WriteRune(char)
		} else if char == ' ' || char == '-' || char == '_' {
			codeBuilder.WriteRune('_')
		}
	}

	baseCode := codeBuilder.String()
	if len(baseCode) > 20 {
		baseCode = baseCode[:20]
	}

	// Agregar un sufijo aleatorio de 6 caracteres para garantizar unicidad
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err == nil {
		randomSuffix := base64.URLEncoding.EncodeToString(randomBytes)[:6]
		// Eliminar caracteres que puedan causar problemas
		randomSuffix = strings.ReplaceAll(randomSuffix, "-", "")
		randomSuffix = strings.ReplaceAll(randomSuffix, "_", "")
		return fmt.Sprintf("%s_%s", baseCode, randomSuffix)
	}

	// Fallback si no se puede generar aleatorio
	return fmt.Sprintf("%s_%d", baseCode, len(name))
}

// CreateBusiness crea un nuevo negocio
func (uc *BusinessUseCase) CreateBusiness(ctx context.Context, request domain.BusinessRequest) (*domain.BusinessResponse, error) {
	// Generar código automáticamente si no se proporciona
	businessCode := request.Code
	if businessCode == "" {
		businessCode = generateCodeFromName(request.Name)
		uc.log.Info().
			Str("name", request.Name).
			Str("generated_code", businessCode).
			Msg("Código generado automáticamente para el negocio")
	}

	uc.log.Info().
		Str("name", request.Name).
		Str("code", businessCode).
		Uint("business_type_id", request.BusinessTypeID).
		Msg("Creando negocio")

	// Validar que el tipo de negocio sea obligatorio y válido
	if request.BusinessTypeID == 0 {
		uc.log.Warn().Str("name", request.Name).Msg("Intento de crear negocio sin tipo de negocio")
		return nil, domain.ErrBusinessTypeIDRequired
	}

	// Validar que el tipo de negocio exista
	existingBusinessType, err := uc.repository.GetBusinessTypeByID(ctx, request.BusinessTypeID)
	if err != nil {
		uc.log.Error().Err(err).
			Uint("business_type_id", request.BusinessTypeID).
			Msg("Error al verificar si el tipo de negocio existe")
		return nil, domain.ErrBusinessTypeIDInvalid
	}
	if existingBusinessType == nil {
		uc.log.Warn().
			Uint("business_type_id", request.BusinessTypeID).
			Msg("El tipo de negocio especificado no existe")
		return nil, domain.ErrBusinessTypeIDInvalid
	}

	// Validar que el código no exista
	existing, err := uc.repository.GetBusinessByCode(ctx, businessCode)
	if err != nil && !errors.Is(err, domain.ErrBusinessNotFound) {
		uc.log.Error().Err(err).Str("code", businessCode).Msg("Error al verificar si el código del negocio ya existe")
		return nil, fmt.Errorf("error al verificar disponibilidad del código: %w", err)
	}

	if existing != nil {
		uc.log.Warn().Str("code", businessCode).Msg("El código del negocio ya está en uso")
		return nil, domain.ErrBusinessCodeAlreadyExists
	}

	// Validar que el dominio personalizado no exista si se proporciona
	if request.CustomDomain != "" {
		domainExists, err := uc.repository.GetBusinessByCustomDomain(ctx, request.CustomDomain)
		if err != nil && !errors.Is(err, domain.ErrBusinessNotFound) {
			uc.log.Error().Err(err).Str("domain", request.CustomDomain).Msg("Error al verificar si el dominio personalizado ya existe")
			return nil, fmt.Errorf("error al verificar disponibilidad del dominio personalizado: %w", err)
		}

		if domainExists != nil {
			uc.log.Warn().Str("domain", request.CustomDomain).Msg("El dominio personalizado ya está en uso")
			return nil, domain.ErrBusinessDomainAlreadyExists
		}
	}

	// Subir logo si viene archivo
	logoURL := ""
	if request.LogoFile != nil {
		uc.log.Info().Str("filename", request.LogoFile.Filename).Msg("Subiendo logo del negocio a S3")
		path, err := uc.s3.UploadImage(ctx, request.LogoFile, "businessLogo")
		if err != nil {
			uc.log.Error().Err(err).Str("filename", request.LogoFile.Filename).Msg("Error al subir logo del negocio a S3")
			return nil, fmt.Errorf("error al subir el logo del negocio: %w", err)
		}
		logoURL = path // Guardar solo path relativo
	}

	// Subir imagen de navbar si viene archivo
	navbarImageURL := ""
	if request.NavbarImageFile != nil {
		uc.log.Info().Str("filename", request.NavbarImageFile.Filename).Msg("Subiendo imagen de navbar a S3")
		path, err := uc.s3.UploadImage(ctx, request.NavbarImageFile, "navbar")
		if err != nil {
			uc.log.Error().Err(err).Str("filename", request.NavbarImageFile.Filename).Msg("Error al subir imagen de navbar a S3")
			return nil, fmt.Errorf("error al subir la imagen de navbar: %w", err)
		}
		navbarImageURL = path
	}

	// Crear entidad
	business := domain.Business{
		Name:               request.Name,
		Code:               businessCode, // Usar el código generado o proporcionado
		BusinessTypeID:     request.BusinessTypeID,
		Timezone:           request.Timezone,
		Address:            request.Address,
		Description:        request.Description,
		LogoURL:            logoURL,
		PrimaryColor:       request.PrimaryColor,
		SecondaryColor:     request.SecondaryColor,
		TertiaryColor:      request.TertiaryColor,
		QuaternaryColor:    request.QuaternaryColor,
		NavbarImageURL:     navbarImageURL,
		CustomDomain:       request.CustomDomain,
		IsActive:           request.IsActive,
		EnableDelivery:     request.EnableDelivery,
		EnablePickup:       request.EnablePickup,
		EnableReservations: request.EnableReservations,
	}

	// Guardar en repositorio (esto ahora crea también las relaciones con recursos)
	businessID, err := uc.repository.CreateBusiness(ctx, business)
	if err != nil {
		uc.log.Error().Err(err).
			Str("name", request.Name).
			Str("code", businessCode).
			Msg("Error al guardar el negocio en la base de datos")
		return nil, fmt.Errorf("error al guardar el negocio en la base de datos: %w", err)
	}

	// Obtener el negocio creado
	created, err := uc.repository.GetBusinessByID(ctx, businessID)
	if err != nil {
		uc.log.Error().Err(err).
			Uint("id", businessID).
			Str("name", request.Name).
			Msg("Error al obtener el negocio recién creado")
		return nil, fmt.Errorf("error al obtener el negocio recién creado: %w", err)
	}

	// Completar URL de logo si es path relativo
	fullLogoURL := created.LogoURL
	if fullLogoURL != "" && !strings.HasPrefix(fullLogoURL, "http") {
		base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
		if base != "" {
			fullLogoURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(fullLogoURL, "/"))
		}
	}
	// Completar URL de imagen de navbar si es path relativo
	fullNavbarImageURL := created.NavbarImageURL
	if fullNavbarImageURL != "" && !strings.HasPrefix(fullNavbarImageURL, "http") {
		base := strings.TrimRight(uc.env.Get("URL_BASE_DOMAIN_S3"), "/")
		if base != "" {
			fullNavbarImageURL = fmt.Sprintf("%s/%s", base, strings.TrimLeft(fullNavbarImageURL, "/"))
		}
	}

	// Mapear BusinessType
	businessType := domain.BusinessTypeResponse{
		ID: created.BusinessTypeID,
	}
	if created.BusinessType != nil {
		businessType = domain.BusinessTypeResponse{
			ID:          created.BusinessType.ID,
			Name:        created.BusinessType.Name,
			Code:        created.BusinessType.Code,
			Description: created.BusinessType.Description,
			Icon:        created.BusinessType.Icon,
			IsActive:    created.BusinessType.IsActive,
			CreatedAt:   created.BusinessType.CreatedAt,
			UpdatedAt:   created.BusinessType.UpdatedAt,
		}
	}

	response := &domain.BusinessResponse{
		ID:                 created.ID,
		Name:               created.Name,
		Code:               created.Code,
		BusinessType:       businessType,
		Timezone:           created.Timezone,
		Address:            created.Address,
		Description:        created.Description,
		LogoURL:            fullLogoURL,
		PrimaryColor:       created.PrimaryColor,
		SecondaryColor:     created.SecondaryColor,
		TertiaryColor:      created.TertiaryColor,
		QuaternaryColor:    created.QuaternaryColor,
		NavbarImageURL:     fullNavbarImageURL,
		CustomDomain:       created.CustomDomain,
		IsActive:           created.IsActive,
		EnableDelivery:     created.EnableDelivery,
		EnablePickup:       created.EnablePickup,
		EnableReservations: created.EnableReservations,
		CreatedAt:          created.CreatedAt,
		UpdatedAt:          created.UpdatedAt,
	}

	uc.log.Info().Uint("id", businessID).Str("name", request.Name).Msg("Negocio creado exitosamente")
	return response, nil
}
