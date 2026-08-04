package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BusinessWebsiteConfig struct {
	gorm.Model
	BusinessID uint     `gorm:"not null;uniqueIndex"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Template string `gorm:"type:varchar(50);default:'default'"`

	SectionsOrder    datatypes.JSON `gorm:"type:jsonb"`
	ThemeContent     datatypes.JSON `gorm:"type:jsonb"`
	HiddenCategories datatypes.JSON `gorm:"type:jsonb"`
	NavbarContent    datatypes.JSON `gorm:"type:jsonb"`

	ShowHero             bool `gorm:"default:true"`
	ShowAbout            bool `gorm:"default:false"`
	ShowFeaturedProducts bool `gorm:"default:true"`
	ShowFullCatalog      bool `gorm:"default:true"`
	ShowTestimonials     bool `gorm:"default:false"`
	ShowLocation         bool `gorm:"default:false"`
	ShowContact          bool `gorm:"default:true"`
	ShowSocialMedia      bool `gorm:"default:false"`
	ShowWhatsApp         bool `gorm:"default:false"`

	HeroContent         datatypes.JSON `gorm:"type:jsonb"`
	AboutContent        datatypes.JSON `gorm:"type:jsonb"`
	TestimonialsContent datatypes.JSON `gorm:"type:jsonb"`
	LocationContent     datatypes.JSON `gorm:"type:jsonb"`
	ContactContent      datatypes.JSON `gorm:"type:jsonb"`
	SocialMediaContent  datatypes.JSON `gorm:"type:jsonb"`
	WhatsAppContent     datatypes.JSON `gorm:"type:jsonb"`
}
