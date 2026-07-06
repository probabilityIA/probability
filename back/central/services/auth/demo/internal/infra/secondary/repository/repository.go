package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	database db.IDatabase
	logger   log.ILogger
	encKey   []byte
}

func New(database db.IDatabase, logger log.ILogger, encryptionKey string) domain.IDemoRepository {
	return &Repository{database: database, logger: logger, encKey: parseEncryptionKey(encryptionKey)}
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.database.Conn(ctx).Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) BusinessCodeExists(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.database.Conn(ctx).Model(&models.Business{}).
		Where("code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) GetBusinessIDByUserID(ctx context.Context, userID uint) (uint, error) {
	var businessID uint
	err := r.database.Conn(ctx).Table("business_staff").Select("business_id").
		Where("user_id = ? AND business_id IS NOT NULL AND deleted_at IS NULL", userID).
		Order("id ASC").Limit(1).Scan(&businessID).Error
	if err != nil {
		return 0, err
	}
	return businessID, nil
}

func (r *Repository) GetDemoRoleID(ctx context.Context) (uint, error) {
	var id uint
	err := r.database.Conn(ctx).Table("role").Select("id").
		Where("name = ? AND deleted_at IS NULL", "demo").Limit(1).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) CreateDemoAccount(ctx context.Context, p domain.CreateDemoAccountParams) (uint, error) {
	var userID uint
	err := r.database.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		business := &models.Business{
			Name:               p.BusinessName,
			Code:               p.BusinessCode,
			BusinessTypeID:     1,
			OrderPrefix:        p.OrderPrefix,
			IsActive:           true,
			SubscriptionStatus: "active",
		}
		if err := tx.Create(business).Error; err != nil {
			return err
		}

		scopeBusiness := uint(2)
		user := &models.User{
			Name:     p.FullName,
			Email:    p.Email,
			Phone:    p.Phone,
			Password: p.PasswordHash,
			IsActive: false,
			ScopeID:  &scopeBusiness,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Update("is_active", false).Error; err != nil {
			return err
		}
		userID = user.ID

		if err := tx.Table("user_businesses").Create(map[string]any{
			"user_id":     user.ID,
			"business_id": business.ID,
		}).Error; err != nil {
			return err
		}

		staff := &models.BusinessStaff{
			UserID:     user.ID,
			BusinessID: &business.ID,
			RoleID:     &p.RoleID,
		}
		if err := tx.Create(staff).Error; err != nil {
			return err
		}

		token := &models.EmailVerificationToken{
			UserID:    user.ID,
			TokenHash: p.TokenHash,
			ExpiresAt: p.ExpiresAt,
		}
		if err := tx.Create(token).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		r.logger.Error().Err(err).Str("email", p.Email).Msg("Error creando cuenta demo")
		return 0, err
	}
	return userID, nil
}

func (r *Repository) GetValidEmailVerificationToken(ctx context.Context, tokenHash string) (*domain.EmailVerificationTokenInfo, error) {
	var token models.EmailVerificationToken
	err := r.database.Conn(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &domain.EmailVerificationTokenInfo{
		ID:        token.ID,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
		UsedAt:    token.UsedAt,
	}, nil
}

func (r *Repository) ActivateUserAndConsumeToken(ctx context.Context, tokenID, userID uint) error {
	now := time.Now()
	return r.database.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", userID).
			Update("is_active", true).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.EmailVerificationToken{}).Where("id = ?", tokenID).
			Update("used_at", now).Error; err != nil {
			return err
		}
		return nil
	})
}
