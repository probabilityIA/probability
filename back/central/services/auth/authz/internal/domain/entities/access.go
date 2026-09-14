package entities

type Access struct {
	UserID             uint
	BusinessID         uint
	BusinessName       string
	IsSuper            bool
	RoleID             uint
	RoleName           string
	RoleCode           string
	SubscriptionStatus string
	Modules            []string
	Permissions        []string
}

type NavItem struct {
	Key         string
	Label       string
	Route       string
	Section     string
	Description string
}

type StaffRole struct {
	BusinessID         *uint
	BusinessName       string
	SubscriptionStatus string
	RoleID             *uint
	RoleName           string
	RoleScopeCode      string
}

type RolePermission struct {
	ResourceName string
	ActionName   string
}

type ResourceConfig struct {
	ResourceName string
	Active       bool
}
