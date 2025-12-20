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
		&models.OrderStatus{}, // Debe estar antes de BusinessNotificationConfigOrderStatus

		// Business Notification Config Order Status (tabla intermedia)
		// Debe ir después de BusinessNotificationConfig y OrderStatus
		&models.BusinessNotificationConfigOrderStatus{},

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

	return r.createDefaultUserIfNotExists(ctx)
}

// createDefaultUserIfNotExists crea el usuario principal solo si no existe
// Las migraciones solo deben crear la estructura de tablas, no datos adicionales
func (r *Repository) createDefaultUserIfNotExists(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// Verificar si el usuario principal ya existe
	var count int64
	if err := db.Model(&models.User{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check for existing user: %w", err)
	}

	// Solo crear el usuario si no existe
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

		user := models.User{
			Model:    gorm.Model{ID: 1},
			Name:     "Admin",
			Email:    email,
			Password: string(hashedPassword),
			IsActive: true,
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create default user: %w", err)
		}
	}

	return nil
}
