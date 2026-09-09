package entities

type StorefrontClient struct {
	ID         uint
	BusinessID uint
	UserID     *uint
	Name       string
	Email      *string
	Phone      string
	Dni        *string
}

type NewUser struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

type StorefrontBusiness struct {
	ID   uint
	Name string
	Code string
}
