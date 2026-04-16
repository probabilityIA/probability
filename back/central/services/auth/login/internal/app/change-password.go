package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// ChangePassword maneja el cambio de contraseña del usuario
func (uc *AuthUseCase) ChangePassword(ctx context.Context, request domain.ChangePasswordRequest) (*domain.ChangePasswordResponse, error) {
	uc.log.Info().Uint("user_id", request.UserID).Msg("Iniciando cambio de contraseña")

	// Obtener usuario por ID
	user, err := uc.repository.GetUserByID(ctx, request.UserID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", request.UserID).Msg("Error al obtener usuario")
		return nil, fmt.Errorf("error interno del servidor")
	}

	if user == nil {
		uc.log.Error().Uint("user_id", request.UserID).Msg("Usuario no encontrado")
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Verificar que el usuario esté activo
	if !user.IsActive {
		uc.log.Error().Uint("user_id", request.UserID).Msg("Usuario inactivo")
		return nil, fmt.Errorf("usuario inactivo")
	}

	// Validar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.CurrentPassword)); err != nil {
		uc.log.Error().Err(err).Uint("user_id", request.UserID).Msg("Contraseña actual inválida")
		return nil, fmt.Errorf("contraseña actual incorrecta")
	}

	// Validar que la nueva contraseña sea diferente a la actual
	if request.CurrentPassword == request.NewPassword {
		uc.log.Error().Uint("user_id", request.UserID).Msg("La nueva contraseña debe ser diferente a la actual")
		return nil, fmt.Errorf("la nueva contraseña debe ser diferente a la actual")
	}

	// Cambiar contraseña
	if err := uc.repository.ChangePassword(ctx, request.UserID, request.NewPassword); err != nil {
		uc.log.Error().Err(err).Uint("user_id", request.UserID).Msg("Error al cambiar contraseña")
		return nil, fmt.Errorf("error al cambiar contraseña")
	}

	// Actualizar último login después del cambio de contraseña exitoso
	if err := uc.repository.UpdateLastLogin(ctx, request.UserID); err != nil {
		uc.log.Warn().Err(err).Uint("user_id", request.UserID).Msg("Error al actualizar último login después del cambio de contraseña")
		// No retornamos error aquí porque el cambio de contraseña ya fue exitoso
	} else {
		uc.log.Info().Uint("user_id", request.UserID).Msg("Último login actualizado después del cambio de contraseña")
	}

	uc.log.Info().Uint("user_id", request.UserID).Msg("Contraseña cambiada exitosamente")

	return &domain.ChangePasswordResponse{
		Success: true,
		Message: "Contraseña cambiada exitosamente",
	}, nil
}
