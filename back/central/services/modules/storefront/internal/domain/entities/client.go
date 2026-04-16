package entities

// StorefrontClient represents a client linked to a storefront user
type StorefrontClient struct {
	ID         uint
	BusinessID uint
	UserID     *uint
	Name       string
	Email      *string
	Phone      string
	Dni        *string
}

// NewUser represents user data for registration
type NewUser struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

// StorefrontBusiness represents minimal business info for registration
type StorefrontBusiness struct {
	ID   uint
	Name string
	Code string
}
