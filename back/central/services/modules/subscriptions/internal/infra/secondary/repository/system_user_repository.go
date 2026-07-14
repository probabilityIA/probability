package repository

import (
	"context"
	"errors"
)

func (r *Repository) FindSuperAdminUserID(ctx context.Context) (uint, error) {
	var userID uint
	err := r.db.Conn(ctx).Raw(`
		SELECT u.id FROM "user" u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN role r ON r.id = ur.role_id
		JOIN scope s ON s.id = r.scope_id
		WHERE s.code = 'platform' AND u.deleted_at IS NULL
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
