package entities

// BusinessPage represents a business's public page info
type BusinessPage struct {
	ID              uint
	Name            string
	Code            string
	Description     string
	LogoURL         string
	PrimaryColor    string
	SecondaryColor  string
	TertiaryColor   string
	QuaternaryColor string
	NavbarImageURL  string
	WebsiteConfig   *WebsiteConfig
}

// WebsiteConfig holds section toggles and content
type WebsiteConfig struct {
	Template             string
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
