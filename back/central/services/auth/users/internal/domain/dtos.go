package domain

import (
	"mime/multipart"
	"time"
)

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

type BusinessInfoEntity struct {
	ID                 uint
	Name               string
	Code               string
	BusinessTypeID     uint
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
	BusinessTypeName   string
	BusinessTypeCode   string
}

type BusinessRoleAssignmentDetailed struct {
	BusinessID   uint
	BusinessName string
	RoleID       uint
	RoleName     string
}
type BusinessRoleAssignment struct {
	BusinessID uint
	RoleID     uint
}

type UserQueryDTO struct {
	ID          uint
	Name        string
	Email       string
	Phone       string
	AvatarURL   string
	IsActive    bool
	LastLoginAt *time.Time
	ScopeID     *uint
	ScopeCode   string // "platform" o "business"
	ScopeName   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
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
type UsersEntity struct {
	ID          uint
	Name        string
	Email       string
	Password    string
	Phone       string
	AvatarURL   string
	IsActive    bool
	ScopeID     *uint // Scope del usuario: platform (1) o business (2)
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
type UserFilters struct {
	Page       int
	PageSize   int
	Name       string
	Email      string
	Phone      string
	UserIDs    []uint // Lista de IDs de usuarios
	IsActive   *bool
	RoleID     *uint
	BusinessID *uint
	CreatedAt  string // formato: "2024-01-01" o "2024-01-01,2024-12-31"
	SortBy     string // "id", "name", "email", "created_at", etc.
	SortOrder  string // "asc" o "desc"

	// Filtros de scope para control de acceso
	ScopeID         *uint  // Filtrar usuarios por scope específico
	ScopeCode       string // Filtrar por código de scope (platform, business)
	RequesterScope  string // Scope del usuario que hace la solicitud
	RequesterUserID uint   // ID del usuario que hace la solicitud
}
type UserListDTO struct {
	Users      []UserDTO
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
type UserDTO struct {
	ID                      uint
	Name                    string
	Email                   string
	Phone                   string
	AvatarURL               string
	IsActive                bool
	LastLoginAt             *time.Time
	ScopeID                 *uint
	ScopeCode               string // "platform" o "business"
	ScopeName               string
	IsSuperUser             bool                             // Indica si es super usuario (scope platform)
	BusinessRoleAssignments []BusinessRoleAssignmentDetailed // Parejas business-rol con información completa
	Roles                   []RoleDTO                        // Mantener por compatibilidad
	Businesses              []BusinessDTO                    // Mantener por compatibilidad
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
}
type RoleDTO struct {
	ID               uint
	Name             string
	Code             string
	Description      string
	Level            int
	IsSystem         bool
	ScopeID          uint
	ScopeName        string // Nombre del scope para mostrar
	ScopeCode        string // Código del scope para mostrar
	BusinessTypeID   uint   // ID del tipo de business
	BusinessTypeName string // Nombre del tipo de business
}
type BusinessDTO struct {
	ID                 uint
	Name               string
	Code               string
	BusinessTypeID     uint
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
	BusinessTypeName   string
	BusinessTypeCode   string
	Role               *RoleDTO // Rol del usuario en este business (desde business_staff)
}

type CreateUserDTO struct {
	Name        string
	Email       string
	Password    string
	Phone       string
	AvatarURL   string                // URL completa (para compatibilidad)
	AvatarFile  *multipart.FileHeader // Archivo de imagen para subir a S3
	IsActive    bool
	ScopeID     *uint  // Scope del usuario (platform=1, business=2)
	BusinessIDs []uint // Businesses a relacionar con el usuario (obligatorio si scope es business)
}

type UpdateUserDTO struct {
	Name         string
	Email        string
	Phone        string
	AvatarURL    string                // URL completa (para compatibilidad)
	AvatarFile   *multipart.FileHeader // Archivo de imagen para subir a S3
	RemoveAvatar bool
	IsActive     bool
	BusinessIDs  []uint // Businesses a mantener (sobrescribe relaciones)
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
	LastUsedAt  *time.Time
	Revoked     bool
	RevokedAt   *time.Time
	RateLimit   int
	IPWhitelist string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type Resource struct {
	ID               uint
	Name             string
	Description      string
	BusinessTypeID   uint
	BusinessTypeName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
type APIKeyInfo struct {
	ID          uint
	UserID      uint
	BusinessID  uint
	Name        string
	Description string
	LastUsedAt  *time.Time
	Revoked     bool
	RateLimit   int
	CreatedAt   time.Time
}
type BusinessStaffRelation struct {
	UserID     uint
	BusinessID *uint               // NULL para super usuarios
	RoleID     *uint               // NULL si aún no tiene rol asignado
	Business   *BusinessInfoEntity // Info del business si business_id no es NULL
}
