package usecaseintegrationtype

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// generateIntegrationTypeCodeFromName genera un código único basado en el nombre del tipo de integración
func generateIntegrationTypeCodeFromName(name string) string {
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

	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err == nil {
		randomSuffix := base64.URLEncoding.EncodeToString(randomBytes)[:6]
		randomSuffix = strings.ReplaceAll(randomSuffix, "-", "")
		randomSuffix = strings.ReplaceAll(randomSuffix, "_", "")
		return fmt.Sprintf("%s_%s", baseCode, randomSuffix)
	}

	return fmt.Sprintf("%s_%d", baseCode, len(name))
}

// CreateIntegrationType crea un nuevo tipo de integración
func (uc *integrationTypeUseCase) CreateIntegrationType(ctx context.Context, dto domain.CreateIntegrationTypeDTO) (*domain.IntegrationType, error) {
	ctx = log.WithFunctionCtx(ctx, "CreateIntegrationType")

	// Validar que el nombre no exista
	existing, err := uc.repo.GetIntegrationTypeByName(ctx, dto.Name)
	if err != nil && !strings.Contains(err.Error(), "no encontrado") {
		uc.log.Error(ctx).Err(err).Str("name", dto.Name).Msg("Error al verificar si el nombre del tipo de integración ya existe")
		return nil, fmt.Errorf("error al verificar disponibilidad del nombre: %w", err)
	}

	if existing != nil {
		uc.log.Warn(ctx).Str("name", dto.Name).Msg("El nombre del tipo de integración ya está en uso")
		return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationTypeNameExists, dto.Name)
	}

	// Generar código automáticamente si no se proporciona
	code := dto.Code
	if code == "" {
		code = generateIntegrationTypeCodeFromName(dto.Name)
		uc.log.Info(ctx).
			Str("name", dto.Name).
			Str("generated_code", code).
			Msg("Código generado automáticamente para el tipo de integración")
	} else {
		// Validar que el código no exista
		existingByCode, err := uc.repo.GetIntegrationTypeByCode(ctx, code)
		if err != nil && !strings.Contains(err.Error(), "no encontrado") {
			uc.log.Error(ctx).Err(err).Str("code", code).Msg("Error al verificar si el código del tipo de integración ya existe")
			return nil, fmt.Errorf("error al verificar disponibilidad del código: %w", err)
		}
		if existingByCode != nil {
			uc.log.Warn(ctx).Str("code", code).Msg("El código del tipo de integración ya está en uso")
			return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationTypeCodeExists, code)
		}
	}

	uc.log.Info(ctx).
		Str("name", dto.Name).
		Str("code", code).
		Uint("category_id", dto.CategoryID).
		Msg("Creando tipo de integración")

	// Procesar imagen si se proporciona
	imageURL := ""
	if dto.ImageFile != nil {
		uc.log.Info(ctx).Str("name", dto.Name).Msg("Subiendo imagen del tipo de integración a S3")

		// Subir imagen a S3 en la carpeta "integration-types"
		// Retorna el path relativo (ej: "integration-types/1234567890_logo.jpg")
		imagePath, err := uc.s3.UploadImage(ctx, dto.ImageFile, "integration-types")
		if err != nil {
			uc.log.Error(ctx).Err(err).Str("name", dto.Name).Msg("Error al subir imagen del tipo de integración")
			return nil, fmt.Errorf("%w: %v", domain.ErrIntegrationTypeImageUploadFailed, err)
		}

		// Guardar solo el path relativo en la base de datos
		imageURL = imagePath
		uc.log.Info(ctx).Str("name", dto.Name).Str("image_path", imagePath).Msg("Imagen del tipo de integración subida exitosamente")
	}

	integrationType := &domain.IntegrationType{
		Name:              dto.Name,
		Code:              code,
		Description:       dto.Description,
		Icon:              dto.Icon,
		ImageURL:          imageURL,
		CategoryID:        dto.CategoryID,
		IsActive:          dto.IsActive,
		ConfigSchema:      dto.ConfigSchema,
		CredentialsSchema: dto.CredentialsSchema,
	}

	if err := uc.repo.CreateIntegrationType(ctx, integrationType); err != nil {
		uc.log.Error(ctx).Err(err).
			Str("name", dto.Name).
			Str("code", code).
			Msg("Error al guardar el tipo de integración en la base de datos")
		return nil, fmt.Errorf("error al guardar el tipo de integración en la base de datos: %w", err)
	}

	uc.log.Info(ctx).
		Uint("id", integrationType.ID).
		Str("name", dto.Name).
		Msg("Tipo de integración creado exitosamente")

	return integrationType, nil
}
