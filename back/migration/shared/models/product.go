package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Product representa un producto único en el catálogo del negocio
// REGLA DE NEGOCIO: Un producto SIEMPRE debe estar asociado a un BusinessID (not null)
// La combinación de BusinessID + SKU debe ser única en el sistema
type Product struct {
	// ID alfanumérico único (generado automáticamente)
	ID        string     `gorm:"type:varchar(64);primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// ID del negocio al que pertenece este producto (REQUERIDO)
	BusinessID uint `gorm:"not null;index;uniqueIndex:idx_business_product_sku,priority:1"`

	// SKU único dentro del negocio
	SKU string `gorm:"size:128;not null;uniqueIndex:idx_business_product_sku,priority:2"`

	// Información Básica
	Name             string `gorm:"size:255;not null" json:"name"`
	Title            string `gorm:"size:500" json:"title"`             // Título descriptivo largo
	Description      string `gorm:"type:text" json:"description"`      // Descripción detallada
	ShortDescription string `gorm:"size:500" json:"short_description"` // Descripción corta
	Slug             string `gorm:"size:255;index" json:"slug"`        // URL-friendly identifier
	ExternalID       string `gorm:"size:255;index" json:"external_id"` // ID en sistemas externos

	// Pricing
	Price          float64  `gorm:"type:decimal(15,2);default:0" json:"price"`            // Precio base
	CompareAtPrice *float64 `gorm:"type:decimal(15,2)" json:"compare_at_price,omitempty"` // Precio de comparación/tachado
	CostPrice      *float64 `gorm:"type:decimal(15,2)" json:"cost_price,omitempty"`       // Precio de costo
	Currency       string   `gorm:"size:10;default:'COP'" json:"currency"`                // Moneda

	// Inventory
	StockQuantity     int  `gorm:"default:0" json:"stock_quantity"`      // Cantidad en stock
	TrackInventory    bool `gorm:"default:false" json:"track_inventory"` // Si trackear inventario
	AllowBackorder    bool `gorm:"default:false" json:"allow_backorder"` // Permitir venta sin stock
	LowStockThreshold *int `json:"low_stock_threshold,omitempty"`        // Umbral de stock bajo

	// Media
	ImageURL string         `gorm:"size:500" json:"image_url"`           // URL de imagen principal
	Images   datatypes.JSON `gorm:"type:jsonb" json:"images,omitempty"`  // Array de URLs de imágenes
	VideoURL *string        `gorm:"size:500" json:"video_url,omitempty"` // URL de video

	// Dimensiones y Peso
	Weight        *float64 `gorm:"type:decimal(10,3)" json:"weight,omitempty"` // Peso
	WeightUnit    string   `gorm:"size:10;default:'kg'" json:"weight_unit"`    // Unidad de peso
	Length        *float64 `gorm:"type:decimal(10,2)" json:"length,omitempty"` // Largo
	Width         *float64 `gorm:"type:decimal(10,2)" json:"width,omitempty"`  // Ancho
	Height        *float64 `gorm:"type:decimal(10,2)" json:"height,omitempty"` // Alto
	DimensionUnit string   `gorm:"size:10;default:'cm'" json:"dimension_unit"` // Unidad de dimensiones

	// Categorización
	Category string         `gorm:"size:255;index" json:"category"`   // Categoría
	Tags     datatypes.JSON `gorm:"type:jsonb" json:"tags,omitempty"` // Array de tags
	Brand    string         `gorm:"size:255;index" json:"brand"`      // Marca

	// Estado
	Status     string `gorm:"size:50;default:'draft';index" json:"status"` // Estado: active, draft, archived
	IsActive   bool   `gorm:"default:false;index" json:"is_active"`        // Si está activo
	IsFeatured bool   `gorm:"default:false;index" json:"is_featured"`      // Si es destacado

	// Metadata
	Metadata datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"` // Datos adicionales flexibles

	// Relaciones
	Business                    Business                     `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProductBusinessIntegrations []ProductBusinessIntegration `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (Product) TableName() string {
	return "products"
}

// BeforeCreate genera el ID hash antes de crear el producto
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = generateProductID()
	}
	return nil
}

// generateProductID genera un ID alfanumérico único para productos
// Formato: PRD_ + 12 caracteres aleatorios (letras y números)
func generateProductID() string {
	// Generar 9 bytes aleatorios (suficiente para 12 caracteres en base64)
	b := make([]byte, 9)
	rand.Read(b)

	// Codificar en base64 URL-safe y tomar los primeros 12 caracteres
	encoded := base64.RawURLEncoding.EncodeToString(b)
	if len(encoded) > 12 {
		encoded = encoded[:12]
	}

	return fmt.Sprintf("PRD_%s", encoded)
}
