package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	BUSINESS TYPES - Tipos de negocios
//
// ───────────────────────────────────────────
type BusinessType struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"`
	Code        string `gorm:"size:50;not null;unique"` // Código interno
	Description string `gorm:"size:500"`
	Icon        string `gorm:"size:100"` // Icono para UI
	IsActive    bool   `gorm:"default:true"`

	// Relación con negocios
	Businesses []Business

	// Relación con roles (un tipo de business puede tener múltiples roles)
	Roles []Role `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// Relación con recursos (un tipo de business puede tener múltiples recursos)
	Resources []Resource `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// Relación con permisos (un tipo de business puede tener múltiples permisos)
	Permissions []Permission `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// ───────────────────────────────────────────
//
//	SCOPES - Ámbitos de permisos y roles
//
// ───────────────────────────────────────────
type Scope struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"`
	Code        string `gorm:"size:50;not null;unique"` // Código interno
	Description string `gorm:"size:500"`
	IsSystem    bool   `gorm:"default:false"` // Si es scope del sistema (no se puede eliminar)

	// Relaciones
	Roles       []Role       `gorm:"foreignKey:ScopeID"`
	Permissions []Permission `gorm:"foreignKey:ScopeID"`
}

// ───────────────────────────────────────────
//
//	BUSINESSES  (multi-tenant) - MARCA BLANCA
//
// ───────────────────────────────────────────
type Business struct {
	gorm.Model
	Name             string `gorm:"size:120;not null"`
	Code             string `gorm:"size:50;not null;unique"` // slug para URL personalizada
	BusinessTypeID   uint   `gorm:"not null;index"`
	ParentBusinessID *uint  `gorm:"index"` // ID del negocio padre (para jerarquía)
	Timezone         string `gorm:"size:40;default:'America/Bogota'"`
	Address          string `gorm:"size:255"`
	Description      string `gorm:"size:500"`

	// Configuración de marca blanca
	LogoURL         string  `gorm:"size:255"`
	PrimaryColor    string  `gorm:"size:7;default:'#1f2937'"` // Hex color
	SecondaryColor  string  `gorm:"size:7;default:'#3b82f6'"` // Hex color
	TertiaryColor   string  `gorm:"size:7;default:'#10b981'"` // Hex color adicional
	QuaternaryColor string  `gorm:"size:7;default:'#fbbf24'"` // Hex color adicional
	NavbarImageURL  string  `gorm:"size:255"`                 // Imagen de fondo para la barra de navegación
	CustomDomain    *string `gorm:"size:100;unique"`          // dominio personalizado
	IsActive        bool    `gorm:"default:true"`

	// Configuración de funcionalidades
	EnableDelivery     bool `gorm:"default:false"`
	EnablePickup       bool `gorm:"default:false"`
	EnableReservations bool `gorm:"default:true"`

	// Configuración de confirmación de órdenes
	RequiresOrderConfirmation bool   `gorm:"default:false"`      // Si requiere confirmación automática
	ConfirmationMethod        string `gorm:"default:'whatsapp'"` // whatsapp, email, sms

	// Relaciones
	BusinessType                BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ParentBusiness              *Business    `gorm:"foreignKey:ParentBusinessID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Negocio padre
	ChildBusinesses             []Business   `gorm:"foreignKey:ParentBusinessID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Negocios hijos
	Staff                       []BusinessStaff
	Clients                     []Client
	Users                       []User                       `gorm:"many2many:user_businesses;"` // Usuarios del negocio (muchos a muchos)
	BusinessResourcesConfigured []BusinessResourceConfigured `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Integrations                []Integration                `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Integraciones del negocio
}

