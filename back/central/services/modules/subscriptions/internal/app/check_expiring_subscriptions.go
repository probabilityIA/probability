package app

import (
	"context"
	"time"
)

const expiryWarnDays = 7

func (uc *UseCase) CheckExpiringSubscriptions(ctx context.Context) error {
	now := time.Now()
	warnUntil := now.AddDate(0, 0, expiryWarnDays)

	expiringSoon, err := uc.repo.ListBusinessesExpiringBetween(ctx, now, warnUntil)
	if err != nil {
		return err
	}
	for _, businessID := range expiringSoon {
		if err := uc.ensureExpiryAnnouncement(ctx, businessID, expiringSoonTitle, expiringSoonMessage, true); err != nil {
			uc.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("failed to ensure expiring soon announcement")
		}
	}

	justExpired, err := uc.repo.ListBusinessesJustExpired(ctx, now)
	if err != nil {
		return err
	}
	for _, businessID := range justExpired {
		if err := uc.ensureExpiryAnnouncement(ctx, businessID, expiredTitle, expiredMessage, true); err != nil {
			uc.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("failed to ensure expired announcement")
		}
		if err := uc.repo.UpdateBusinessSubscriptionStatus(ctx, businessID, "expired", nil); err != nil {
			uc.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("failed to mark business as expired")
		}
	}

	return nil
}

const (
	expiringSoonTitle   = "Tu suscripcion esta por vencer"
	expiringSoonMessage = "Tu suscripcion vence en menos de 7 dias. Realiza el pago para evitar interrupciones en tus modulos contratados."
	expiredTitle        = "Tu suscripcion vencio"
	expiredMessage      = "Tu suscripcion ya vencio. Realiza el pago para seguir disfrutando de todos los modulos contratados."
)

func (uc *UseCase) ensureExpiryAnnouncement(ctx context.Context, businessID uint, title, message string, daily bool) error {
	existing, err := uc.announcements.FindActiveBusinessAlert(ctx, businessID, title)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	systemUserID, err := uc.resolveSystemUserID(ctx)
	if err != nil {
		return err
	}

	_, err = uc.announcements.CreateBusinessAlert(ctx, businessID, title, message, systemUserID, daily)
	return err
}

func (uc *UseCase) deactivateExpiryAnnouncements(ctx context.Context, businessID uint) {
	for _, title := range []string{expiringSoonTitle, expiredTitle} {
		id, err := uc.announcements.FindActiveBusinessAlert(ctx, businessID, title)
		if err != nil || id == nil {
			continue
		}
		if err := uc.announcements.DeactivateAnnouncement(ctx, *id); err != nil {
			uc.log.Warn(ctx).Err(err).Uint("business_id", businessID).Msg("failed to deactivate expiry announcement")
		}
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
