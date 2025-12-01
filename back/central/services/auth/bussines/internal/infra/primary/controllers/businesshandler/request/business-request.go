package request

import "mime/multipart"

// BusinessRequest representa la solicitud para crear/actualizar un negocio
type BusinessRequest struct {
	Name           string `form:"name" binding:"required"`
	Code           string `form:"code" binding:"required"`
	BusinessTypeID uint   `form:"business_type_id" binding:"required"`
	Timezone       string `form:"timezone"`
	Address        string `form:"address"`
	Description    string `form:"description"`

	// Configuración de marca blanca
	LogoFile        *multipart.FileHeader `form:"logo_file"`
	PrimaryColor    string                `form:"primary_color"`
	SecondaryColor  string                `form:"secondary_color"`
	TertiaryColor   string                `form:"tertiary_color"`
	QuaternaryColor string                `form:"quaternary_color"`
	NavbarImageFile *multipart.FileHeader `form:"navbar_image_file"`
	CustomDomain    string                `form:"custom_domain"`
	IsActive        bool                  `form:"is_active"`

	// Configuración de funcionalidades
	EnableDelivery     bool `form:"enable_delivery"`
	EnablePickup       bool `form:"enable_pickup"`
	EnableReservations bool `form:"enable_reservations"`
}