// ───────────────────────────────────────────
//
//	BUSINESS RESOURCE CONFIGURED – recursos del negocio configurados para un negocio
//
// ───────────────────────────────────────────
type BusinessResourceConfigured struct {
	gorm.Model
	BusinessID uint `gorm:"not null;index;uniqueIndex:idx_business_resource_config,priority:1"`
	ResourceID uint `gorm:"not null;index;uniqueIndex:idx_business_resource_config,priority:2"`
	Active     bool `gorm:"default:true"`

	// Relaciones
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Resource Resource `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// ───────────────────────────────────────────
//
//	RESOURCES – recursos del negocio
//
// ───────────────────────────────────────────
type Resource struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"`
	Description string `gorm:"size:500"`

	// Relación con tipo de business (opcional para recursos genéricos)
	BusinessTypeID *uint         `gorm:"index"`                                                                   // Tipo de business (null = genérico, aplica a todos)
	BusinessType   *BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Relación con tipo de business

	// Relaciones
	BusinessResourcesConfigured []BusinessResourceConfigured `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permissions                 []Permission                 `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// ───────────────────────────────────────────
//
//	ROLES DEL SISTEMA
//
// ───────────────────────────────────────────
type Role struct {
	gorm.Model
	Name        string `gorm:"size:50;not null;unique"`
	Description string `gorm:"size:255"`
	Level       int    `gorm:"not null;default:1"` // Nivel jerárquico (1=super, 2=admin, 3=manager, 4=staff)
	IsSystem    bool   `gorm:"default:false"`      // Si es rol del sistema (no se puede eliminar)

	// Scope del rol
	ScopeID uint  `gorm:"not null;index"`
	Scope   Scope `gorm:"foreignKey:ScopeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Relación con tipo de business (un rol solo puede estar en un tipo de business)
	BusinessTypeID *uint         `gorm:"index"` // Temporalmente opcional para migración
	BusinessType   *BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Users       []User       `gorm:"many2many:user_roles;"`
}

// ───────────────────────────────────────────
//
//	PERMISOS DEL SISTEMA
//
// ───────────────────────────────────────────
type Permission struct {
	gorm.Model
	Name        string `gorm:"size:50;unique"`
	Description string `gorm:"size:500"`
	ResourceID  uint   `gorm:"not null;index"`
	ActionID    uint   `gorm:"not null;index"`
	ScopeID     uint   `gorm:"not null;index"`

	// Relación con tipo de business (opcional para permisos genéricos)
	BusinessTypeID *uint         `gorm:"index"`                                                                   // Tipo de business (null = genérico, aplica a todos)
	BusinessType   *BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // Relación con tipo de business

	Scope    Scope    `gorm:"foreignKey:ScopeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Roles    []Role   `gorm:"many2many:role_permissions;"`
	Resource Resource `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Action   Action   `gorm:"foreignKey:ActionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// ───────────────────────────────────────────
//
//	USUARIOS DEL SISTEMA
//
// ───────────────────────────────────────────
type User struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	Email       string `gorm:"size:255;not null;unique"`
	Password    string `gorm:"size:255;not null"`
	Phone       string `gorm:"size:20"`
	AvatarURL   string `gorm:"size:255"`
	IsActive    bool   `gorm:"default:true"`
	LastLoginAt *time.Time

	// Scope del usuario: platform (super admin) o business (usuario de negocio)
	// Si ScopeID es NULL, se asume "business" por defecto
	ScopeID *uint  `gorm:"index"`
	Scope   *Scope `gorm:"foreignKey:ScopeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Relación con negocios (un usuario puede estar en múltiples negocios)
	Businesses []Business `gorm:"many2many:user_businesses;"`

	// Roles del usuario (RELACIÓN MANY-TO-MANY)
	Roles []Role `gorm:"many2many:user_roles;"`

	// Relaciones existentes
	StaffOf []BusinessStaff
}

// ───────────────────────────────────────────
//
//	BUSINESS STAFF  (N:M usuario ↔ negocio)
//
// ───────────────────────────────────────────
type BusinessStaff struct {
	gorm.Model
	UserID     uint  `gorm:"not null;index;uniqueIndex:idx_user_business,priority:1"`
	BusinessID *uint `gorm:"index;uniqueIndex:idx_user_business,priority:2"`
	// Rol puede asignarse después: permitir NULL
	RoleID *uint `gorm:"index"` // Referencia a Role; NULL si aún no tiene rol asignado

	User     User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role     Role     `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// ───────────────────────────────────────────
//
//	CLIENTS – personas que hacen la reserva
//
// ───────────────────────────────────────────
type Client struct {
	gorm.Model
	BusinessID uint    `gorm:"not null;index;uniqueIndex:idx_business_client_email,priority:1"`
	Name       string  `gorm:"size:255;not null"`
	Email      string  `gorm:"size:255;uniqueIndex:idx_business_client_email,priority:2"`
	Phone      string  `gorm:"size:20"`
	Dni        *string `gorm:"size:30;uniqueIndex:idx_business_client_dni,priority:2"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// ───────────────────────────────────────────
//
//	ACTIONS – acciones que se pueden realizar en el sistema
//
// ───────────────────────────────────────────
type Action struct {
	gorm.Model
	Name        string `gorm:"size:20;not null;unique"`
	Description string `gorm:"size:255"`

	// Relaciones
	Permissions []Permission `gorm:"foreignKey:ActionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// ───────────────────────────────────────────
//
//	API KEYS - Claves de API para integraciones
//
// ───────────────────────────────────────────
type APIKey struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index"`    // Usuario para el cual se genera la API Key
	BusinessID  uint   `gorm:"not null;index"`    // Business asociado
	CreatedByID uint   `gorm:"not null;index"`    // Usuario que creó la API Key (super admin)
	Name        string `gorm:"size:255;not null"` // Nombre de referencia (ej. "API para sitio web")
	KeyHash     string `gorm:"size:255;not null"` // Hash de la API Key (bcrypt)
	Description string `gorm:"size:500"`          // Descripción opcional

	// Control de uso
	LastUsedAt *time.Time `gorm:"index"`               // Última vez que se usó
	Revoked    bool       `gorm:"default:false;index"` // Si está revocada
	RevokedAt  *time.Time // Cuándo fue revocada

	// Configuración opcional
	RateLimit   int    `gorm:"default:1000"` // Límite de requests por hora
	IPWhitelist string `gorm:"size:1000"`    // IPs permitidas (separadas por coma)

	// Relaciones
	User      User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business  Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedBy User     `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// ───────────────────────────────────────────
//
//	INTEGRATION CATEGORIES - Categorías de integraciones
//
// ───────────────────────────────────────────
type IntegrationCategory struct {
	gorm.Model
	Code             string `gorm:"size:50;not null;unique;index"` // "ecommerce", "invoicing", "messaging"
	Name             string `gorm:"size:100;not null"`             // "E-commerce", "Facturación Electrónica"
	Description      string `gorm:"size:500"`                      // Descripción de la categoría
	Icon             string `gorm:"size:100"`                      // Icono para UI
	Color            string `gorm:"size:20"`                       // Color hexadecimal para UI
	DisplayOrder     int    `gorm:"default:0"`                     // Orden de visualización
	ParentCategoryID *uint  `gorm:"index"`                         // Para categorías anidadas (futuro)
	IsActive         bool   `gorm:"default:true;index"`            // Si la categoría está activa
	IsVisible        bool   `gorm:"default:true"`                  // Si se muestra en UI

	// Relaciones
	ParentCategory   *IntegrationCategory `gorm:"foreignKey:ParentCategoryID"`
	IntegrationTypes []IntegrationType    `gorm:"foreignKey:CategoryID"`
}

// TableName especifica el nombre de la tabla para IntegrationCategory
func (IntegrationCategory) TableName() string {
	return "integration_categories"
}

// ───────────────────────────────────────────
//
//	INTEGRATION TYPES - Tipos de integraciones disponibles
//
// ───────────────────────────────────────────
type IntegrationType struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"` // "WhatsApp", "Shopify", "Mercado Libre"
	Code        string `gorm:"size:50;not null;unique"`  // "whatsapp", "shopify", "mercado_libre"
	Description string `gorm:"size:500"`                 // Descripción del tipo de integración
	Icon        string `gorm:"size:100"`                 // Icono para UI
	ImageURL    string `gorm:"size:500"`                 // URL de la imagen del logo (path relativo en S3)
	IsActive    bool   `gorm:"default:true"`             // Si el tipo está activo y disponible

	// Relación con IntegrationCategory
	CategoryID *uint                `gorm:"index"`
	Category   *IntegrationCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Configuración requerida (JSON schema - define qué campos de config son necesarios)
	// Ejemplo: {"required_fields": ["phone_number_id"], "optional_fields": ["webhook_url"]}
	ConfigSchema datatypes.JSON `gorm:"type:jsonb"`

	// Credenciales requeridas (JSON schema - define qué credenciales son necesarias)
	// Ejemplo: {"required_fields": ["access_token"], "optional_fields": ["refresh_token"]}
	CredentialsSchema datatypes.JSON `gorm:"type:jsonb"`

	// Instrucciones paso a paso para configurar la integración
	SetupInstructions string `gorm:"type:text"`

	// Relaciones
	Integrations []Integration `gorm:"foreignKey:IntegrationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla para IntegrationType
func (IntegrationType) TableName() string {
	return "integration_types"
}

// ───────────────────────────────────────────
//
//	INTEGRATIONS – Integraciones del sistema (WhatsApp, Shopify, Mercado Libre, etc.)
//
// ───────────────────────────────────────────
type Integration struct {
	gorm.Model

	// Identificación
	Name     string `gorm:"size:100;not null"`       // "WhatsApp Principal", "Shopify Store 1"
	Code     string `gorm:"size:50;not null;unique"` // "whatsapp_platform", "shopify_store_1"
	Category string `gorm:"size:50;not null;index"`  // Código de categoría (derivado de IntegrationType.Category.Code)
	StoreID  string `gorm:"size:150;index"`          // Identificador externo (p.e. shop domain)

	// Relación con IntegrationType (obligatorio)
	IntegrationTypeID uint             `gorm:"not null;index"`
	IntegrationType   *IntegrationType `gorm:"foreignKey:IntegrationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Relación con Business
	// NULL = integración global (como WhatsApp - una sola para toda la plataforma)
	// NOT NULL = integración específica de un business (como Shopify - puede haber múltiples)
	BusinessID *uint     `gorm:"index"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Estado
	IsActive  bool `gorm:"default:true;index"`
	IsDefault bool `gorm:"default:false;index"` // Si es la integración por defecto para este tipo

	// Configuración (JSON flexible - no contiene información sensible)
	// Ejemplo WhatsApp: {"phone_number_id": "123", "webhook_url": "...", "template_language": "es"}
	// Ejemplo Shopify: {"store_name": "mi-tienda", "api_version": "2024-01"}
	Config datatypes.JSON `gorm:"type:jsonb"`

	// Credenciales encriptadas (JSON)
	// Contiene tokens, API keys, secrets encriptados
	// Ejemplo: {"access_token": "encrypted_value", "api_key": "encrypted_value"}
	Credentials datatypes.JSON `gorm:"type:jsonb"`

	// Metadata
	Description string `gorm:"size:500"`
	CreatedByID uint   `gorm:"index"`
	UpdatedByID *uint  `gorm:"index"`

	// Relaciones
	CreatedBy           User                            `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy           *User                           `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	NotificationConfigs []IntegrationNotificationConfig `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla para Integration
func (Integration) TableName() string {
	return "integrations"
}

// ───────────────────────────────────────────
//
//	INTEGRATION NOTIFICATION CONFIG - Configuraciones de notificaciones por integración
//
// ───────────────────────────────────────────
type IntegrationNotificationConfig struct {
	gorm.Model

	// Relación con Integration
	IntegrationID uint        `gorm:"not null;index"`
	Integration   Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Tipo de notificación
	// "whatsapp" | "email" | "sms"
	NotificationType string `gorm:"size:20;not null;index"`

	// Estado de la notificación
	IsActive bool `gorm:"default:true;index"`

	// Condiciones/Eventos que disparan la notificación
	// JSON que define cuándo se debe enviar la notificación
	// Estructura:
	//   {
	//     "trigger": "order.created" | "order.updated" | "order.status_changed",
	//     "statuses": ["pending", "processing"], // opcional, vacío = todos
	//     "payment_methods": [1, 3, 5],          // opcional, vacío = todos
	//     "source_integration_id": 2             // opcional, null = todas las integraciones
	//   }
	// Ejemplos:
	//   {"trigger": "order.created"}
	//   {"trigger": "order.status_changed", "statuses": ["delivered"]}
	//   {"trigger": "order.created", "payment_methods": [1], "source_integration_id": 2}
	Conditions datatypes.JSON `gorm:"type:jsonb;not null"`

	// Configuración adicional de la notificación
	// Contiene templates, destinatarios, configuración específica del canal
	// Ejemplo WhatsApp:
	//   {"template_id": "order_status_update", "language": "es", "recipient_type": "customer"}
	// Ejemplo Email:
	//   {"template": "order_confirmation", "subject": "Tu orden está en camino", "recipient_type": "customer"}
	// Ejemplo SMS:
	//   {"message_template": "Tu orden #{{order_number}} está en camino", "recipient_type": "customer"}
	Config datatypes.JSON `gorm:"type:jsonb"`

	// Descripción opcional
	Description string `gorm:"size:500"`

	// Prioridad (para cuando hay múltiples configuraciones que coinciden)
	// Mayor número = mayor prioridad
	Priority int `gorm:"default:0;index"`
}

// TableName especifica el nombre de la tabla para IntegrationNotificationConfig
func (IntegrationNotificationConfig) TableName() string {
	return "integration_notification_configs"
}
