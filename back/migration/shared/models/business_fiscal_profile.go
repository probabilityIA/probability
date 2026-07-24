package models

import "gorm.io/gorm"

type BusinessFiscalProfile struct {
	gorm.Model
	BusinessID        uint    `gorm:"not null;uniqueIndex"`
	DocumentType      string  `gorm:"size:10;default:'NIT'"`
	DocumentNumber    string  `gorm:"size:30"`
	DV                string  `gorm:"size:2"`
	PersonType        string  `gorm:"size:10;default:'JURIDICA'"`
	TaxRegime         string  `gorm:"size:60"`
	MunicipalityID    string  `gorm:"size:10"`
	Address           string  `gorm:"size:255"`
	Phone             string  `gorm:"size:30"`
	BillingEmail      string  `gorm:"size:255"`
	AppliesIVA        bool    `gorm:"default:false"`
	IVARate           float64 `gorm:"type:decimal(8,4);default:19"`
	AppliesReteFuente bool    `gorm:"default:false"`
	ReteFuenteRate    float64 `gorm:"type:decimal(8,4);default:11"`
	AppliesReteICA    bool    `gorm:"default:false"`
	ReteICARate       float64 `gorm:"type:decimal(8,4);default:0.966"`
	Notes             string  `gorm:"size:500"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (BusinessFiscalProfile) TableName() string { return "business_fiscal_profiles" }
