package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BusinessWebsiteConfig stores the public website configuration for a business
type BusinessWebsiteConfig struct {
	gorm.Model
	BusinessID uint     `gorm:"not null;uniqueIndex"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Template selection
	Template string `gorm:"type:varchar(50);default:'default'"`

	// Section toggles
	ShowHero             bool `gorm:"default:true"`
	ShowAbout            bool `gorm:"default:false"`
	ShowFeaturedProducts bool `gorm:"default:true"`
	ShowFullCatalog      bool `gorm:"default:true"`
	ShowTestimonials     bool `gorm:"default:false"`
	ShowLocation         bool `gorm:"default:false"`
	ShowContact          bool `gorm:"default:true"`
	ShowSocialMedia      bool `gorm:"default:false"`
	ShowWhatsApp         bool `gorm:"default:false"`

	// Section content (JSONB)
	HeroContent         datatypes.JSON `gorm:"type:jsonb"` // {title, subtitle, cta_text, background_image}
	AboutContent        datatypes.JSON `gorm:"type:jsonb"` // {text, image, mission, vision}
	TestimonialsContent datatypes.JSON `gorm:"type:jsonb"` // [{name, text, rating, avatar}]
	LocationContent     datatypes.JSON `gorm:"type:jsonb"` // {lat, lng, address, hours}
	ContactContent      datatypes.JSON `gorm:"type:jsonb"` // {email, phone, form_enabled, contacts:[{name,role,phone}]}
	SocialMediaContent  datatypes.JSON `gorm:"type:jsonb"` // {facebook, instagram, twitter, tiktok}
	WhatsAppContent     datatypes.JSON `gorm:"type:jsonb"` // {number, message, show_floating_button}
}
