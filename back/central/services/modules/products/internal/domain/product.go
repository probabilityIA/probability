package domain

import (
	"time"

	"gorm.io/datatypes"
)

// ProductFamily representa el producto padre o familia de variantes.
type ProductFamily struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	BusinessID uint

	Name         string
	Title        string
	Description  string
	Slug         string
	Category     string
	Brand        string
	ImageURL     string
	Status       string
	IsActive     bool
	VariantAxes  datatypes.JSON
	Metadata     datatypes.JSON
	VariantCount int64
}

// Product representa un producto en el dominio
type Product struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Identificadores
	BusinessID uint
	SKU        string
	ExternalID string
	Barcode    *string
	FamilyID   *uint

	// Información Básica
	Name              string
	Title             string
	Description       string
	ShortDescription  string
	Slug              string
	VariantLabel      string
	VariantAttributes datatypes.JSON
	VariantSignature  string

	// Pricing
	Price          float64
	CompareAtPrice *float64
	CostPrice      *float64
	Currency       string

	// Inventory
	StockQuantity     int
	TrackInventory    bool
	AllowBackorder    bool
	LowStockThreshold *int

	// Media
	ImageURL string
	Images   datatypes.JSON
	VideoURL *string

	// Dimensiones y Peso
	Weight        *float64
	WeightUnit    string
	Length        *float64
	Width         *float64
	Height        *float64
	DimensionUnit string

	// Categorización
	Category string
	Tags     datatypes.JSON
	Brand    string

	// Estado
	Status     string
	IsActive   bool
	IsFeatured bool

	// Metadata
	Metadata datatypes.JSON

	// Relaciones
	Family *ProductFamily
}

// ProductBusinessIntegration representa la asociación de un producto con una integración
type ProductBusinessIntegration struct {
	ID                uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
	ProductID         string
	BusinessID        uint
	IntegrationID     uint
	ExternalProductID string
	ExternalVariantID *string
	ExternalSKU       *string
	ExternalBarcode   *string

	// Información de la integración (opcional, se incluye cuando se hace Preload)
	IntegrationName string
	IntegrationType string
}
