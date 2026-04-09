package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerSummary struct {
	gorm.Model
	CustomerID        uint    `gorm:"not null;uniqueIndex:idx_customer_summary_biz_cust,priority:2"`
	BusinessID        uint    `gorm:"not null;uniqueIndex:idx_customer_summary_biz_cust,priority:1;index"`
	TotalOrders       int     `gorm:"not null;default:0"`
	DeliveredOrders   int     `gorm:"not null;default:0"`
	CancelledOrders   int     `gorm:"not null;default:0"`
	InProgressOrders  int     `gorm:"not null;default:0"`
	TotalSpent        float64 `gorm:"type:decimal(14,2);not null;default:0"`
	AvgTicket         float64 `gorm:"type:decimal(14,2);not null;default:0"`
	TotalPaidOrders   int     `gorm:"not null;default:0"`
	AvgDeliveryScore  float64 `gorm:"type:decimal(5,2);not null;default:0"`
	FirstOrderAt      *time.Time
	LastOrderAt       *time.Time
	PreferredPlatform string    `gorm:"size:50"`
	LastUpdatedAt     time.Time `gorm:"autoUpdateTime"`

	Client Client `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CustomerAddress struct {
	gorm.Model
	CustomerID uint   `gorm:"not null;index;uniqueIndex:idx_cust_addr_unique,priority:1"`
	BusinessID uint   `gorm:"not null;index;uniqueIndex:idx_cust_addr_unique,priority:2"`
	Street     string `gorm:"size:500;uniqueIndex:idx_cust_addr_unique,priority:3"`
	City       string `gorm:"size:128;uniqueIndex:idx_cust_addr_unique,priority:4"`
	State      string `gorm:"size:128;uniqueIndex:idx_cust_addr_unique,priority:5"`
	Country    string `gorm:"size:128;uniqueIndex:idx_cust_addr_unique,priority:6"`
	PostalCode string `gorm:"size:32;uniqueIndex:idx_cust_addr_unique,priority:7"`
	TimesUsed  int    `gorm:"not null;default:1"`
	LastUsedAt time.Time

	Client Client `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CustomerProductHistory struct {
	gorm.Model
	CustomerID    uint    `gorm:"not null;index;uniqueIndex:idx_cust_prod_unique,priority:1"`
	BusinessID    uint    `gorm:"not null;index;uniqueIndex:idx_cust_prod_unique,priority:2"`
	ProductID     string  `gorm:"size:64;not null;index;uniqueIndex:idx_cust_prod_unique,priority:3"`
	ProductName   string  `gorm:"size:255"`
	ProductSKU    string  `gorm:"size:255"`
	ProductImage  *string `gorm:"size:512"`
	TimesOrdered  int     `gorm:"not null;default:0"`
	TotalQuantity int     `gorm:"not null;default:0"`
	TotalSpent    float64 `gorm:"type:decimal(14,2);not null;default:0"`
	FirstOrderedAt time.Time
	LastOrderedAt  time.Time

	Client Client `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CustomerOrderItem struct {
	gorm.Model
	CustomerID  uint    `gorm:"not null;index"`
	BusinessID  uint    `gorm:"not null;index"`
	OrderID     string  `gorm:"size:36;not null;index"`
	OrderNumber string  `gorm:"size:128"`
	ProductID   *string `gorm:"size:64;index"`
	ProductName string  `gorm:"size:255"`
	ProductSKU  string  `gorm:"size:255"`
	ProductImage *string `gorm:"size:512"`
	Quantity    int     `gorm:"not null;default:0"`
	UnitPrice   float64 `gorm:"type:decimal(14,2);not null;default:0"`
	TotalPrice  float64 `gorm:"type:decimal(14,2);not null;default:0"`
	OrderStatus string  `gorm:"size:64"`
	OrderedAt   time.Time

	Client Client `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
