package response

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
)

type ConfigResponse struct {
	ID                   uint            `json:"id"`
	BusinessID           uint            `json:"business_id"`
	Template             string          `json:"template"`
	ShowHero             bool            `json:"show_hero"`
	ShowAbout            bool            `json:"show_about"`
	ShowFeaturedProducts bool            `json:"show_featured_products"`
	ShowFullCatalog      bool            `json:"show_full_catalog"`
	ShowTestimonials     bool            `json:"show_testimonials"`
	ShowLocation         bool            `json:"show_location"`
	ShowContact          bool            `json:"show_contact"`
	ShowSocialMedia      bool            `json:"show_social_media"`
	ShowWhatsApp         bool            `json:"show_whatsapp"`
	HeroContent          json.RawMessage `json:"hero_content"`
	AboutContent         json.RawMessage `json:"about_content"`
	TestimonialsContent  json.RawMessage `json:"testimonials_content"`
	LocationContent      json.RawMessage `json:"location_content"`
	ContactContent       json.RawMessage `json:"contact_content"`
	SocialMediaContent   json.RawMessage `json:"social_media_content"`
	WhatsAppContent      json.RawMessage `json:"whatsapp_content"`
}

func ConfigFromEntity(e *entities.WebsiteConfig) ConfigResponse {
	return ConfigResponse{
		ID:                   e.ID,
		BusinessID:           e.BusinessID,
		Template:             e.Template,
		ShowHero:             e.ShowHero,
		ShowAbout:            e.ShowAbout,
		ShowFeaturedProducts: e.ShowFeaturedProducts,
		ShowFullCatalog:      e.ShowFullCatalog,
		ShowTestimonials:     e.ShowTestimonials,
		ShowLocation:         e.ShowLocation,
		ShowContact:          e.ShowContact,
		ShowSocialMedia:      e.ShowSocialMedia,
		ShowWhatsApp:         e.ShowWhatsApp,
		HeroContent:          e.HeroContent,
		AboutContent:         e.AboutContent,
		TestimonialsContent:  e.TestimonialsContent,
		LocationContent:      e.LocationContent,
		ContactContent:       e.ContactContent,
		SocialMediaContent:   e.SocialMediaContent,
		WhatsAppContent:      e.WhatsAppContent,
	}
}
