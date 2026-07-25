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
	Kind        string  `gorm:"size:15;not null;default:'CHARGE';index"`
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

type AccountingInvoice struct {
	gorm.Model
	Number     string         `gorm:"size:30;uniqueIndex"`
	BusinessID uint           `gorm:"not null;index"`
	ConceptID  uint           `gorm:"not null;index"`
	IssueDate  time.Time      `gorm:"type:date;not null"`
	DueDate    *time.Time     `gorm:"type:date"`
	Status     string         `gorm:"size:15;not null;default:'DRAFT';index"`
	Notes      string         `gorm:"size:1000"`
	Subtotal   float64        `gorm:"type:decimal(15,2);not null;default:0"`
	TaxTotal   float64        `gorm:"type:decimal(15,2);not null;default:0"`
	Total      float64        `gorm:"type:decimal(15,2);not null;default:0"`
	TaxDetail  datatypes.JSON `gorm:"type:jsonb"`
	EmailTo    string         `gorm:"size:255"`
	SentAt     *time.Time
	PaidAt     *time.Time

	CustomerDocument string `gorm:"size:30"`
	CustomerPhone    string `gorm:"size:30"`
	CustomerAddress  string `gorm:"size:255"`
	DianStatus       string `gorm:"size:15;not null;default:'NONE';index"`
	Cufe             string `gorm:"size:120"`
	DianNumber       string `gorm:"size:40"`
	DianExternalID   string `gorm:"size:40"`
	DianQR           string `gorm:"type:text"`
	DianEmittedAt    *time.Time
	IsTest           bool `gorm:"default:false;index"`

	WithholdingTotal  float64        `gorm:"type:decimal(15,2);not null;default:0"`
	WithholdingDetail datatypes.JSON `gorm:"type:jsonb"`

	Business Business                `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Concept  AccountingConcept       `gorm:"foreignKey:ConceptID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Items    []AccountingInvoiceItem `gorm:"foreignKey:InvoiceID"`
}

func (AccountingInvoice) TableName() string { return "accounting_invoices" }

type AccountingInvoiceItem struct {
	gorm.Model
	InvoiceID   uint    `gorm:"not null;index"`
	ServiceID   *uint   `gorm:"index"`
	UnspscCode  string  `gorm:"size:10"`
	Description string  `gorm:"size:300;not null"`
	Quantity    float64 `gorm:"type:decimal(12,2);not null;default:1"`
	UnitPrice   float64 `gorm:"type:decimal(15,2);not null"`
	Amount      float64 `gorm:"type:decimal(15,2);not null"`

	Invoice AccountingInvoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (AccountingInvoiceItem) TableName() string { return "accounting_invoice_items" }

type AccountingService struct {
	gorm.Model
	Code       string  `gorm:"size:30;not null;uniqueIndex"`
	Name       string  `gorm:"size:200;not null"`
	UnspscCode string  `gorm:"size:10"`
	UnitPrice  float64 `gorm:"type:decimal(15,2);not null;default:0"`
	IsActive   bool    `gorm:"default:true;index"`
}

func (AccountingService) TableName() string { return "accounting_services" }
