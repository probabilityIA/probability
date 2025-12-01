package response

import "time"

// LoginResponse representa la respuesta simplificada del login
type LoginResponse struct {
	User                  UserInfo       `json:"user"`
	Token                 string         `json:"token"`
	RequirePasswordChange bool           `json:"require_password_change"`
	Businesses            []BusinessInfo `json:"businesses"`
	Scope                 string         `json:"scope"`          // Scope del usuario (platform, business, etc.)
	IsSuperAdmin          bool           `json:"is_super_admin"` // Indica si es super admin (scope platform o scope_id 1)
}

// UserInfo representa la información del usuario en la respuesta
type UserInfo struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	AvatarURL   string     `json:"avatar_url"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// RoleInfo representa la información del rol en la respuesta
type RoleInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PermissionInfo representa la información del permiso en la respuesta
type Action struct {
	Description string `json:"description"`
	Action      string `json:"action"`
}

// ResourcePermissions representa los permisos agrupados por recurso
type ResourcePermissions struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
	Active   bool     `json:"active"` // Indica si el recurso está activo para el business
}

// UserRolesPermissionsResponse representa la respuesta de roles y permisos del usuario
type UserRolesPermissionsResponse struct {
	IsSuper          bool                  `json:"is_super"`
	BusinessID       uint                  `json:"business_id"`
	BusinessName     string                `json:"business_name"`
	BusinessTypeID   uint                  `json:"business_type_id"`
	BusinessTypeName string                `json:"business_type_name"`
	Role             RoleInfo              `json:"role"`
	Resources        []ResourcePermissions `json:"resources"`
}

// LoginSuccessResponse representa la respuesta exitosa del login para Swagger
type LoginSuccessResponse struct {
	Success bool          `json:"success"`
	Data    LoginResponse `json:"data"`
}

// UserRolesPermissionsSuccessResponse representa la respuesta exitosa de roles y permisos
type UserRolesPermissionsSuccessResponse struct {
	Success bool                         `json:"success"`
	Data    UserRolesPermissionsResponse `json:"data"`
}

// LoginErrorResponse representa la respuesta de error del login para Swagger
type LoginErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// LoginBadRequestResponse representa la respuesta de error 400 para Swagger
type LoginBadRequestResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

// BusinessInfo representa la información del negocio en la respuesta
type BusinessInfo struct {
	ID                 uint             `json:"id"`
	Name               string           `json:"name"`
	Code               string           `json:"code"`
	BusinessTypeID     uint             `json:"business_type_id"`
	BusinessType       BusinessTypeInfo `json:"business_type"`
	Timezone           string           `json:"timezone"`
	Address            string           `json:"address"`
	Description        string           `json:"description"`
	LogoURL            string           `json:"logo_url"`
	PrimaryColor       string           `json:"primary_color"`
	SecondaryColor     string           `json:"secondary_color"`
	TertiaryColor      string           `json:"tertiary_color"`
	QuaternaryColor    string           `json:"quaternary_color"`
	NavbarImageURL     string           `json:"navbar_image_url"`
	CustomDomain       string           `json:"custom_domain"`
	IsActive           bool             `json:"is_active"`
	EnableDelivery     bool             `json:"enable_delivery"`
	EnablePickup       bool             `json:"enable_pickup"`
	EnableReservations bool             `json:"enable_reservations"`
}

// BusinessTypeInfo representa la información del tipo de negocio
type BusinessTypeInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
