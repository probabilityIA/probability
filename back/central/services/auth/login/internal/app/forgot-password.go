package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
)

const passwordResetTokenTTL = time.Hour
const otpTTL = 10 * time.Minute
const maxOTPAttempts = 5

func (uc *AuthUseCase) ForgotPassword(ctx context.Context, request domain.ForgotPasswordRequest) (*domain.ForgotPasswordResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))

	genericResponse := &domain.ForgotPasswordResponse{
		Success: true,
		Message: "Si el correo esta registrado, recibiras un enlace para restablecer tu contrasena.",
	}

	if email == "" {
		return genericResponse, nil
	}

	user, err := uc.repository.GetUserByEmail(ctx, email)
	if err != nil {
		uc.log.Error().Err(err).Str("email", email).Msg("Error buscando usuario para recuperacion de contrasena")
		return genericResponse, nil
	}
	if user == nil || !user.IsActive {
		uc.log.Info().Str("email", email).Msg("Solicitud de recuperacion para email inexistente o inactivo")
		return genericResponse, nil
	}

	channel := strings.ToLower(strings.TrimSpace(request.Channel))
	if channel == "" {
		channel = "email"
	}

	if err := uc.repository.InvalidateUserPasswordResetTokens(ctx, user.ID); err != nil {
		uc.log.Warn().Err(err).Uint("user_id", user.ID).Msg("Error invalidando tokens previos")
	}

	if channel == "whatsapp" {
		return uc.forgotPasswordWhatsApp(ctx, user, genericResponse)
	}
	return uc.forgotPasswordEmail(ctx, user, genericResponse)
}

func (uc *AuthUseCase) forgotPasswordEmail(ctx context.Context, user *domain.UserAuthInfo, genericResponse *domain.ForgotPasswordResponse) (*domain.ForgotPasswordResponse, error) {
	rawToken, tokenHash, err := generateResetToken()
	if err != nil {
		uc.log.Error().Err(err).Msg("Error generando token de recuperacion")
		return genericResponse, nil
	}

	expiresAt := time.Now().Add(passwordResetTokenTTL)
	if err := uc.repository.CreatePasswordResetToken(ctx, user.ID, tokenHash, "email", expiresAt); err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.ID).Msg("Error guardando token de recuperacion")
		return genericResponse, nil
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", uc.frontendBaseURL(), rawToken)
	html := buildResetPasswordEmail(user.Name, resetURL)

	if err := uc.emailSender.SendHTML(ctx, user.Email, "Restablece tu contrasena", html); err != nil {
		uc.log.Error().Err(err).Str("email", user.Email).Msg("Error enviando correo de recuperacion")
		return genericResponse, nil
	}

	uc.log.Info().Uint("user_id", user.ID).Str("email", user.Email).Msg("Correo de recuperacion enviado")
	return genericResponse, nil
}

func (uc *AuthUseCase) forgotPasswordWhatsApp(ctx context.Context, user *domain.UserAuthInfo, genericResponse *domain.ForgotPasswordResponse) (*domain.ForgotPasswordResponse, error) {
	phone := strings.TrimSpace(user.Phone)
	if phone == "" {
		uc.log.Info().Uint("user_id", user.ID).Msg("Solicitud OTP WhatsApp para usuario sin telefono")
		return genericResponse, nil
	}

	code, err := generateOTPCode()
	if err != nil {
		uc.log.Error().Err(err).Msg("Error generando codigo OTP")
		return genericResponse, nil
	}

	tokenHash := hashOTPCode(user.ID, code)
	expiresAt := time.Now().Add(otpTTL)
	if err := uc.repository.CreatePasswordResetToken(ctx, user.ID, tokenHash, "whatsapp", expiresAt); err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.ID).Msg("Error guardando codigo OTP")
		return genericResponse, nil
	}

	event := domain.PasswordResetOTPEvent{
		Phone:          phone,
		Code:           code,
		UserName:       user.Name,
		ExpiresMinutes: int(otpTTL.Minutes()),
	}
	if err := uc.otpPublisher.PublishPasswordResetOTP(ctx, event); err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.ID).Msg("Error publicando evento OTP WhatsApp")
		return genericResponse, nil
	}

	uc.log.Info().Uint("user_id", user.ID).Msg("Codigo OTP de recuperacion publicado a WhatsApp")
	return genericResponse, nil
}

func (uc *AuthUseCase) frontendBaseURL() string {
	base := strings.TrimRight(uc.env.Get("FRONTEND_BASE_URL"), "/")
	if base == "" {
		base = "http://localhost:3000"
	}
	return base
}

func generateResetToken() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	rawToken := hex.EncodeToString(bytes)
	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])
	return rawToken, tokenHash, nil
}

func hashResetToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func generateOTPCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func hashOTPCode(userID uint, code string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%d:%s", userID, code)))
	return hex.EncodeToString(sum[:])
}

func maskPhone(phone string) string {
	digits := strings.TrimSpace(phone)
	if len(digits) <= 3 {
		return "***"
	}
	return "***" + digits[len(digits)-3:]
}

func buildResetPasswordEmail(name, resetURL string) string {
	displayName := strings.TrimSpace(name)
	if displayName == "" {
		displayName = "Hola"
	} else {
		displayName = "Hola " + displayName
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<body style="margin:0;padding:0;background-color:#f4f4f7;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f7;padding:24px 0;">
    <tr>
      <td align="center">
        <table width="480" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:12px;overflow:hidden;">
          <tr>
            <td style="padding:32px 32px 8px 32px;">
              <h1 style="margin:0;font-size:20px;color:#111827;">%s,</h1>
              <p style="margin:16px 0 0 0;font-size:15px;line-height:22px;color:#374151;">
                Recibimos una solicitud para restablecer la contrasena de tu cuenta. Haz clic en el boton para crear una nueva contrasena.
              </p>
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:24px 32px;">
              <a href="%s" style="display:inline-block;background-color:#4f46e5;color:#ffffff;text-decoration:none;padding:12px 28px;border-radius:8px;font-size:15px;font-weight:bold;">Restablecer contrasena</a>
            </td>
          </tr>
          <tr>
            <td style="padding:0 32px 32px 32px;">
              <p style="margin:0;font-size:13px;line-height:20px;color:#6b7280;">
                Este enlace caduca en 1 hora. Si no solicitaste este cambio, puedes ignorar este correo de forma segura.
              </p>
              <p style="margin:16px 0 0 0;font-size:12px;line-height:18px;color:#9ca3af;word-break:break-all;">
                Si el boton no funciona, copia y pega este enlace: %s
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, displayName, resetURL, resetURL)
}
