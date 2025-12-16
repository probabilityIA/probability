package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
)

// generateRandomPassword genera una contraseña segura de longitud n
func generateRandomPassword(n int) (string, error) {
	const (
		lower    = "abcdefghijklmnopqrstuvwxyz"
		upper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits   = "0123456789"
		specials = "!@#$%^&*()-_=+[]{}|;:,.<>?"
		all      = lower + upper + digits + specials
	)
	if n < 8 {
		n = 8 // mínimo recomendado
	}
	password := make([]byte, n)
	// Garantizar al menos un carácter de cada tipo
	sets := []string{lower, upper, digits, specials}
	for i, set := range sets {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		if err != nil {
			return "", err
		}
		password[i] = set[idx.Int64()]
	}
	// Rellenar el resto
	for i := len(sets); i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", err
		}
		password[i] = all[idx.Int64()]
	}
	// Mezclar
	for i := range password {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}
	return string(password), nil
}

// CreateUser crea un nuevo usuario y retorna email y contraseña generada
func (uc *UserUseCase) CreateUser(ctx context.Context, userDTO domain.CreateUserDTO) (string, string, string, error) {
	// Normalizar email a minúsculas
	normalizedEmail := strings.ToLower(strings.TrimSpace(userDTO.Email))
	uc.log.Info().Str("email", normalizedEmail).Msg("Iniciando caso de uso: crear usuario")

	// Validar que el email no exista (buscar en minúsculas)
	existingUser, err := uc.repository.GetUserByEmail(ctx, normalizedEmail)
	if err != nil {
		uc.log.Error().Err(err).Str("email", normalizedEmail).Msg("Error al verificar email existente")
		return "", "", "", fmt.Errorf("error al verificar email: %w", err)
	}
	if existingUser != nil {
		uc.log.Error().Str("email", normalizedEmail).Msg("Email ya existe")
		return "", "", "", domain.ErrUserEmailExists
	}

	// Generar contraseña aleatoria
	generatedPassword, err := generateRandomPassword(12)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error al generar contraseña aleatoria")
		return "", "", "", fmt.Errorf("%w: %v", domain.ErrUserPasswordError, err)
	}
	uc.log.Info().Str("email", normalizedEmail).Msg("Contraseña aleatoria generada exitosamente")

	// Procesar imagen de avatar si se proporciona
	avatarURL := userDTO.AvatarURL
	if userDTO.AvatarFile != nil {
		uc.log.Info().Str("email", userDTO.Email).Msg("Subiendo imagen de avatar a S3")

		// Subir imagen a S3 en la carpeta "avatars"
		// Retorna el path relativo (ej: "avatars/1234567890_imagen.jpg")
		avatarPath, err := uc.s3.UploadImage(ctx, userDTO.AvatarFile, "avatars")
		if err != nil {
			uc.log.Error().Err(err).Str("email", userDTO.Email).Msg("Error al subir imagen de avatar")
			return "", "", "", fmt.Errorf("%w: %v", domain.ErrUserAvatarUploadFailed, err)
		}

		// Guardar solo el path relativo en la base de datos
		avatarURL = avatarPath
		uc.log.Info().Str("email", userDTO.Email).Str("avatar_path", avatarPath).Msg("Imagen de avatar subida exitosamente")
	}

	// Convertir DTO a entidad, usando la contraseña generada y email normalizado
	user := domain.UsersEntity{
		Name:      userDTO.Name,
		Email:     normalizedEmail,   // Siempre en minúsculas
		Password:  generatedPassword, // El repo la hashea
		Phone:     userDTO.Phone,
		AvatarURL: avatarURL, // URL relativa o completa según corresponda
		IsActive:  userDTO.IsActive,
		ScopeID:   userDTO.ScopeID, // Scope del usuario: platform (1) o business (2)
	}

	// Crear usuario y obtener el ID
	userID, err := uc.repository.CreateUser(ctx, user)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error al crear usuario desde el repositorio")
		// Verificar si es un error de duplicado de email
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "email") {
			return "", "", "", domain.ErrUserEmailExists
		}
		return "", "", "", fmt.Errorf("error al crear usuario: %w", err)
	}

	uc.log.Info().Uint("user_id", userID).Msg("Usuario creado exitosamente, asignando businesses")

	if len(userDTO.BusinessIDs) > 0 {
		uc.log.Info().Uint("user_id", userID).Any("business_ids", userDTO.BusinessIDs).Msg("Asignando businesses al usuario")
		if err := uc.repository.AssignBusinessesToUser(ctx, userID, userDTO.BusinessIDs); err != nil {
			uc.log.Error().Err(err).Uint("user_id", userID).Any("business_ids", userDTO.BusinessIDs).Msg("Error al asignar businesses al usuario")
			if strings.Contains(err.Error(), "algunos businesses no existen") {
				return "", "", "", domain.ErrBusinessesNotFound
			}
			return "", "", "", fmt.Errorf("error al asignar businesses: %w", err)
		}
		uc.log.Info().Uint("user_id", userID).Int("businesses_count", len(userDTO.BusinessIDs)).Msg("Businesses asignados exitosamente")
	}

	message := fmt.Sprintf("Usuario creado con ID: %d", userID)
	uc.log.Info().
		Uint("user_id", userID).
		Str("email", normalizedEmail).
		Bool("password_generated", generatedPassword != "").
		Msg("Usuario creado exitosamente - retornando email y contraseña generada")
	return normalizedEmail, generatedPassword, message, nil
}
