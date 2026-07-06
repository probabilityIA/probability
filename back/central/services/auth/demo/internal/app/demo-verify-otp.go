package app

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
)

func (uc *UseCase) DemoVerifyOTP(ctx context.Context, request domain.DemoVerifyOTPRequest) (*domain.DemoVerifyOTPResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	code := strings.TrimSpace(request.Code)

	invalid := &domain.DemoVerifyOTPResponse{Success: false, Message: "Codigo invalido o expirado"}
	if email == "" || code == "" {
		return invalid, nil
	}

	tokenHash := hashOTP(email, code)
	tokenInfo, err := uc.repository.GetValidEmailVerificationToken(ctx, tokenHash)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error consultando codigo de verificacion demo")
		return invalid, nil
	}
	if tokenInfo == nil {
		return invalid, nil
	}

	alreadyVerified := tokenInfo.UsedAt != nil
	if !alreadyVerified && time.Now().After(tokenInfo.ExpiresAt) {
		return invalid, nil
	}

	if !alreadyVerified {
		if err := uc.repository.ActivateUserAndConsumeToken(ctx, tokenInfo.ID, tokenInfo.UserID); err != nil {
			uc.log.Error().Err(err).Uint("user_id", tokenInfo.UserID).Msg("Error activando usuario demo por OTP")
			return &domain.DemoVerifyOTPResponse{Success: false, Message: "No se pudo verificar la cuenta"}, nil
		}
	}

	uc.provisionDemo(ctx, tokenInfo.UserID)

	uc.log.Info().Uint("user_id", tokenInfo.UserID).Bool("already_verified", alreadyVerified).Msg("Cuenta demo verificada por OTP")
	return &domain.DemoVerifyOTPResponse{
		Success: true,
		Message: "Cuenta verificada. Ya puedes iniciar sesion y usar la plataforma.",
	}, nil
}
