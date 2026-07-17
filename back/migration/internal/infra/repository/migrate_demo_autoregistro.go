package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func nowUTC() time.Time {
	return time.Now().UTC()
}

const demoRoleName = "demo"

func (r *Repository) migrateDemoAutoregistro(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.EmailVerificationToken{}); err != nil {
		return fmt.Errorf("failed to auto-migrate email_verification_tokens: %w", err)
	}

	var businessScopeID uint
	if err := db.Table("scope").Select("id").Where("code = ?", "business").Limit(1).Scan(&businessScopeID).Error; err != nil {
		return fmt.Errorf("failed to read business scope: %w", err)
	}
	if businessScopeID == 0 {
		return fmt.Errorf("business scope not found")
	}

	var businessTypeID uint
	if err := db.Table("business_type").Select("id").Order("id ASC").Limit(1).Scan(&businessTypeID).Error; err != nil {
		return fmt.Errorf("failed to read business type: %w", err)
	}

	var roleID uint
	if err := db.Table("role").Select("id").Where("name = ? AND deleted_at IS NULL", demoRoleName).Limit(1).Scan(&roleID).Error; err != nil {
		return fmt.Errorf("failed to read demo role: %w", err)
	}
	if roleID == 0 {
		insert := map[string]any{
			"name":             demoRoleName,
			"description":      "Rol de demostracion: acceso limitado a ordenes, facturacion, guias e inventario",
			"level":            4,
			"is_system":        false,
			"scope_id":         businessScopeID,
			"business_type_id": businessTypeID,
			"created_at":       nowUTC(),
			"updated_at":       nowUTC(),
		}
		if err := db.Table("role").Create(insert).Error; err != nil {
			return fmt.Errorf("failed to create demo role: %w", err)
		}
		if err := db.Table("role").Select("id").Where("name = ?", demoRoleName).Limit(1).Scan(&roleID).Error; err != nil {
			return fmt.Errorf("failed to re-read demo role: %w", err)
		}
	}

	createOrdenesID, err := r.ensurePermission(ctx, "Crear Ordenes", "Permiso para crear ordenes", 6, 1, businessScopeID)
	if err != nil {
		return err
	}

	readBilleteraID, err := r.ensurePermission(ctx, "Read Billetera", "Permiso para ver la billetera", 21, 2, businessScopeID)
	if err != nil {
		return err
	}

	createBilleteraID, err := r.ensurePermission(ctx, "Create Billetera", "Permiso para recargar saldo en la billetera", 21, 1, businessScopeID)
	if err != nil {
		return err
	}

	permissionIDs := []uint{1, 2, createOrdenesID, 5, 66, 67, 68, 69, 6, 20, readBilleteraID, createBilleteraID}
	for _, pid := range permissionIDs {
		if pid == 0 {
			continue
		}
		var exists int64
		if err := db.Table("role_permissions").Where("role_id = ? AND permission_id = ?", roleID, pid).Count(&exists).Error; err != nil {
			return fmt.Errorf("failed to check role_permissions: %w", err)
		}
		if exists > 0 {
			continue
		}
		if err := db.Table("role_permissions").Create(map[string]any{"role_id": roleID, "permission_id": pid}).Error; err != nil {
			return fmt.Errorf("failed to link permission %d to demo role: %w", pid, err)
		}
	}

	return nil
}

func (r *Repository) ensurePermission(ctx context.Context, name, description string, resourceID, actionID, scopeID uint) (uint, error) {
	db := r.db.Conn(ctx)
	var id uint
	if err := db.Table("permission").Select("id").
		Where("resource_id = ? AND action_id = ? AND scope_id = ? AND deleted_at IS NULL", resourceID, actionID, scopeID).
		Limit(1).Scan(&id).Error; err != nil {
		return 0, fmt.Errorf("failed to read permission: %w", err)
	}
	if id != 0 {
		return id, nil
	}
	insert := map[string]any{
		"name":        name,
		"description": description,
		"resource_id": resourceID,
		"action_id":   actionID,
		"scope_id":    scopeID,
		"created_at":  nowUTC(),
		"updated_at":  nowUTC(),
	}
	if err := db.Table("permission").Create(insert).Error; err != nil {
		return 0, fmt.Errorf("failed to create permission %s: %w", name, err)
	}
	if err := db.Table("permission").Select("id").
		Where("resource_id = ? AND action_id = ? AND scope_id = ?", resourceID, actionID, scopeID).
		Limit(1).Scan(&id).Error; err != nil {
		return 0, fmt.Errorf("failed to re-read permission: %w", err)
	}
	return id, nil
}
