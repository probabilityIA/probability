package app

import (
	"context"
	"crypto/subtle"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
)

func (uc *AuthUseCase) VerifyOTP(ctx context.Context, request domain.VerifyOTPRequest) (*domain.VerifyOTPResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	code := strings.TrimSpace(request.Code)

	invalidResponse := &domain.VerifyOTPResponse{
		Success: false,
		Message: "Codigo invalido o expirado",
	}

	if email == "" || code == "" {
		return invalidResponse, nil
	}

	user, err := uc.repository.GetUserByEmail(ctx, email)
	if err != nil {
		uc.log.Error().Err(err).Str("email", email).Msg("Error buscando usuario en verificacion OTP")
		return invalidResponse, nil
	}
	if user == nil || !user.IsActive {
		return invalidResponse, nil
	}

	tokenInfo, err := uc.repository.GetActiveOTPToken(ctx, user.ID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.ID).Msg("Error obteniendo codigo OTP activo")
		return invalidResponse, nil
	}
	if tokenInfo == nil {
		return invalidResponse, nil
	}
	if time.Now().After(tokenInfo.ExpiresAt) {
		return invalidResponse, nil
	}
	if tokenInfo.Attempts >= maxOTPAttempts {
		if err := uc.repository.MarkPasswordResetTokenUsed(ctx, tokenInfo.ID); err != nil {
			uc.log.Warn().Err(err).Uint("token_id", tokenInfo.ID).Msg("Error invalidando OTP por exceso de intentos")
		}
		return invalidResponse, nil
	}

	expectedHash := hashOTPCode(user.ID, code)
	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(tokenInfo.TokenHash)) != 1 {
		if err := uc.repository.IncrementPasswordResetTokenAttempts(ctx, tokenInfo.ID); err != nil {
			uc.log.Warn().Err(err).Uint("token_id", tokenInfo.ID).Msg("Error incrementando intentos OTP")
		}
		return invalidResponse, nil
	}

	if err := uc.repository.MarkPasswordResetTokenUsed(ctx, tokenInfo.ID); err != nil {
		uc.log.Warn().Err(err).Uint("token_id", tokenInfo.ID).Msg("Error marcando OTP como usado")
	}

	rawToken, tokenHash, err := generateResetToken()
	if err != nil {
		uc.log.Error().Err(err).Msg("Error generando token tras verificar OTP")
		return &domain.VerifyOTPResponse{Success: false, Message: "Error interno del servidor"}, nil
	}

	expiresAt := time.Now().Add(15 * time.Minute)
	if err := uc.repository.CreatePasswordResetToken(ctx, user.ID, tokenHash, "email", expiresAt); err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.ID).Msg("Error guardando token tras verificar OTP")
		return &domain.VerifyOTPResponse{Success: false, Message: "Error interno del servidor"}, nil
	}

	uc.log.Info().Uint("user_id", user.ID).Msg("Codigo OTP verificado exitosamente")
	return &domain.VerifyOTPResponse{
		Success: true,
		Message: "Codigo verificado",
		Token:   rawToken,
	}, nil
}
