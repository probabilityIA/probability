package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
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

// GeneratePassword genera una nueva contraseña aleatoria para un usuario y la reemplaza
func (uc *AuthUseCase) GeneratePassword(ctx context.Context, request domain.GeneratePasswordRequest) (*domain.GeneratePasswordResponse, error) {
	uc.log.Info().Uint("user_id", request.UserID).Msg("Iniciando generación de nueva contraseña aleatoria")

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

	// Generar contraseña aleatoria
	generatedPassword, err := generateRandomPassword(12)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", request.UserID).Msg("Error al generar contraseña aleatoria")
		return nil, fmt.Errorf("error al generar contraseña")
	}

	uc.log.Info().Uint("user_id", request.UserID).Msg("Contraseña aleatoria generada exitosamente")

	// Actualizar contraseña en la base de datos (el repositorio la hashea)
	if err := uc.repository.ChangePassword(ctx, request.UserID, generatedPassword); err != nil {
		uc.log.Error().Err(err).Uint("user_id", request.UserID).Msg("Error al actualizar contraseña")
		return nil, fmt.Errorf("error al actualizar contraseña")
	}

	uc.log.Info().Uint("user_id", request.UserID).Str("email", user.Email).Msg("Contraseña actualizada exitosamente")

	return &domain.GeneratePasswordResponse{
		Success:  true,
		Email:    user.Email,
		Password: generatedPassword,
		Message:  fmt.Sprintf("Nueva contraseña generada para el usuario %s", user.Email),
	}, nil
}
