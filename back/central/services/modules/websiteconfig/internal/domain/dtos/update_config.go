package dtos

import "encoding/json"

// UpdateConfigDTO for updating website configuration
type UpdateConfigDTO struct {
	Template             *string         `json:"template"`
	ShowHero             *bool           `json:"show_hero"`
	ShowAbout            *bool           `json:"show_about"`
	ShowFeaturedProducts *bool           `json:"show_featured_products"`
	ShowFullCatalog      *bool           `json:"show_full_catalog"`
	ShowTestimonials     *bool           `json:"show_testimonials"`
	ShowLocation         *bool           `json:"show_location"`
	ShowContact          *bool           `json:"show_contact"`
	ShowSocialMedia      *bool           `json:"show_social_media"`
	ShowWhatsApp         *bool           `json:"show_whatsapp"`
	HeroContent          json.RawMessage `json:"hero_content"`
	AboutContent         json.RawMessage `json:"about_content"`
	TestimonialsContent  json.RawMessage `json:"testimonials_content"`
	LocationContent      json.RawMessage `json:"location_content"`
	ContactContent       json.RawMessage `json:"contact_content"`
	SocialMediaContent   json.RawMessage `json:"social_media_content"`
	WhatsAppContent      json.RawMessage `json:"whatsapp_content"`
}
