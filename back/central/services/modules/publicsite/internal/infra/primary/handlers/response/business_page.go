package response

import (
	"encoding/json"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
)

type WebsiteConfigResponse struct {
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

type BusinessPageResponse struct {
	ID               uint                   `json:"id"`
	Name             string                 `json:"name"`
	Code             string                 `json:"code"`
	Description      string                 `json:"description"`
	LogoURL          string                 `json:"logo_url"`
	PrimaryColor     string                 `json:"primary_color"`
	SecondaryColor   string                 `json:"secondary_color"`
	TertiaryColor    string                 `json:"tertiary_color"`
	QuaternaryColor  string                 `json:"quaternary_color"`
	NavbarImageURL   string                 `json:"navbar_image_url"`
	WebsiteConfig    *WebsiteConfigResponse `json:"website_config"`
	FeaturedProducts []ProductResponse      `json:"featured_products"`
}

func buildFullImageURL(relativePath string, imageURLBase string) string {
	if relativePath == "" || imageURLBase == "" {
		return relativePath
	}
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}
	return imageURLBase + relativePath
}

func BusinessPageFromEntity(b *entities.BusinessPage, featured []entities.PublicProduct, imageURLBase string) BusinessPageResponse {
	resp := BusinessPageResponse{
		ID:              b.ID,
		Name:            b.Name,
		Code:            b.Code,
		Description:     b.Description,
		LogoURL:         buildFullImageURL(b.LogoURL, imageURLBase),
		PrimaryColor:    b.PrimaryColor,
		SecondaryColor:  b.SecondaryColor,
		TertiaryColor:   b.TertiaryColor,
		QuaternaryColor: b.QuaternaryColor,
		NavbarImageURL:  buildFullImageURL(b.NavbarImageURL, imageURLBase),
	}

	if b.WebsiteConfig != nil {
		wc := &WebsiteConfigResponse{
			Template:             b.WebsiteConfig.Template,
			ShowHero:             b.WebsiteConfig.ShowHero,
			ShowAbout:            b.WebsiteConfig.ShowAbout,
			ShowFeaturedProducts: b.WebsiteConfig.ShowFeaturedProducts,
			ShowFullCatalog:      b.WebsiteConfig.ShowFullCatalog,
			ShowTestimonials:     b.WebsiteConfig.ShowTestimonials,
			ShowLocation:         b.WebsiteConfig.ShowLocation,
			ShowContact:          b.WebsiteConfig.ShowContact,
			ShowSocialMedia:      b.WebsiteConfig.ShowSocialMedia,
			ShowWhatsApp:         b.WebsiteConfig.ShowWhatsApp,
			HeroContent:          b.WebsiteConfig.HeroContent,
			AboutContent:         b.WebsiteConfig.AboutContent,
			TestimonialsContent:  b.WebsiteConfig.TestimonialsContent,
			LocationContent:      b.WebsiteConfig.LocationContent,
			ContactContent:       b.WebsiteConfig.ContactContent,
			SocialMediaContent:   b.WebsiteConfig.SocialMediaContent,
			WhatsAppContent:      b.WebsiteConfig.WhatsAppContent,
		}
		resp.WebsiteConfig = wc
	}

	if featured != nil {
		resp.FeaturedProducts = ProductsFromEntities(featured, imageURLBase)
	}

	return resp
}
