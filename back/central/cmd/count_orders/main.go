package main

import (
	"context"
	"fmt"
	"os"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Order struct {
	ID            uint   `gorm:"primaryKey"`
	IntegrationID uint   `gorm:"index"`
	OrderNumber   string `gorm:"size:64"`
}

func (Order) TableName() string {
	return "orders"
}

func main() {
	logger := log.New()
	environment := env.New(logger)
	database := db.New(logger, environment)
	ctx := context.Background()

	var totalCount int64
	var integration2Count int64

	// Count strictly in the orders table
	if err := database.Conn(ctx).Model(&Order{}).Count(&totalCount).Error; err != nil {
		fmt.Printf("Error counting total orders: %v\n", err)
		os.Exit(1)
	}

	if err := database.Conn(ctx).Model(&Order{}).Where("integration_id = ?", 2).Count(&integration2Count).Error; err != nil {
		fmt.Printf("Error counting integration 2 orders: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n--- DATABASE ORDER COUNT ---\n")
	fmt.Printf("Total Orders in DB: %d\n", totalCount)
	fmt.Printf("Orders for Integration 2: %d\n", integration2Count)
	fmt.Printf("----------------------------\n")
}
