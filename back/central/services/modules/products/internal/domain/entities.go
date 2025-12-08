package domain

import (
	"time"

	"gorm.io/datatypes"
)

// ───────────────────────────────────────────
//
//	PRODUCT DTOs
//
// ───────────────────────────────────────────

// CreateProductRequest representa la solicitud para crear un producto
type CreateProductRequest struct {
	// Identificadores
	BusinessID uint   `json:"business_id" binding:"required"`
	SKU        string `json:"sku" binding:"required,max=128"`
	ExternalID string `json:"external_id" binding:"omitempty,max=255"`

	// Información Básica
	Name             string `json:"name" binding:"required,max=255"`
	Title            string `json:"title" binding:"omitempty,max=500"`
	Description      string `json:"description" binding:"omitempty"`
	ShortDescription string `json:"short_description" binding:"omitempty,max=500"`
	Slug             string `json:"slug" binding:"omitempty,max=255"`

	// Pricing
	Price          float64  `json:"price" binding:"omitempty,min=0"`
	CompareAtPrice *float64 `json:"compare_at_price" binding:"omitempty,min=0"`
	CostPrice      *float64 `json:"cost_price" binding:"omitempty,min=0"`
	Currency       string   `json:"currency" binding:"omitempty,max=10"`

	// Inventory
	StockQuantity     int  `json:"stock_quantity" binding:"omitempty,min=0"`
	TrackInventory    bool `json:"track_inventory"`
	AllowBackorder    bool `json:"allow_backorder"`
	LowStockThreshold *int `json:"low_stock_threshold" binding:"omitempty,min=0"`

	// Media
	ImageURL string         `json:"image_url" binding:"omitempty,max=500"`
	Images   datatypes.JSON `json:"images" binding:"omitempty"`
	VideoURL *string        `json:"video_url" binding:"omitempty,max=500"`

	// Dimensiones y Peso
	Weight        *float64 `json:"weight" binding:"omitempty,min=0"`
	WeightUnit    string   `json:"weight_unit" binding:"omitempty,max=10"`
	Length        *float64 `json:"length" binding:"omitempty,min=0"`
	Width         *float64 `json:"width" binding:"omitempty,min=0"`
	Height        *float64 `json:"height" binding:"omitempty,min=0"`
	DimensionUnit string   `json:"dimension_unit" binding:"omitempty,max=10"`

	// Categorización
	Category string         `json:"category" binding:"omitempty,max=255"`
	Tags     datatypes.JSON `json:"tags" binding:"omitempty"`
	Brand    string         `json:"brand" binding:"omitempty,max=255"`

	// Estado
	Status     string `json:"status" binding:"omitempty,oneof=active draft archived"`
	IsActive   bool   `json:"is_active"`
	IsFeatured bool   `json:"is_featured"`

	// Metadata
	Metadata datatypes.JSON `json:"metadata" binding:"omitempty"`
}

// UpdateProductRequest representa la solicitud para actualizar un producto
type UpdateProductRequest struct {
	// Identificadores
	SKU        *string `json:"sku" binding:"omitempty,max=128"`
	ExternalID *string `json:"external_id" binding:"omitempty,max=255"`

	// Información Básica
	Name             *string `json:"name" binding:"omitempty,max=255"`
	Title            *string `json:"title" binding:"omitempty,max=500"`
	Description      *string `json:"description" binding:"omitempty"`
	ShortDescription *string `json:"short_description" binding:"omitempty,max=500"`
	Slug             *string `json:"slug" binding:"omitempty,max=255"`

	// Pricing
	Price          *float64 `json:"price" binding:"omitempty,min=0"`
	CompareAtPrice *float64 `json:"compare_at_price" binding:"omitempty,min=0"`
	CostPrice      *float64 `json:"cost_price" binding:"omitempty,min=0"`
	Currency       *string  `json:"currency" binding:"omitempty,max=10"`

	// Inventory
	StockQuantity     *int  `json:"stock_quantity" binding:"omitempty,min=0"`
	TrackInventory    *bool `json:"track_inventory"`
	AllowBackorder    *bool `json:"allow_backorder"`
	LowStockThreshold *int  `json:"low_stock_threshold" binding:"omitempty,min=0"`

	// Media
	ImageURL *string        `json:"image_url" binding:"omitempty,max=500"`
	Images   datatypes.JSON `json:"images" binding:"omitempty"`
	VideoURL *string        `json:"video_url" binding:"omitempty,max=500"`

	// Dimensiones y Peso
	Weight        *float64 `json:"weight" binding:"omitempty,min=0"`
	WeightUnit    *string  `json:"weight_unit" binding:"omitempty,max=10"`
	Length        *float64 `json:"length" binding:"omitempty,min=0"`
	Width         *float64 `json:"width" binding:"omitempty,min=0"`
	Height        *float64 `json:"height" binding:"omitempty,min=0"`
	DimensionUnit *string  `json:"dimension_unit" binding:"omitempty,max=10"`

	// Categorización
	Category *string        `json:"category" binding:"omitempty,max=255"`
	Tags     datatypes.JSON `json:"tags" binding:"omitempty"`
	Brand    *string        `json:"brand" binding:"omitempty,max=255"`

	// Estado
	Status     *string `json:"status" binding:"omitempty,oneof=active draft archived"`
	IsActive   *bool   `json:"is_active"`
	IsFeatured *bool   `json:"is_featured"`

	// Metadata
	Metadata datatypes.JSON `json:"metadata" binding:"omitempty"`
}

// ProductResponse representa la respuesta de un producto
type ProductResponse struct {
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

// ProductsListResponse representa la respuesta paginada de productos
type ProductsListResponse struct {
	Data       []ProductResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// ───────────────────────────────────────────
//
//	PRODUCT INTEGRATION DTOs
//
// ───────────────────────────────────────────

// AddProductIntegrationRequest representa la solicitud para asociar un producto con una integración
type AddProductIntegrationRequest struct {
	IntegrationID     uint   `json:"integration_id" binding:"required"`
	ExternalProductID string `json:"external_product_id" binding:"required,max=255"`
}

// RemoveProductIntegrationRequest representa la solicitud para remover una integración de un producto
type RemoveProductIntegrationRequest struct {
	IntegrationID uint `json:"integration_id" binding:"required"`
}

// ProductIntegrationResponse representa la respuesta de una integración asociada a un producto
type ProductIntegrationResponse struct {
	ID                uint      `json:"id"`
	ProductID         string    `json:"product_id"`
	IntegrationID     uint      `json:"integration_id"`
	IntegrationType   string    `json:"integration_type,omitempty"`
	IntegrationName   string    `json:"integration_name,omitempty"`
	ExternalProductID string    `json:"external_product_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ProductResponseWithIntegrations extiende ProductResponse con información de integraciones
type ProductResponseWithIntegrations struct {
	ProductResponse
	Integrations []ProductIntegrationResponse `json:"integrations,omitempty"`
}
