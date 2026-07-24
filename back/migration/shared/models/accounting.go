package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AccountingConcept struct {
	gorm.Model
	Code         string `gorm:"size:50;not null;uniqueIndex"`
	Name         string `gorm:"size:120;not null"`
	Description  string `gorm:"type:text"`
	Kind         string `gorm:"size:10;not null;index"`
	IsRealIncome bool   `gorm:"default:true"`
	IsAutomatic  bool   `gorm:"default:false"`
	SourceType   string `gorm:"size:50;index"`
	IsActive     bool   `gorm:"default:true;index"`

	Taxes []AccountingConceptTax `gorm:"foreignKey:ConceptID"`
}

func (AccountingConcept) TableName() string { return "accounting_concepts" }

type AccountingTax struct {
	gorm.Model
	Code        string  `gorm:"size:50;not null;uniqueIndex"`
	Name        string  `gorm:"size:120;not null"`
	Description string  `gorm:"type:text"`
	RatePercent float64 `gorm:"type:decimal(8,4);not null;default:0"`
	IsActive    bool    `gorm:"default:true;index"`
}

func (AccountingTax) TableName() string { return "accounting_taxes" }

type AccountingConceptTax struct {
	gorm.Model
	ConceptID uint `gorm:"not null;uniqueIndex:idx_acc_concept_tax,priority:1"`
	TaxID     uint `gorm:"not null;uniqueIndex:idx_acc_concept_tax,priority:2"`
	IsActive  bool `gorm:"default:true"`

	Concept AccountingConcept `gorm:"foreignKey:ConceptID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tax     AccountingTax     `gorm:"foreignKey:TaxID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (AccountingConceptTax) TableName() string { return "accounting_concept_taxes" }

type AccountingEntry struct {
	gorm.Model
	ConceptID   uint           `gorm:"not null;index"`
	BusinessID  *uint          `gorm:"index"`
	EntryDate   time.Time      `gorm:"type:date;not null;index"`
	Amount      float64        `gorm:"type:decimal(15,2);not null"`
	Kind        string         `gorm:"size:10;not null;index"`
	SourceType  string         `gorm:"size:50;not null;default:'MANUAL';index"`
	SourceID    string         `gorm:"size:64;not null;default:''"`
	Description string         `gorm:"size:500"`
	IsAutomatic bool           `gorm:"default:false"`
	TaxTotal    float64        `gorm:"type:decimal(15,2);not null;default:0"`
	TaxDetail   datatypes.JSON `gorm:"type:jsonb"`

	Concept AccountingConcept `gorm:"foreignKey:ConceptID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (AccountingEntry) TableName() string { return "accounting_entries" }
