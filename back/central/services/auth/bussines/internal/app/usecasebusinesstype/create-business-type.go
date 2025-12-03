package usecasebusinesstype

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

// generateCodeFromName genera un código único basado en el nombre del tipo de negocio
func generateBusinessTypeCodeFromName(name string) string {
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

// CreateBusinessType crea un nuevo tipo de negocio
func (uc *BusinessTypeUseCase) CreateBusinessType(ctx context.Context, request domain.BusinessTypeRequest) (*domain.BusinessTypeResponse, error) {
	// Validar que el nombre no exista
	existing, err := uc.repository.GetBusinessTypeByName(ctx, request.Name)
	if err != nil && !errors.Is(err, domain.ErrBusinessTypeNotFound) {
		uc.log.Error().Err(err).Str("name", request.Name).Msg("Error al verificar si el nombre del tipo de negocio ya existe")
		return nil, fmt.Errorf("error al verificar disponibilidad del nombre: %w", err)
	}

	if existing != nil {
		uc.log.Warn().Str("name", request.Name).Msg("El nombre del tipo de negocio ya está en uso")
		return nil, domain.ErrBusinessTypeNameAlreadyExists
	}

	// Generar código automáticamente si no se proporciona
	businessTypeCode := request.Code
	if businessTypeCode == "" {
		businessTypeCode = generateBusinessTypeCodeFromName(request.Name)
		uc.log.Info().
			Str("name", request.Name).
			Str("generated_code", businessTypeCode).
			Msg("Código generado automáticamente para el tipo de negocio")
	}

	uc.log.Info().
		Str("name", request.Name).
		Str("code", businessTypeCode).
		Msg("Creando tipo de negocio")

	// Crear entidad
	businessType := domain.BusinessType{
		Name:        request.Name,
		Code:        businessTypeCode,
		Description: request.Description,
		Icon:        request.Icon,
		IsActive:    request.IsActive,
	}

	// Guardar en repositorio
	_, err = uc.repository.CreateBusinessType(ctx, businessType)
	if err != nil {
		uc.log.Error().Err(err).
			Str("name", request.Name).
			Str("code", businessTypeCode).
			Msg("Error al guardar el tipo de negocio en la base de datos")
		return nil, fmt.Errorf("error al guardar el tipo de negocio en la base de datos: %w", err)
	}

	// Obtener el tipo de negocio creado por nombre (más confiable que por código)
	created, err := uc.repository.GetBusinessTypeByName(ctx, request.Name)
	if err != nil {
		uc.log.Error().Err(err).
			Str("name", request.Name).
			Str("code", businessTypeCode).
			Msg("Error al obtener el tipo de negocio recién creado")
		return nil, fmt.Errorf("error al obtener el tipo de negocio recién creado: %w", err)
	}

	response := &domain.BusinessTypeResponse{
		ID:          created.ID,
		Name:        created.Name,
		Code:        created.Code,
		Description: created.Description,
		Icon:        created.Icon,
		IsActive:    created.IsActive,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}

	uc.log.Info().Uint("id", created.ID).Str("name", request.Name).Msg("Tipo de negocio creado exitosamente")
	return response, nil
}
