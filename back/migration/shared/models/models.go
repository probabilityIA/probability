package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BusinessType struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"`
	Code        string `gorm:"size:50;not null;unique"`
	Description string `gorm:"size:500"`
	Icon        string `gorm:"size:100"`
	IsActive    bool   `gorm:"default:true"`

	Businesses []Business

	Roles []Role `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Resources []Resource `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Permissions []Permission `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Scope struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"`
	Code        string `gorm:"size:50;not null;unique"`
	Description string `gorm:"size:500"`
	IsSystem    bool   `gorm:"default:false"`

	Roles       []Role       `gorm:"foreignKey:ScopeID"`
	Permissions []Permission `gorm:"foreignKey:ScopeID"`
}

type Business struct {
	gorm.Model
	Name             string `gorm:"size:120;not null"`
	Code             string `gorm:"size:50;not null;unique"`
	BusinessTypeID   uint   `gorm:"not null;index"`
	ParentBusinessID *uint  `gorm:"index"`
	Timezone         string `gorm:"size:40;default:'America/Bogota'"`
	Address          string `gorm:"size:255"`
	Description      string `gorm:"size:500"`

	LogoURL         string  `gorm:"size:255"`
	PrimaryColor    string  `gorm:"size:7;default:'#1f2937'"`
	SecondaryColor  string  `gorm:"size:7;default:'#3b82f6'"`
	TertiaryColor   string  `gorm:"size:7;default:'#10b981'"`
	QuaternaryColor string  `gorm:"size:7;default:'#fbbf24'"`
	SidebarColor    string  `gorm:"size:7"`
	TopbarColor     string  `gorm:"size:7"`
	NavbarImageURL  string  `gorm:"size:255"`
	CustomDomain    *string `gorm:"size:100;unique"`
	IsActive        bool    `gorm:"default:true"`
	IsDemo          bool    `gorm:"default:false;index"`

	OrderPrefix string `gorm:"size:8;index"`

	SubscriptionStatus             string `gorm:"size:20;default:'active'"`
	SubscriptionEndDate            *time.Time
	SubscriptionCutoffDay          *int              `gorm:"column:subscription_cutoff_day"`
	SubscriptionCourtesyUntil      *time.Time        `gorm:"column:subscription_courtesy_until"`
	SubscriptionTypeID             *uint             `gorm:"index"`
	SubscriptionType               *SubscriptionType `gorm:"foreignKey:SubscriptionTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	SubscriptionAutoPaymentEnabled bool              `gorm:"column:subscription_auto_payment_enabled;default:false"`

	EnableDelivery     bool `gorm:"default:false"`
	EnablePickup       bool `gorm:"default:false"`
	EnableReservations bool `gorm:"default:true"`

	RequiresOrderConfirmation bool   `gorm:"default:false"`
	ConfirmationMethod        string `gorm:"default:'whatsapp'"`

	BusinessType                BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ParentBusiness              *Business    `gorm:"foreignKey:ParentBusinessID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ChildBusinesses             []Business   `gorm:"foreignKey:ParentBusinessID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Staff                       []BusinessStaff
	Clients                     []Client
	Users                       []User                       `gorm:"many2many:user_businesses;"`
	BusinessResourcesConfigured []BusinessResourceConfigured `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Integrations                []Integration                `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type BusinessResourceConfigured struct {
	gorm.Model
	BusinessID uint `gorm:"not null;index;uniqueIndex:idx_business_resource_config,priority:1"`
	ResourceID uint `gorm:"not null;index;uniqueIndex:idx_business_resource_config,priority:2"`
	Active     bool `gorm:"default:true"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Resource Resource `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Resource struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"`
	Description string `gorm:"size:500"`

	BusinessTypeID *uint         `gorm:"index"`
	BusinessType   *BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	BusinessResourcesConfigured []BusinessResourceConfigured `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permissions                 []Permission                 `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Role struct {
	gorm.Model
	Name        string `gorm:"size:50;not null;unique"`
	Description string `gorm:"size:255"`
	Level       int    `gorm:"not null;default:1"`
	IsSystem    bool   `gorm:"default:false"`

	ScopeID uint  `gorm:"not null;index"`
	Scope   Scope `gorm:"foreignKey:ScopeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	BusinessTypeID *uint         `gorm:"index"`
	BusinessType   *BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Users       []User       `gorm:"many2many:user_roles;"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"size:50;unique"`
	Description string `gorm:"size:500"`
	ResourceID  uint   `gorm:"not null;index"`
	ActionID    uint   `gorm:"not null;index"`
	ScopeID     uint   `gorm:"not null;index"`

	BusinessTypeID *uint         `gorm:"index"`
	BusinessType   *BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Scope    Scope    `gorm:"foreignKey:ScopeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Roles    []Role   `gorm:"many2many:role_permissions;"`
	Resource Resource `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Action   Action   `gorm:"foreignKey:ActionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type User struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	Email       string `gorm:"size:255;not null;unique"`
	Password    string `gorm:"size:255;not null"`
	Phone       string `gorm:"size:20"`
	AvatarURL   string `gorm:"size:255"`
	IsActive    bool   `gorm:"default:true"`
	LastLoginAt *time.Time

	GoogleID *string `gorm:"size:64;index"`

	ScopeID *uint  `gorm:"index"`
	Scope   *Scope `gorm:"foreignKey:ScopeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Businesses []Business `gorm:"many2many:user_businesses;"`

	Roles []Role `gorm:"many2many:user_roles;"`

	StaffOf []BusinessStaff
}

type BusinessStaff struct {
	gorm.Model
	UserID     uint  `gorm:"not null;index;uniqueIndex:idx_user_business,priority:1"`
	BusinessID *uint `gorm:"index;uniqueIndex:idx_user_business,priority:2"`

	RoleID *uint `gorm:"index"`

	User     User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role     Role     `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type Client struct {
	gorm.Model
	BusinessID uint    `gorm:"not null;index;uniqueIndex:idx_business_client_email,priority:1;uniqueIndex:idx_business_client_dni,priority:1"`
	Name       string  `gorm:"size:255;not null"`
	Email      *string `gorm:"size:255;uniqueIndex:idx_business_client_email,priority:2"`
	Phone      string  `gorm:"size:50"`
	Dni        *string `gorm:"size:50;uniqueIndex:idx_business_client_dni,priority:2"`
	UserID     *uint   `gorm:"index:idx_client_user_id"`

	AcceptsMarketing   bool       `gorm:"not null;default:true;index"`
	MarketingOptOutAt  *time.Time
	MarketingOptOutVia string `gorm:"size:32"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User     *User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Action struct {
	gorm.Model
	Name        string `gorm:"size:20;not null;unique"`
	Description string `gorm:"size:255"`

	Permissions []Permission `gorm:"foreignKey:ActionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type APIKey struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index"`
	BusinessID  uint   `gorm:"not null;index"`
	CreatedByID uint   `gorm:"not null;index"`
	Name        string `gorm:"size:255;not null"`
	KeyHash     string `gorm:"size:255;not null"`
	Description string `gorm:"size:500"`

	LastUsedAt *time.Time `gorm:"index"`
	Revoked    bool       `gorm:"default:false;index"`
	RevokedAt  *time.Time

	RateLimit   int    `gorm:"default:1000"`
	IPWhitelist string `gorm:"size:1000"`

	User      User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business  Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedBy User     `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type IntegrationCategory struct {
	gorm.Model
	Code             string `gorm:"size:50;not null;unique;index"`
	Name             string `gorm:"size:100;not null"`
	Description      string `gorm:"size:500"`
	Icon             string `gorm:"size:100"`
	Color            string `gorm:"size:20"`
	DisplayOrder     int    `gorm:"default:0"`
	ParentCategoryID *uint  `gorm:"index"`
	IsActive         bool   `gorm:"default:true;index"`
	IsVisible        bool   `gorm:"default:true"`

	ParentCategory   *IntegrationCategory `gorm:"foreignKey:ParentCategoryID"`
	IntegrationTypes []IntegrationType    `gorm:"foreignKey:CategoryID"`
}

func (IntegrationCategory) TableName() string {
	return "integration_categories"
}

type IntegrationType struct {
	gorm.Model
	Name             string `gorm:"size:100;not null;unique"`
	Code             string `gorm:"size:50;not null;unique"`
	Description      string `gorm:"size:500"`
	Icon             string `gorm:"size:100"`
	ImageURL         string `gorm:"size:500"`
	IsActive         bool   `gorm:"default:true"`
	InDevelopment    bool   `gorm:"default:false"`
	IsSystemProvider bool   `gorm:"default:false;index"`

	CategoryID *uint                `gorm:"index"`
	Category   *IntegrationCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	ConfigSchema datatypes.JSON `gorm:"type:jsonb"`

	CredentialsSchema datatypes.JSON `gorm:"type:jsonb"`

	SetupInstructions string `gorm:"type:text"`

	BaseURL     string `gorm:"column:base_url;size:500"`
	BaseURLTest string `gorm:"column:base_url_test;size:500"`

	PlatformCredentialsEncrypted []byte `gorm:"column:platform_credentials_encrypted;type:bytea"`

	ProductMatchOptions datatypes.JSON `gorm:"column:product_match_options;type:jsonb"`

	DefaultProductMatchRules datatypes.JSON `gorm:"column:default_product_match_rules;type:jsonb"`

	Integrations []Integration `gorm:"foreignKey:IntegrationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (IntegrationType) TableName() string {
	return "integration_types"
}

type Integration struct {
	gorm.Model

	Name     string `gorm:"size:100;not null"`
	Code     string `gorm:"size:50;not null;unique"`
	Category string `gorm:"size:50;not null;index"`
	StoreID  string `gorm:"size:150;index"`

	IntegrationTypeID uint             `gorm:"not null;index"`
	IntegrationType   *IntegrationType `gorm:"foreignKey:IntegrationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	BusinessID *uint     `gorm:"index"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	IsActive  bool `gorm:"default:true;index"`
	IsDefault bool `gorm:"default:false;index"`
	IsTesting bool `gorm:"default:false;index"`

	Config datatypes.JSON `gorm:"type:jsonb"`

	Credentials datatypes.JSON `gorm:"type:jsonb"`

	ProductMatchRules datatypes.JSON `gorm:"column:product_match_rules;type:jsonb"`

	Description string `gorm:"size:500"`
	CreatedByID uint   `gorm:"index"`
	UpdatedByID *uint  `gorm:"index"`

	CreatedBy           User                            `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy           *User                           `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	NotificationConfigs []IntegrationNotificationConfig `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Integration) TableName() string {
	return "integrations"
}

type IntegrationNotificationConfig struct {
	gorm.Model

	IntegrationID uint        `gorm:"not null;index"`
	Integration   Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	NotificationType string `gorm:"size:20;not null;index"`

	IsActive bool `gorm:"default:true;index"`

	Conditions datatypes.JSON `gorm:"type:jsonb;not null"`

	Config datatypes.JSON `gorm:"type:jsonb"`

	Description string `gorm:"size:500"`

	Priority int `gorm:"default:0;index"`
}

func (IntegrationNotificationConfig) TableName() string {
	return "integration_notification_configs"
}
