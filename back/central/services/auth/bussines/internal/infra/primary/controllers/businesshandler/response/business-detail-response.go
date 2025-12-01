package response

import "time"

// BusinessTypeDetailResponse represents a business type in detail responses
type BusinessTypeDetailResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BusinessDetailResponse represents full business information for editing
type BusinessDetailResponse struct {
	ID                 uint                       `json:"id"`
	Name               string                     `json:"name"`
	Code               string                     `json:"code"`
	BusinessType       BusinessTypeDetailResponse `json:"business_type"`
	Timezone           string                     `json:"timezone"`
	Address            string                     `json:"address"`
	Description        string                     `json:"description"`
	LogoURL            string                     `json:"logo_url"`
	PrimaryColor       string                     `json:"primary_color"`
	SecondaryColor     string                     `json:"secondary_color"`
	TertiaryColor      string                     `json:"tertiary_color"`
	QuaternaryColor    string                     `json:"quaternary_color"`
	NavbarImageURL     string                     `json:"navbar_image_url"`
	CustomDomain       string                     `json:"custom_domain"`
	IsActive           bool                       `json:"is_active"`
	EnableDelivery     bool                       `json:"enable_delivery"`
	EnablePickup       bool                       `json:"enable_pickup"`
	EnableReservations bool                       `json:"enable_reservations"`
	CreatedAt          time.Time                  `json:"created_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
}

// GetBusinessByIDResponse represents the response for getting a business by ID with full information
type GetBusinessByIDResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    BusinessDetailResponse `json:"data"`
}
