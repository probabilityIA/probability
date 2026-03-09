package request

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/dtos"
)

type UpdateConfigRequest struct {
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

func (r *UpdateConfigRequest) ToDTO() *dtos.UpdateConfigDTO {
	return &dtos.UpdateConfigDTO{
		Template:             r.Template,
		ShowHero:             r.ShowHero,
		ShowAbout:            r.ShowAbout,
		ShowFeaturedProducts: r.ShowFeaturedProducts,
		ShowFullCatalog:      r.ShowFullCatalog,
		ShowTestimonials:     r.ShowTestimonials,
		ShowLocation:         r.ShowLocation,
		ShowContact:          r.ShowContact,
		ShowSocialMedia:      r.ShowSocialMedia,
		ShowWhatsApp:         r.ShowWhatsApp,
		HeroContent:          r.HeroContent,
		AboutContent:         r.AboutContent,
		TestimonialsContent:  r.TestimonialsContent,
		LocationContent:      r.LocationContent,
		ContactContent:       r.ContactContent,
		SocialMediaContent:   r.SocialMediaContent,
		WhatsAppContent:      r.WhatsAppContent,
	}
}
