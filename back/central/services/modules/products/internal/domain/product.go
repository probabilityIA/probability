package domain

import (
	"time"

	"gorm.io/datatypes"
)

// Product representa un producto en el dominio
type Product struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Identificadores
	BusinessID uint   `json:"business_id"`
	SKU        string `json:"sku"`
	ExternalID string `json:"external_id"`

	// Información Básica
	Name             string `json:"name"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Slug             string `json:"slug"`

	// Pricing
	Price          float64  `json:"price"`
	CompareAtPrice *float64 `json:"compare_at_price,omitempty"`
	CostPrice      *float64 `json:"cost_price,omitempty"`
	Currency       string   `json:"currency"`

	// Inventory
	StockQuantity     int  `json:"stock_quantity"`
	TrackInventory    bool `json:"track_inventory"`
	AllowBackorder    bool `json:"allow_backorder"`
	LowStockThreshold *int `json:"low_stock_threshold,omitempty"`

	// Media
	ImageURL string         `json:"image_url"`
	Images   datatypes.JSON `json:"images,omitempty"`
	VideoURL *string        `json:"video_url,omitempty"`

	// Dimensiones y Peso
	Weight        *float64 `json:"weight,omitempty"`
	WeightUnit    string   `json:"weight_unit"`
	Length        *float64 `json:"length,omitempty"`
	Width         *float64 `json:"width,omitempty"`
	Height        *float64 `json:"height,omitempty"`
	DimensionUnit string   `json:"dimension_unit"`

	// Categorización
	Category string         `json:"category"`
	Tags     datatypes.JSON `json:"tags,omitempty"`
	Brand    string         `json:"brand"`

	// Estado
	Status     string `json:"status"`
	IsActive   bool   `json:"is_active"`
	IsFeatured bool   `json:"is_featured"`

	// Metadata
	Metadata datatypes.JSON `json:"metadata,omitempty"`
}

// ProductBusinessIntegration representa la asociación de un producto con una integración
type ProductBusinessIntegration struct {
	ID                uint       `json:"id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	ProductID         string     `json:"product_id"`
	BusinessID        uint       `json:"business_id"`
	IntegrationID     uint       `json:"integration_id"`
	ExternalProductID string     `json:"external_product_id"`

	// Información de la integración (opcional, se incluye cuando se hace Preload)
	IntegrationName string `json:"integration_name,omitempty"`
	IntegrationType string `json:"integration_type,omitempty"`
}
