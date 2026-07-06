package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreatePasswordResetToken(ctx context.Context, userID uint, tokenHash string, channel string, expiresAt time.Time) error {
	if channel == "" {
		channel = "email"
	}
	token := &models.PasswordResetToken{
		UserID:    userID,
		TokenHash: tokenHash,
		Channel:   channel,
		ExpiresAt: expiresAt,
	}
	if err := r.database.Conn(ctx).Create(token).Error; err != nil {
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error creando token de recuperacion")
		return err
	}
	return nil
}

func (r *Repository) InvalidateUserPasswordResetTokens(ctx context.Context, userID uint) error {
	now := time.Now()
	if err := r.database.Conn(ctx).
		Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", now).Error; err != nil {
		r.logger.Error().Uint("user_id", userID).Err(err).Msg("Error invalidando tokens de recuperacion")
		return err
	}
	return nil
}

func (r *Repository) GetValidPasswordResetToken(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
	var token models.PasswordResetToken
	if err := r.database.Conn(ctx).
		Where("token_hash = ?", tokenHash).
		First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Msg("Error obteniendo token de recuperacion")
		return nil, err
	}
	return &domain.PasswordResetTokenInfo{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		Channel:   token.Channel,
		Attempts:  token.Attempts,
		ExpiresAt: token.ExpiresAt,
		UsedAt:    token.UsedAt,
	}, nil
}

func (r *Repository) GetActiveOTPToken(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error) {
	var token models.PasswordResetToken
	if err := r.database.Conn(ctx).
		Where("user_id = ? AND channel = ? AND used_at IS NULL", userID, "whatsapp").
		Order("created_at DESC").
		First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Error obteniendo token OTP activo")
		return nil, err
	}
	return &domain.PasswordResetTokenInfo{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		Channel:   token.Channel,
		Attempts:  token.Attempts,
		ExpiresAt: token.ExpiresAt,
		UsedAt:    token.UsedAt,
	}, nil
}

func (r *Repository) IncrementPasswordResetTokenAttempts(ctx context.Context, tokenID uint) error {
	if err := r.database.Conn(ctx).
		Model(&models.PasswordResetToken{}).
		Where("id = ?", tokenID).
		UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error; err != nil {
		r.logger.Error().Err(err).Uint("token_id", tokenID).Msg("Error incrementando intentos de OTP")
		return err
	}
	return nil
}

func (r *Repository) MarkPasswordResetTokenUsed(ctx context.Context, tokenID uint) error {
	now := time.Now()
	if err := r.database.Conn(ctx).
		Model(&models.PasswordResetToken{}).
		Where("id = ?", tokenID).
		Update("used_at", now).Error; err != nil {
		r.logger.Error().Uint("token_id", tokenID).Err(err).Msg("Error marcando token como usado")
		return err
	}
	return nil
}
