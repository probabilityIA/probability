package models

import (
	"time"

	"gorm.io/gorm"
)

type BusinessSubscription struct {
	gorm.Model
	BusinessID          uint    `gorm:"not null;index"`
	SubscriptionTypeID  *uint   `gorm:"index"`
	Months              *int    `gorm:"column:months"`
	Amount              float64 `gorm:"type:numeric(12,2);not null"`
	StartDate           time.Time
	EndDate             time.Time
	Status              string     `gorm:"size:20;default:'pending'"`
	PaymentReference    *string    `gorm:"size:255"`
	Notes               *string    `gorm:"type:text"`
	AutoPaymentEnabled  bool       `gorm:"default:true"`
	PaymentMethod       string     `gorm:"size:50;default:'WALLET'"`
	PaymentMode         *string    `gorm:"size:20"`
	OverageAccepted     bool       `gorm:"column:overage_accepted;not null;default:false"`
	OverageAcceptedAt   *time.Time `gorm:"column:overage_accepted_at"`
	OverageAmountDue    *float64   `gorm:"column:overage_amount_due;type:numeric(12,2)"`
	OverageAmountPaidAt *time.Time `gorm:"column:overage_amount_paid_at"`

	Business         Business          `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubscriptionType *SubscriptionType `gorm:"foreignKey:SubscriptionTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (BusinessSubscription) TableName() string {
	return "business_subscriptions"
}
