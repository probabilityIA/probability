package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
)

func (uc *AuthUseCase) ResetPassword(ctx context.Context, request domain.ResetPasswordRequest) (*domain.ResetPasswordResponse, error) {
	rawToken := strings.TrimSpace(request.Token)
	if rawToken == "" {
		return nil, fmt.Errorf("token invalido")
	}
	if len(request.NewPassword) < 6 {
		return nil, fmt.Errorf("la contrasena debe tener al menos 6 caracteres")
	}

	tokenHash := hashResetToken(rawToken)

	tokenInfo, err := uc.repository.GetValidPasswordResetToken(ctx, tokenHash)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error consultando token de recuperacion")
		return nil, fmt.Errorf("error interno del servidor")
	}
	if tokenInfo == nil {
		return nil, fmt.Errorf("token invalido o expirado")
	}
	if tokenInfo.UsedAt != nil {
		return nil, fmt.Errorf("token invalido o expirado")
	}
	if time.Now().After(tokenInfo.ExpiresAt) {
		return nil, fmt.Errorf("token invalido o expirado")
	}

	if err := uc.repository.ChangePassword(ctx, tokenInfo.UserID, request.NewPassword); err != nil {
		uc.log.Error().Err(err).Uint("user_id", tokenInfo.UserID).Msg("Error actualizando contrasena")
		return nil, fmt.Errorf("error al actualizar la contrasena")
	}

	if err := uc.repository.MarkPasswordResetTokenUsed(ctx, tokenInfo.ID); err != nil {
		uc.log.Warn().Err(err).Uint("token_id", tokenInfo.ID).Msg("Error marcando token como usado")
	}

	uc.log.Info().Uint("user_id", tokenInfo.UserID).Msg("Contrasena restablecida exitosamente")

	return &domain.ResetPasswordResponse{
		Success: true,
		Message: "Contrasena restablecida exitosamente",
	}, nil
}
