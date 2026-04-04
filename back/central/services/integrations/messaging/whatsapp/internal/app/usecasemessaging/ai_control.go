package usecasemessaging

import (
	"context"
	"fmt"
)

// PauseAI pausa el bot AI para un número de teléfono y activa la sesión humana.
// Llamado cuando el humano decide tomar control del chat desde el dashboard.
func (u *usecases) PauseAI(ctx context.Context, conversationID, phoneNumber string, businessID uint) error {
	phoneNumber = NormalizePhoneNumber(phoneNumber)

	// 1. Marcar AI como pausado en Redis
	if err := u.conversationCache.SetAIPaused(ctx, phoneNumber, conversationID, businessID); err != nil {
		return fmt.Errorf("error pausando AI: %w", err)
	}

	// 2. Activar sesión humana para que las respuestas lleguen al dashboard
	if err := u.conversationCache.ActivateHumanSession(ctx, phoneNumber, conversationID, businessID); err != nil {
		u.log.Error(ctx).Err(err).Msg("[PauseAI] - error activando human session")
		// No retornamos: la pausa ya está aplicada
	}

	u.log.Info(ctx).
		Str("conversation_id", conversationID).
		Str("phone_number", phoneNumber).
		Uint("business_id", businessID).
		Msg("[WhatsApp UseCase] - AI pausado, humano toma control")

	return nil
}

// ResumeAI reactiva el bot AI para un número de teléfono.
func (u *usecases) ResumeAI(ctx context.Context, conversationID, phoneNumber string, businessID uint) error {
	phoneNumber = NormalizePhoneNumber(phoneNumber)

	if err := u.conversationCache.ClearAIPaused(ctx, phoneNumber); err != nil {
		return fmt.Errorf("error reactivando AI: %w", err)
	}

	u.log.Info(ctx).
		Str("conversation_id", conversationID).
		Str("phone_number", phoneNumber).
		Uint("business_id", businessID).
		Msg("[WhatsApp UseCase] - AI reactivado")

	return nil
}
