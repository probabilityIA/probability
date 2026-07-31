package entities

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

type WebsiteConfig struct {
	Template             string
	SectionsOrder        []byte
	ThemeContent         []byte
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
