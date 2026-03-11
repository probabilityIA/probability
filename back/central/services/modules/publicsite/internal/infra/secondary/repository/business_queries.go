package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) GetBusinessBySlug(ctx context.Context, slug string) (*entities.BusinessPage, error) {
	var business models.Business
	err := r.db.Conn(ctx).
		Where("code = ? AND is_active = true", slug).
		First(&business).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	page := &entities.BusinessPage{
		ID:              business.ID,
		Name:            business.Name,
		Code:            business.Code,
		Description:     business.Description,
		LogoURL:         business.LogoURL,
		PrimaryColor:    business.PrimaryColor,
		SecondaryColor:  business.SecondaryColor,
		TertiaryColor:   business.TertiaryColor,
		QuaternaryColor: business.QuaternaryColor,
		NavbarImageURL:  business.NavbarImageURL,
	}

	// Load website config if exists
	var config models.BusinessWebsiteConfig
	err = r.db.Conn(ctx).Where("business_id = ?", business.ID).First(&config).Error
	if err == nil {
		template := config.Template
		if template == "" {
			template = "default"
		}
		page.WebsiteConfig = &entities.WebsiteConfig{
			Template:             template,
			ShowHero:             config.ShowHero,
			ShowAbout:            config.ShowAbout,
			ShowFeaturedProducts: config.ShowFeaturedProducts,
			ShowFullCatalog:      config.ShowFullCatalog,
			ShowTestimonials:     config.ShowTestimonials,
			ShowLocation:         config.ShowLocation,
			ShowContact:          config.ShowContact,
			ShowSocialMedia:      config.ShowSocialMedia,
			ShowWhatsApp:         config.ShowWhatsApp,
			HeroContent:          config.HeroContent,
			AboutContent:         config.AboutContent,
			TestimonialsContent:  config.TestimonialsContent,
			LocationContent:      config.LocationContent,
			ContactContent:       config.ContactContent,
			SocialMediaContent:   config.SocialMediaContent,
			WhatsAppContent:      config.WhatsAppContent,
		}
	}
	// If no config exists, WebsiteConfig remains nil (defaults apply on frontend)

	return page, nil
}
