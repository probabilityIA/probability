package app

import (
	"context"
	"fmt"
)

// DeleteUser elimina un usuario
func (uc *UserUseCase) DeleteUser(ctx context.Context, id uint) (string, error) {
	uc.log.Info().Uint("id", id).Msg("Iniciando caso de uso: eliminar usuario")

	// Verificar que el usuario existe
	existingUser, err := uc.repository.GetUserByID(ctx, id)
	if err != nil || existingUser == nil {
		uc.log.Error().Uint("id", id).Msg("Usuario no encontrado")
		return "", fmt.Errorf("usuario no encontrado")
	}

	// Eliminar usuario
	message, err := uc.repository.DeleteUser(ctx, id)
	if err != nil {
		uc.log.Error().Uint("id", id).Err(err).Msg("Error al eliminar usuario desde el repositorio")
		return "", err
	}

	uc.log.Info().Uint("user_id", id).Msg("Usuario eliminado exitosamente")
	return message, nil
}
