package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
)

// UpdateUser actualiza un usuario existente
func (uc *UserUseCase) UpdateUser(ctx context.Context, id uint, userDTO domain.UpdateUserDTO) (string, error) {
	uc.log.Info().Uint("id", id).Msg("Iniciando caso de uso: actualizar usuario")

	// Verificar que el usuario existe
	existingUser, err := uc.repository.GetUserByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al buscar usuario")
		return "", fmt.Errorf("error al buscar usuario: %w", err)
	}
	if existingUser == nil {
		uc.log.Error().Uint("id", id).Msg("Usuario no encontrado")
		return "", domain.ErrUserNotFound
	}

	// Normalizar email a minúsculas si se proporciona
	normalizedEmail := userDTO.Email
	if normalizedEmail != "" {
		normalizedEmail = strings.ToLower(strings.TrimSpace(normalizedEmail))
	}

	// Verificar que el email no esté en uso por otro usuario
	if normalizedEmail != "" && normalizedEmail != strings.ToLower(existingUser.Email) {
		userWithEmail, err := uc.repository.GetUserByEmail(ctx, normalizedEmail)
		if err != nil {
			uc.log.Error().Err(err).Str("email", normalizedEmail).Msg("Error al verificar email existente")
			return "", fmt.Errorf("error al verificar email: %w", err)
		}
		if userWithEmail != nil && userWithEmail.ID != id {
			uc.log.Error().Str("email", normalizedEmail).Msg("Email ya existe en otro usuario")
			return "", domain.ErrUserEmailExists
		}
	}

	// Procesar imagen de avatar si se proporciona una nueva
	avatarURL := userDTO.AvatarURL
	if userDTO.AvatarFile != nil {
		uc.log.Info().Uint("user_id", id).Msg("Subiendo nueva imagen de avatar a S3")

		// Subir nueva imagen a S3 en la carpeta "avatars"
		// Retorna el path relativo (ej: "avatars/1234567890_imagen.jpg")
		avatarPath, err := uc.s3.UploadImage(ctx, userDTO.AvatarFile, "avatars")
		if err != nil {
			uc.log.Error().Err(err).Uint("user_id", id).Msg("Error al subir nueva imagen de avatar")
			return "", fmt.Errorf("%w: %v", domain.ErrUserAvatarUploadFailed, err)
		}

		// Guardar solo el path relativo en la base de datos
		avatarURL = avatarPath
		uc.log.Info().Uint("user_id", id).Str("avatar_path", avatarPath).Msg("Nueva imagen de avatar subida exitosamente")

		// Eliminar imagen anterior si existe y es diferente
		if existingUser.AvatarURL != "" && existingUser.AvatarURL != avatarPath {
			// Verificar si la imagen anterior es un path relativo (no URL completa)
			if !strings.HasPrefix(existingUser.AvatarURL, "http") {
				uc.log.Info().Uint("user_id", id).Str("old_avatar", existingUser.AvatarURL).Msg("Eliminando imagen anterior de avatar")
				if err := uc.s3.DeleteImage(ctx, existingUser.AvatarURL); err != nil {
					uc.log.Warn().Err(err).Str("old_avatar", existingUser.AvatarURL).Msg("Error al eliminar imagen anterior (no crítico)")
					// No fallar la actualización si no se puede eliminar la imagen anterior
				}
			}
		}
	} else if userDTO.RemoveAvatar {
		// Eliminar avatar solo si el cliente lo solicita explícitamente
		uc.log.Info().Uint("user_id", id).Str("old_avatar", existingUser.AvatarURL).Msg("Eliminando imagen de avatar")

		// Verificar si la imagen anterior es un path relativo (no URL completa)
		if !strings.HasPrefix(existingUser.AvatarURL, "http") {
			if err := uc.s3.DeleteImage(ctx, existingUser.AvatarURL); err != nil {
				uc.log.Warn().Err(err).Str("old_avatar", existingUser.AvatarURL).Msg("Error al eliminar imagen anterior (no crítico)")
				// No fallar la actualización si no se puede eliminar la imagen
			}
		}
		avatarURL = "" // Limpiar la URL
	}

	// Convertir DTO a entidad
	user := domain.UsersEntity{
		Name:      userDTO.Name,
		Email:     normalizedEmail, // Usar email normalizado
		Phone:     userDTO.Phone,
		AvatarURL: avatarURL, // URL relativa o vacía según corresponda
		IsActive:  userDTO.IsActive,
	}

	// No se permite actualizar contraseña por este endpoint

	// Actualizar usuario
	message, err := uc.repository.UpdateUser(ctx, id, user)
	if err != nil {
		uc.log.Error().Uint("id", id).Err(err).Msg("Error al actualizar usuario desde el repositorio")
		// Verificar si es un error de duplicado de email
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "email") {
			return "", domain.ErrUserEmailExists
		}
		return "", fmt.Errorf("error al actualizar usuario: %w", err)
	}

	// Actualizar relación de businesses
	if len(userDTO.BusinessIDs) > 0 {
		uc.log.Info().Uint("user_id", id).Any("business_ids", userDTO.BusinessIDs).Msg("Actualizando businesses del usuario")
		if err := uc.repository.AssignBusinessesToUser(ctx, id, userDTO.BusinessIDs); err != nil {
			uc.log.Error().Err(err).Uint("user_id", id).Any("business_ids", userDTO.BusinessIDs).Msg("Error al asignar businesses al usuario")
			if strings.Contains(err.Error(), "algunos businesses no existen") {
				return "", domain.ErrBusinessesNotFound
			}
			return "", fmt.Errorf("error al asignar businesses: %w", err)
		}
		uc.log.Info().Uint("user_id", id).Int("businesses_count", len(userDTO.BusinessIDs)).Msg("Businesses actualizados exitosamente")
	}

	uc.log.Info().Uint("user_id", id).Msg("Usuario actualizado exitosamente")
	return message, nil
}
