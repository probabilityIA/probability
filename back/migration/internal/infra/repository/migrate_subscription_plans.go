package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSubscriptionPlans(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.SubscriptionType{},
		&models.BusinessModuleOverride{},
		&models.BusinessSubscription{},
		&models.Business{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate subscription plan tables: %w", err)
	}

	return r.seedDefaultSubscriptionTypes(ctx)
}

func (r *Repository) seedDefaultSubscriptionTypes(ctx context.Context) error {
	db := r.db.Conn(ctx)

	types := []struct {
		Name                 string
		Code                 string
		Description          string
		Price                float64
		Modules              []string
		MaxEcommerceChannels int
	}{
		{
			Name:                 "Basico",
			Code:                 "basico",
			Description:          "Plan de entrada con los modulos esenciales",
			Price:                50000,
			Modules:              []string{"orders", "customers"},
			MaxEcommerceChannels: 1,
		},
		{
			Name:                 "Pro",
			Code:                 "pro",
			Description:          "Plan intermedio con logistica e inventario",
			Price:                120000,
			Modules:              []string{"orders", "shipments", "inventory", "invoicing", "customers", "wallet"},
			MaxEcommerceChannels: 2,
		},
		{
			Name:        "Enterprise",
			Code:        "enterprise",
			Description: "Plan completo con todos los modulos",
			Price:       250000,
			Modules: []string{
				"iam", "orders", "shipments", "inventory", "invoicing", "delivery",
				"customers", "storefront", "wallet", "announcements", "tickets",
				"integrations", "notification_config",
			},
			MaxEcommerceChannels: 3,
		},
	}

	for _, t := range types {
		var existing models.SubscriptionType
		result := db.Where("code = ?", t.Code).First(&existing)
		if result.RowsAffected > 0 {
			if existing.MaxEcommerceChannels == 0 {
				if err := db.Model(&existing).Update("max_ecommerce_channels", t.MaxEcommerceChannels).Error; err != nil {
					return fmt.Errorf("failed to backfill max_ecommerce_channels for %s: %w", t.Code, err)
				}
			}
			continue
		}

		featuresJSON, err := marshalModuleCodes(t.Modules)
		if err != nil {
			return fmt.Errorf("failed to marshal modules for %s: %w", t.Code, err)
		}

		subType := models.SubscriptionType{
			Name:                 t.Name,
			Code:                 t.Code,
			Description:          t.Description,
			Price:                t.Price,
			BillingPeriod:        "monthly",
			Active:               true,
			Features:             featuresJSON,
			MaxEcommerceChannels: t.MaxEcommerceChannels,
		}
		if err := db.Create(&subType).Error; err != nil {
			return fmt.Errorf("failed to seed subscription type %s: %w", t.Code, err)
		}
	}

	return nil
}
