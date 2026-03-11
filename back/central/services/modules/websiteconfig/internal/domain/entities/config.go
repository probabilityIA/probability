package entities

// WebsiteConfig represents the website configuration for a business
type WebsiteConfig struct {
	ID         uint
	BusinessID uint
	Template   string

	ShowHero             bool
	ShowAbout            bool
	ShowFeaturedProducts bool
	ShowFullCatalog      bool
	ShowTestimonials     bool
	ShowLocation         bool
	ShowContact          bool
	ShowSocialMedia      bool
	ShowWhatsApp         bool

	HeroContent         []byte
	AboutContent        []byte
	TestimonialsContent []byte
	LocationContent     []byte
	ContactContent      []byte
	SocialMediaContent  []byte
	WhatsAppContent     []byte
}
