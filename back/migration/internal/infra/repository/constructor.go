package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	db  db.IDatabase
	cfg env.IConfig
}

func New(db db.IDatabase, cfg env.IConfig) *Repository {
	return &Repository{
		db:  db,
		cfg: cfg,
	}
}

func (r *Repository) Migrate(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.BusinessType{},
		&models.Scope{},
		&models.Business{},
		&models.BusinessNotificationConfig{},
		&models.BusinessResourceConfigured{},
		&models.Resource{},
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.BusinessStaff{},
		&models.Client{},
		&models.Action{},
		&models.APIKey{},
		&models.IntegrationType{},
		&models.Integration{},

		// Integration Notification Configs (debe ir después de Integration)
		&models.IntegrationNotificationConfig{},

		// Payment Methods
		&models.PaymentMethod{},
		&models.PaymentMethodMapping{},
		&models.OrderStatusMapping{},
		&models.Product{},

		// Orders
		&models.Order{},
		&models.OrderHistory{},
		&models.OrderError{},

		// Order Channel Metadata
		&models.OrderChannelMetadata{},

		// Order Items
		&models.OrderItem{},

		// Addresses
		&models.Address{},

		// Payments
		&models.Payment{},

		// Shipments
		&models.Shipment{},
	); err != nil {
		return err
	}

	return r.seedInitialData(ctx)
}

func (r *Repository) seedInitialData(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Seed Scope "platform"
	var platformScope models.Scope
	if err := db.Where("code = ?", "platform").FirstOrCreate(&platformScope, models.Scope{
		Name:        "Platform",
		Code:        "platform",
		Description: "Scope for platform-wide permissions",
		IsSystem:    true,
	}).Error; err != nil {
		return fmt.Errorf("failed to seed platform scope: %w", err)
	}

	// 2. Seed Role "Super Admin"
	var superAdminRole models.Role
	if err := db.Where("name = ? AND scope_id = ?", "Super Admin", platformScope.ID).FirstOrCreate(&superAdminRole, models.Role{
		Name:        "Super Admin",
		Description: "Super Administrator with full access",
		Level:       1,
		IsSystem:    true,
		ScopeID:     platformScope.ID,
	}).Error; err != nil {
		return fmt.Errorf("failed to seed super admin role: %w", err)
	}

	// 3. Seed Default User
	var user models.User
	var count int64
	if err := db.Model(&models.User{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check for existing user: %w", err)
	}

	if count == 0 {
		email := r.cfg.Get("EMAIL_USER_DEFAULT")
		password := r.cfg.Get("USER_PASS_DEFAULT")

		if email == "" || password == "" {
			return fmt.Errorf("EMAIL_USER_DEFAULT or USER_PASS_DEFAULT not set")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user = models.User{
			Model:    gorm.Model{ID: 1},
			Name:     "Admin",
			Email:    email,
			Password: string(hashedPassword),
			IsActive: true,
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create default user: %w", err)
		}
	} else {
		// If user exists, fetch it to get the ID for association
		if err := db.First(&user, 1).Error; err != nil {
			return fmt.Errorf("failed to fetch existing user: %w", err)
		}
	}

	// 4. Assign Super Admin Role to User (in BusinessStaff with BusinessID = nil)
	var staffCount int64
	if err := db.Model(&models.BusinessStaff{}).
		Where("user_id = ? AND role_id = ? AND business_id IS NULL", user.ID, superAdminRole.ID).
		Count(&staffCount).Error; err != nil {
		return fmt.Errorf("failed to check for existing staff association: %w", err)
	}

	if staffCount == 0 {
		staff := models.BusinessStaff{
			UserID:     user.ID,
			BusinessID: nil, // Platform level
			RoleID:     &superAdminRole.ID,
		}
		if err := db.Create(&staff).Error; err != nil {
			return fmt.Errorf("failed to assign super admin role to user: %w", err)
		}
	}

	return nil
}
