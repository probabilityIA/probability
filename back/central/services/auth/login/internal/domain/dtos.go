package domain

import (
	"time"

	"github.com/secamc93/probability/back/central/shared/jwt"
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Success               bool
	Message               string
	User                  UserInfo
	Token                 string
	RequirePasswordChange bool
	Businesses            []BusinessInfo
	Scope                 string // Scope del usuario (platform, business, etc.)
	IsSuperAdmin          bool   // Indica si es super admin (scope platform o scope_id 1)
}

type UserInfo struct {
	ID          uint
	Name        string
	Email       string
	Phone       string
	AvatarURL   string
	IsActive    bool
	LastLoginAt *time.Time
}

type UserAuthInfo struct {
	ID          uint
	Name        string
	Email       string
	Password    string
	Phone       string
	AvatarURL   string
	IsActive    bool
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type BusinessInfo struct {
	ID                 uint
	Name               string
	Code               string
	BusinessTypeID     uint
	BusinessType       BusinessTypeInfo
	Timezone           string
	Address            string
	Description        string
	LogoURL            string
	PrimaryColor       string
	SecondaryColor     string
	TertiaryColor      string
	QuaternaryColor    string
	NavbarImageURL     string
	CustomDomain       string
	IsActive           bool
	EnableDelivery     bool
	EnablePickup       bool
	EnableReservations bool
}

type BusinessInfoEntity struct {
	ID                 uint
	Name               string
	Code               string
	BusinessTypeID     uint
	BusinessTypeName   string
	BusinessTypeCode   string
	Timezone           string
	Address            string
	Description        string
	LogoURL            string
	PrimaryColor       string
	SecondaryColor     string
	TertiaryColor      string
	QuaternaryColor    string
	NavbarImageURL     string
	CustomDomain       string
	IsActive           bool
	EnableDelivery     bool
	EnablePickup       bool
	EnableReservations bool
}

type BusinessTypeInfo struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Icon        string
}

type Role struct {
	ID               uint
	Name             string
	Description      string
	Level            int
	IsSystem         bool
	ScopeID          uint
	ScopeName        string
	ScopeCode        string
	BusinessTypeID   uint
	BusinessTypeName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Permission struct {
	ID               uint
	Name             string
	Description      string
	Resource         string
	Action           string
	ResourceID       uint
	ActionID         uint
	ScopeID          uint
	BusinessTypeID   uint
	BusinessTypeName string
}

type APIKey struct {
	ID          uint
	UserID      uint
	BusinessID  uint
	CreatedByID uint
	Name        string
	Description string
	KeyHash     string
	Revoked     bool
	RevokedAt   *time.Time
	LastUsedAt  *time.Time
	RateLimit   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type APIKeyInfo struct {
	ID          uint
	UserID      uint
	BusinessID  uint
	Name        string
	Description string
	RateLimit   int
	CreatedAt   time.Time
}

type BusinessStaffRelation struct {
	UserID     uint
	BusinessID *uint
	RoleID     *uint
	Business   *BusinessInfoEntity
}

type UserRolesPermissionsResponse struct {
	Success          bool
	Message          string
	UserID           uint
	Email            string
	IsSuper          bool
	BusinessID       uint
	BusinessName     string
	BusinessTypeID   uint
	BusinessTypeName string
	Role             RoleInfo
	Permissions      []PermissionInfo
}
type PermissionInfo struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Resource    string
	Action      string
	Scope       string
	Active      bool // Indica si el recurso está activo para el business
}
type RoleInfo struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Level       int
	IsSystem    bool
	Scope       string
}
type ChangePasswordRequest struct {
	UserID          uint
	CurrentPassword string
	NewPassword     string
}
type ChangePasswordResponse struct {
	Success bool
	Message string
}
type GeneratePasswordRequest struct {
	UserID uint
}
type GeneratePasswordResponse struct {
	Success  bool
	Email    string
	Password string
	Message  string
}
type ValidateAPIKeyRequest struct {
	APIKey string
}
type ValidateAPIKeyResponse struct {
	Success    bool
	Message    string
	UserID     uint
	Email      string
	BusinessID uint
	Roles      []string
	APIKeyID   uint
}

type JWTClaims = jwt.JWTClaims
type GenerateAPIKeyRequest struct {
	UserID      uint
	BusinessID  uint
	Name        string
	Description string
	RequesterID uint
}
type GenerateAPIKeyResponse struct {
	Success    bool
	Message    string
	APIKey     string
	APIKeyInfo APIKeyInfo
}
