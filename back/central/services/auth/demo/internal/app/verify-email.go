package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
)

func (uc *UseCase) VerifyEmail(ctx context.Context, request domain.VerifyEmailRequest) (*domain.VerifyEmailResponse, error) {
	raw := strings.TrimSpace(request.Token)
	if raw == "" {
		return nil, fmt.Errorf("token invalido")
	}

	tokenHash := hashToken(raw)
	tokenInfo, err := uc.repository.GetValidEmailVerificationToken(ctx, tokenHash)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error consultando token de verificacion")
		return nil, fmt.Errorf("error interno del servidor")
	}
	if tokenInfo == nil {
		return nil, fmt.Errorf("token invalido o expirado")
	}
	alreadyVerified := tokenInfo.UsedAt != nil
	if !alreadyVerified && time.Now().After(tokenInfo.ExpiresAt) {
		return nil, fmt.Errorf("token invalido o expirado")
	}

	if !alreadyVerified {
		if err := uc.repository.ActivateUserAndConsumeToken(ctx, tokenInfo.ID, tokenInfo.UserID); err != nil {
			uc.log.Error().Err(err).Uint("user_id", tokenInfo.UserID).Msg("Error activando usuario demo")
			return nil, fmt.Errorf("no se pudo verificar la cuenta")
		}
	}

	uc.provisionDemo(ctx, tokenInfo.UserID)

	uc.log.Info().Uint("user_id", tokenInfo.UserID).Bool("already_verified", alreadyVerified).Msg("Cuenta demo verificada")
	if alreadyVerified {
		return &domain.VerifyEmailResponse{Success: true, Message: "Tu cuenta ya estaba verificada. Puedes iniciar sesion."}, nil
	}
	return &domain.VerifyEmailResponse{
		Success: true,
		Message: "Cuenta verificada. Ya puedes iniciar sesion y usar la plataforma.",
	}, nil
}

func (uc *UseCase) provisionDemo(ctx context.Context, userID uint) {
	businessID, err := uc.repository.GetBusinessIDByUserID(ctx, userID)
	if err != nil || businessID == 0 {
		uc.log.Warn().Err(err).Uint("user_id", userID).Msg("No se pudo resolver el negocio para aprovisionar integraciones demo")
		return
	}
	if err := uc.repository.ProvisionDemoIntegrations(ctx, businessID, userID); err != nil {
		uc.log.Error().Err(err).Uint("business_id", businessID).Msg("Error aprovisionando integraciones demo")
	}
}

func buildVerificationEmail(name, businessName, verifyURL string) string {
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
    <tr><td align="center">
      <table width="480" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:12px;overflow:hidden;">
        <tr><td style="padding:32px 32px 8px 32px;">
          <h1 style="margin:0;font-size:20px;color:#111827;">%s,</h1>
          <p style="margin:16px 0 0 0;font-size:15px;line-height:22px;color:#374151;">
            Tu demo de Probability para <b>%s</b> esta casi lista. Verifica tu correo para activar la cuenta y empezar a probar la plataforma.
          </p>
        </td></tr>
        <tr><td align="center" style="padding:24px 32px;">
          <a href="%s" style="display:inline-block;background-color:#4f46e5;color:#ffffff;text-decoration:none;padding:12px 28px;border-radius:8px;font-size:15px;font-weight:bold;">Verificar mi cuenta</a>
        </td></tr>
        <tr><td style="padding:0 32px 32px 32px;">
          <p style="margin:0;font-size:13px;line-height:20px;color:#6b7280;">
            Este enlace caduca en 24 horas. Es un entorno de demostracion: las integraciones de facturacion y envios son simuladas.
          </p>
          <p style="margin:16px 0 0 0;font-size:12px;line-height:18px;color:#9ca3af;word-break:break-all;">
            Si el boton no funciona, copia y pega: %s
          </p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, displayName, businessName, verifyURL, verifyURL)
}
