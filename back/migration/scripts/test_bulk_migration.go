package main

import (
	"fmt"
	"log"

	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=probability port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate bulk job tables
	fmt.Println("Running AutoMigrate for BulkInvoiceJob and BulkInvoiceJobItem...")
	if err := db.AutoMigrate(&models.BulkInvoiceJob{}, &models.BulkInvoiceJobItem{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("✅ Migration completed successfully!")
}
