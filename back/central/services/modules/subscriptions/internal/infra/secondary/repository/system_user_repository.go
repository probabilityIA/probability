package repository

import (
	"context"
	"errors"
	"fmt"
)

func (r *Repository) ResolveUserLabel(ctx context.Context, userID uint) (string, error) {
	var label string
	err := r.db.Conn(ctx).Raw(`
		SELECT COALESCE(NULLIF(u.name, ''), u.email, '')
		FROM "user" u WHERE u.id = ? AND u.deleted_at IS NULL
	`, userID).Scan(&label).Error
	if err != nil {
		return "", err
	}
	if label == "" {
		return fmt.Sprintf("Usuario #%d", userID), nil
	}
	return label, nil
}

func (r *Repository) FindSuperAdminUserID(ctx context.Context) (uint, error) {
	var userID uint
	err := r.db.Conn(ctx).Raw(`
		SELECT u.id FROM "user" u
		JOIN scope s ON s.id = u.scope_id
		WHERE s.code = 'platform' AND u.deleted_at IS NULL
		ORDER BY u.id
		LIMIT 1
	`).Scan(&userID).Error
	if err != nil {
		return 0, err
	}
	if userID == 0 {
		return 0, errors.New("no super admin user found")
	}
	return userID, nil
}
