package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateFreeTrialSubscriptionPlans(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.SubscriptionType{},
		&models.BusinessSubscription{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate free/trial subscription fields: %w", err)
	}

	return r.seedFreeTrialSubscriptionTypes(ctx)
}

func (r *Repository) seedFreeTrialSubscriptionTypes(ctx context.Context) error {
	db := r.db.Conn(ctx)

	freeShipments := 50
	freeOveragePrice := 1000.0
	trialDays := 15

	types := []struct {
		Name                 string
		Code                 string
		Description          string
		Modules              []string
		MaxEcommerceChannels int
		Payable              bool
		IncludedShipments    *int
		ShipmentOveragePrice *float64
		TrialDurationDays    *int
	}{
		{
			Name:                 "Gratuito",
			Code:                 "free",
			Description:          "Plan gratuito con 50 envios incluidos por mes",
			Modules:              []string{"orders", "customers"},
			MaxEcommerceChannels: 1,
			Payable:              true,
			IncludedShipments:    &freeShipments,
			ShipmentOveragePrice: &freeOveragePrice,
		},
		{
			Name:        "Prueba",
			Code:        "trial",
			Description: "Plan de prueba con acceso a todos los modulos por 15 dias",
			Modules: []string{
				"iam", "orders", "shipments", "inventory", "invoicing", "delivery",
				"customers", "storefront", "wallet", "announcements", "tickets",
				"integrations", "notification_config",
			},
			MaxEcommerceChannels: 3,
			Payable:              false,
			TrialDurationDays:    &trialDays,
		},
	}

	for _, t := range types {
		var existing models.SubscriptionType
		result := db.Where("code = ?", t.Code).First(&existing)
		if result.RowsAffected > 0 {
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
			Price:                0,
			BillingPeriod:        "monthly",
			Active:               true,
			Features:             featuresJSON,
			MaxEcommerceChannels: t.MaxEcommerceChannels,
			Payable:              t.Payable,
			IncludedShipments:    t.IncludedShipments,
			ShipmentOveragePrice: t.ShipmentOveragePrice,
			TrialDurationDays:    t.TrialDurationDays,
		}
		if err := db.Create(&subType).Error; err != nil {
			return fmt.Errorf("failed to seed subscription type %s: %w", t.Code, err)
		}
	}

	return nil
}
