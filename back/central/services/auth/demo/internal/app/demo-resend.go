package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
)

const resendCooldown = 60 * time.Second

func (uc *UseCase) DemoResend(ctx context.Context, request domain.DemoResendRequest) (*domain.DemoResendResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	channel := strings.ToLower(strings.TrimSpace(request.Channel))
	if channel == "" {
		channel = "email"
	}
	phone := strings.TrimSpace(request.Phone)

	if channel == "whatsapp" && phone == "" {
		return nil, fmt.Errorf("el telefono es obligatorio para verificar por WhatsApp")
	}

	generic := &domain.DemoResendResponse{
		Success: true,
		Message: "Si existe una cuenta pendiente de verificacion, te enviamos un nuevo enlace.",
	}
	if email == "" {
		return generic, nil
	}

	user, err := uc.repository.GetDemoUserByEmail(ctx, email)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error consultando usuario para reenvio de verificacion")
		return generic, nil
	}
	if user == nil || user.IsActive {
		return generic, nil
	}

	if user.LastTokenCreatedAt != nil && time.Since(*user.LastTokenCreatedAt) < resendCooldown {
		uc.log.Warn().Uint("user_id", user.UserID).Msg("Reenvio de verificacion demo bloqueado por cooldown")
		return generic, nil
	}

	var rawToken, tokenHash, code string
	ttl := emailVerificationTTL
	if channel == "whatsapp" {
		code, err = generateOTPCode()
		if err != nil {
			return generic, nil
		}
		tokenHash = hashOTP(email, code)
		ttl = otpVerificationTTL
	} else {
		rawToken, tokenHash, err = generateToken()
		if err != nil {
			return generic, nil
		}
	}

	if err := uc.repository.InvalidateEmailVerificationTokens(ctx, user.UserID); err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.UserID).Msg("Error invalidando tokens previos de verificacion demo")
		return generic, nil
	}
	if err := uc.repository.CreateEmailVerificationToken(ctx, user.UserID, tokenHash, time.Now().Add(ttl)); err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.UserID).Msg("Error creando token de verificacion demo")
		return generic, nil
	}

	if channel == "whatsapp" {
		if phone != user.Phone {
			if err := uc.repository.UpdateUserPhone(ctx, user.UserID, phone); err != nil {
				uc.log.Error().Err(err).Uint("user_id", user.UserID).Msg("Error actualizando telefono para reenvio demo")
			}
		}
		event := domain.DemoOTPEvent{
			Phone:          phone,
			Code:           code,
			UserName:       user.FullName,
			ExpiresMinutes: int(ttl.Minutes()),
		}
		if err := uc.otpPublisher.PublishDemoOTP(ctx, event); err != nil {
			uc.log.Error().Err(err).Uint("user_id", user.UserID).Msg("Error publicando codigo OTP demo por WhatsApp")
			return generic, nil
		}
		uc.log.Info().Uint("user_id", user.UserID).Msg("Codigo de verificacion demo reenviado por WhatsApp")
		return generic, nil
	}

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", uc.frontendBaseURL(), rawToken)
	html := buildVerificationEmail(user.FullName, user.BusinessName, verifyURL)
	if err := uc.emailSender.SendHTML(ctx, email, "Verifica tu cuenta demo de Probability", html); err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.UserID).Msg("Error reenviando correo de verificacion demo")
		return generic, nil
	}

	uc.log.Info().Uint("user_id", user.UserID).Msg("Correo de verificacion demo reenviado")
	return generic, nil
}
