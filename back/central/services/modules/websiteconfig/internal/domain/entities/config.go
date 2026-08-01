package entities

type WebsiteConfig struct {
	ID         uint
	BusinessID uint
	Template   string

	SectionsOrder    []byte
	ThemeContent     []byte
	HiddenCategories []byte
	NavbarContent    []byte

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

type ProductCategory struct {
	Name         string
	ProductCount int64
}
