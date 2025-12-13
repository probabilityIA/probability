package main

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type BusinessType struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(255);not null"`
	Code      string `gorm:"type:varchar(50);not null"`
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (BusinessType) TableName() string {
	return "business_type" // Singular due to NamingStrategy
}

type Business struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"type:varchar(255);not null"`
	Code           string `gorm:"type:varchar(50);not null"`
	BusinessTypeID uint
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Business) TableName() string {
	return "business" // Singular due to NamingStrategy
}

// ... existing code ...

type PaymentMethod struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"size:64;unique;not null"`
	Name      string `gorm:"size:128;not null"`
	Category  string `gorm:"size:64"`
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}

func main() {
	ctx := context.Background()
	logger := log.New()
	environment := env.New(logger)
	database := db.New(logger, environment)

	// 1. Seed BusinessType
	var countType int64
	if err := database.Conn(ctx).Model(&BusinessType{}).Where("id = ?", 1).Count(&countType).Error; err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to count business types")
	}

	if countType == 0 {
		logger.Info(ctx).Msg("BusinessType ID 1 not found. Creating default type...")
		bt := BusinessType{
			ID:        1,
			Name:      "Default Type",
			Code:      "default",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := database.Conn(ctx).Create(&bt).Error; err != nil {
			logger.Fatal(ctx).Err(err).Msg("Failed to create default business type")
		}
		logger.Info(ctx).Msg("Successfully created BusinessType ID 1")
	}

	// 2. Seed Business
	var count int64
	if err := database.Conn(ctx).Model(&Business{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to count businesses")
	}

	if count == 0 {
		logger.Info(ctx).Msg("Business ID 1 not found. Creating default business...")
		business := Business{
			ID:             1,
			Name:           "Default Business",
			Code:           "default_biz",
			BusinessTypeID: 1, // Link to the created type
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := database.Conn(ctx).Create(&business).Error; err != nil {
			logger.Fatal(ctx).Err(err).Msg("Failed to create default business")
		}
		logger.Info(ctx).Msg("Successfully created Business ID 1: Default Business")
	} else {
		logger.Info(ctx).Msg("Business ID 1 already exists. Skipping.")
	}

	// 3. Seed PaymentMethod
	var countPM int64
	if err := database.Conn(ctx).Model(&PaymentMethod{}).Where("id = ?", 1).Count(&countPM).Error; err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to count payment methods")
	}

	if countPM == 0 {
		logger.Info(ctx).Msg("PaymentMethod ID 1 not found. Creating default method...")
		pm := PaymentMethod{
			ID:        1,
			Code:      "default_payment",
			Name:      "Default Payment",
			Category:  "cash",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := database.Conn(ctx).Create(&pm).Error; err != nil {
			logger.Fatal(ctx).Err(err).Msg("Failed to create default payment method")
		}
		logger.Info(ctx).Msg("Successfully created PaymentMethod ID 1")
	} else {
		logger.Info(ctx).Msg("PaymentMethod ID 1 already exists. Skipping.")
	}
}
