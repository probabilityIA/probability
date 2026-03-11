package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (r *Repository) GetConfig(ctx context.Context, businessID uint) (*entities.WebsiteConfig, error) {
	var config models.BusinessWebsiteConfig
	err := r.db.Conn(ctx).Where("business_id = ?", businessID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	template := config.Template
	if template == "" {
		template = "default"
	}

	return &entities.WebsiteConfig{
		ID:                   config.ID,
		BusinessID:           config.BusinessID,
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
	}, nil
}

func (r *Repository) UpsertConfig(ctx context.Context, businessID uint, dto *dtos.UpdateConfigDTO) (*entities.WebsiteConfig, error) {
	var config models.BusinessWebsiteConfig
	err := r.db.Conn(ctx).Where("business_id = ?", businessID).First(&config).Error

	if err == gorm.ErrRecordNotFound {
		// Create new config
		config = models.BusinessWebsiteConfig{
			BusinessID: businessID,
		}
	} else if err != nil {
		return nil, err
	}

	// Apply updates (only non-nil fields)
	if dto.Template != nil {
		config.Template = *dto.Template
	}
	if dto.ShowHero != nil {
		config.ShowHero = *dto.ShowHero
	}
	if dto.ShowAbout != nil {
		config.ShowAbout = *dto.ShowAbout
	}
	if dto.ShowFeaturedProducts != nil {
		config.ShowFeaturedProducts = *dto.ShowFeaturedProducts
	}
	if dto.ShowFullCatalog != nil {
		config.ShowFullCatalog = *dto.ShowFullCatalog
	}
	if dto.ShowTestimonials != nil {
		config.ShowTestimonials = *dto.ShowTestimonials
	}
	if dto.ShowLocation != nil {
		config.ShowLocation = *dto.ShowLocation
	}
	if dto.ShowContact != nil {
		config.ShowContact = *dto.ShowContact
	}
	if dto.ShowSocialMedia != nil {
		config.ShowSocialMedia = *dto.ShowSocialMedia
	}
	if dto.ShowWhatsApp != nil {
		config.ShowWhatsApp = *dto.ShowWhatsApp
	}

	if dto.HeroContent != nil {
		config.HeroContent = datatypes.JSON(dto.HeroContent)
	}
	if dto.AboutContent != nil {
		config.AboutContent = datatypes.JSON(dto.AboutContent)
	}
	if dto.TestimonialsContent != nil {
		config.TestimonialsContent = datatypes.JSON(dto.TestimonialsContent)
	}
	if dto.LocationContent != nil {
		config.LocationContent = datatypes.JSON(dto.LocationContent)
	}
	if dto.ContactContent != nil {
		config.ContactContent = datatypes.JSON(dto.ContactContent)
	}
	if dto.SocialMediaContent != nil {
		config.SocialMediaContent = datatypes.JSON(dto.SocialMediaContent)
	}
	if dto.WhatsAppContent != nil {
		config.WhatsAppContent = datatypes.JSON(dto.WhatsAppContent)
	}

	if config.ID == 0 {
		err = r.db.Conn(ctx).Create(&config).Error
	} else {
		err = r.db.Conn(ctx).Save(&config).Error
	}

	if err != nil {
		return nil, err
	}

	upsertTemplate := config.Template
	if upsertTemplate == "" {
		upsertTemplate = "default"
	}

	return &entities.WebsiteConfig{
		ID:                   config.ID,
		BusinessID:           config.BusinessID,
		Template:             upsertTemplate,
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
	}, nil
}
