package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/push/internal/domain/errors"
)

var allowedPlatforms = map[string]bool{"android": true, "ios": true, "web": true}

func (u *UseCase) RegisterDevice(ctx context.Context, dto dtos.RegisterDeviceDTO) error {
	dto.Token = strings.TrimSpace(dto.Token)
	if dto.Token == "" {
		return domainerrors.ErrTokenRequired
	}

	dto.Platform = strings.ToLower(strings.TrimSpace(dto.Platform))
	if !allowedPlatforms[dto.Platform] {
		return domainerrors.ErrPlatformInvalid
	}

	if err := u.repo.UpsertDeviceToken(ctx, dto); err != nil {
		u.log.Error(ctx).Err(err).Uint("user_id", dto.UserID).Msg("no se pudo registrar el dispositivo")
		return err
	}

	u.log.Info(ctx).
		Uint("user_id", dto.UserID).
		Str("platform", dto.Platform).
		Msg("dispositivo registrado para notificaciones push")
	return nil
}

func (u *UseCase) UnregisterDevice(ctx context.Context, userID uint, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return domainerrors.ErrTokenRequired
	}
	return u.repo.DeactivateToken(ctx, userID, token)
}

func (u *UseCase) ListDevices(ctx context.Context, userID uint) ([]entities.DeviceToken, error) {
	return u.repo.ListDevicesByUser(ctx, userID)
}
