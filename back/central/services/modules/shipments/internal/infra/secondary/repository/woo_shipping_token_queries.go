package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func randomWooSalt() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *Repository) GetWooShippingToken(ctx context.Context, integrationID uint) (string, bool, bool, error) {
	var m models.WooShippingToken
	err := r.db.Conn(ctx).Where("integration_id = ?", integrationID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, false, nil
	}
	if err != nil {
		return "", false, false, err
	}
	return m.Salt, m.Revoked, true, nil
}

func (r *Repository) EnsureWooShippingToken(ctx context.Context, integrationID uint) (string, bool, error) {
	var m models.WooShippingToken
	err := r.db.Conn(ctx).Where("integration_id = ?", integrationID).First(&m).Error
	if err == nil {
		return m.Salt, m.Revoked, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, err
	}

	salt, gerr := randomWooSalt()
	if gerr != nil {
		return "", false, gerr
	}
	m = models.WooShippingToken{IntegrationID: integrationID, Salt: salt, Revoked: false}
	if cerr := r.db.Conn(ctx).Create(&m).Error; cerr != nil {
		if rerr := r.db.Conn(ctx).Where("integration_id = ?", integrationID).First(&m).Error; rerr == nil {
			return m.Salt, m.Revoked, nil
		}
		return "", false, cerr
	}
	return m.Salt, m.Revoked, nil
}

func (r *Repository) RotateWooShippingToken(ctx context.Context, integrationID uint) (string, error) {
	salt, err := randomWooSalt()
	if err != nil {
		return "", err
	}

	var m models.WooShippingToken
	e := r.db.Conn(ctx).Where("integration_id = ?", integrationID).First(&m).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		m = models.WooShippingToken{IntegrationID: integrationID, Salt: salt, Revoked: false}
		if cerr := r.db.Conn(ctx).Create(&m).Error; cerr != nil {
			return "", cerr
		}
		return salt, nil
	}
	if e != nil {
		return "", e
	}

	m.Salt = salt
	m.Revoked = false
	if uerr := r.db.Conn(ctx).Save(&m).Error; uerr != nil {
		return "", uerr
	}
	return salt, nil
}

func (r *Repository) RevokeWooShippingToken(ctx context.Context, integrationID uint) error {
	var m models.WooShippingToken
	e := r.db.Conn(ctx).Where("integration_id = ?", integrationID).First(&m).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		salt, gerr := randomWooSalt()
		if gerr != nil {
			return gerr
		}
		m = models.WooShippingToken{IntegrationID: integrationID, Salt: salt, Revoked: true}
		return r.db.Conn(ctx).Create(&m).Error
	}
	if e != nil {
		return e
	}

	m.Revoked = true
	return r.db.Conn(ctx).Save(&m).Error
}
