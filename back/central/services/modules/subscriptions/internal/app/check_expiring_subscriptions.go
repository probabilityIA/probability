package app

import (
	"context"
	"time"
)

// El aviso de vencimiento de suscripcion (ticker/modal) se elimino a pedido
// del usuario (2026-08-28): era invasivo y ademas se mostraba sin segmentar
// (bug de scoping en el front, ver AnnouncementTicker/AnnouncementModal) a
// cualquier super admin. CheckExpiringSubscriptions ya no crea anuncios; solo
// downgrade de trials, settle de ciclos free y marcar como expirado.
func (uc *UseCase) CheckExpiringSubscriptions(ctx context.Context) error {
	now := time.Now()

	if err := uc.DowngradeExpiredTrials(ctx); err != nil {
		uc.log.Error(ctx).Err(err).Msg("failed to downgrade expired trials")
	}
	if err := uc.SettleFreePlanCycles(ctx); err != nil {
		uc.log.Error(ctx).Err(err).Msg("failed to settle free plan cycles")
	}

	justExpired, err := uc.repo.ListBusinessesJustExpired(ctx, now)
	if err != nil {
		return err
	}
	for _, business := range justExpired {
		if uc.autoRenewIfEnabled(ctx, business) {
			continue
		}
		if !cutoffReached(business, now) {
			continue
		}
		if err := uc.repo.MarkExpiredIfStillActive(ctx, business.BusinessID, now); err != nil {
			uc.log.Error(ctx).Err(err).Uint("business_id", business.BusinessID).Msg("failed to mark business as expired")
		}
	}

	return nil
}

// Planes pagos: el negocio debe pagar para renovar, el mensaje lo dice
// explicitamente. Planes free/trial: no son pagables de la misma forma (el
// trial degrada solo a free, el free se renueva solo salvo excedente sin
// pagar, que ya se avisa aparte en el guard de creacion de guia), asi que el
// mensaje evita decir "realiza el pago".
const (
	paidExpiringSoonTitle   = "Tu suscripcion esta por vencer"
	paidExpiringSoonMessage = "Tu suscripcion vence en menos de 7 dias. Realiza el pago para evitar interrupciones en tus modulos contratados."
	paidExpiredTitle        = "Tu suscripcion vencio"
	paidExpiredMessage      = "Tu suscripcion ya vencio. Realiza el pago para seguir disfrutando de todos los modulos contratados."

	freeTrialExpiringSoonTitle   = "Tu plan actual esta por finalizar"
	freeTrialExpiringSoonMessage = "Tu plan actual finaliza en menos de 7 dias. No necesitas hacer nada: tu cuenta continuara activa en el plan gratuito de Probability."
	freeTrialExpiredTitle        = "Tu plan actual finalizo"
	freeTrialExpiredMessage      = "Tu plan actual ya finalizo. Tu cuenta continua activa en el plan gratuito de Probability."
)

func (uc *UseCase) deactivateExpiryAnnouncements(ctx context.Context, businessID uint) {
	titles := []string{
		paidExpiringSoonTitle, paidExpiredTitle,
		freeTrialExpiringSoonTitle, freeTrialExpiredTitle,
	}
	deactivated := make(map[uint]bool)
	for _, title := range titles {
		id, err := uc.announcements.FindActiveBusinessAlert(ctx, businessID, title)
		if err != nil || id == nil || deactivated[*id] {
			continue
		}
		if err := uc.announcements.DeactivateAnnouncement(ctx, *id); err != nil {
			uc.log.Warn(ctx).Err(err).Uint("business_id", businessID).Msg("failed to deactivate expiry announcement")
			continue
		}
		deactivated[*id] = true
	}
}

func (uc *UseCase) resolveSystemUserID(ctx context.Context) (uint, error) {
	if uc.systemUserID > 0 {
		return uc.systemUserID, nil
	}
	id, err := uc.repo.FindSuperAdminUserID(ctx)
	if err != nil {
		return 0, err
	}
	uc.systemUserID = id
	return id, nil
}
