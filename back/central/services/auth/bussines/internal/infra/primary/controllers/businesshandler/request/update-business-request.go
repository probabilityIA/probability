package request

import "mime/multipart"

// UpdateBusinessRequest representa la solicitud para actualizar un negocio (todos los campos opcionales)
type UpdateBusinessRequest struct {
	Name           string `form:"name" json:"name"`
	Code           string `form:"code" json:"code"`
	BusinessTypeID uint   `form:"business_type_id" json:"business_type_id"`
	Timezone       string `form:"timezone" json:"timezone"`
	Address        string `form:"address" json:"address"`
	Description    string `form:"description" json:"description"`

	LogoFile        *multipart.FileHeader `form:"logo_file" json:"-"`
	PrimaryColor    string                `form:"primary_color" json:"primary_color"`
	SecondaryColor  string                `form:"secondary_color" json:"secondary_color"`
	TertiaryColor   string                `form:"tertiary_color" json:"tertiary_color"`
	QuaternaryColor string                `form:"quaternary_color" json:"quaternary_color"`
	NavbarImageFile *multipart.FileHeader `form:"navbar_image_file" json:"-"`
	CustomDomain    string                `form:"custom_domain" json:"custom_domain"`
	IsActive        bool                  `form:"is_active" json:"is_active"`

	EnableDelivery     bool `form:"enable_delivery" json:"enable_delivery"`
	EnablePickup       bool `form:"enable_pickup" json:"enable_pickup"`
	EnableReservations bool `form:"enable_reservations" json:"enable_reservations"`
}
