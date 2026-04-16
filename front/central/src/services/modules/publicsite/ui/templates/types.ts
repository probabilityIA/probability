import { PublicBusiness, HeroContent, AboutContent, Testimonial, LocationContent, ContactContent, SocialMediaContent, WhatsAppContent, PublicProduct, WebsiteConfig } from '../../domain/types';

/**
 * Contract that every template must fulfill.
 * Each component receives props using the existing domain types.
 */
export interface TemplateComponents {
    // Layout components
    Nav: React.ComponentType<{ business: PublicBusiness }>;
    Footer: React.ComponentType<{ business: PublicBusiness }>;
    WhatsAppButton: React.ComponentType<{ content: WhatsAppContent }>;

    // Homepage section components
    HeroSection: React.ComponentType<{ content: HeroContent | null; business: PublicBusiness; slug: string }>;
    AboutSection: React.ComponentType<{ content: AboutContent }>;
    FeaturedProducts: React.ComponentType<{ products: PublicProduct[]; slug: string }>;
    TestimonialsSection: React.ComponentType<{ content: Testimonial[] }>;
    LocationSection: React.ComponentType<{ content: LocationContent }>;
    ContactSection: React.ComponentType<{ slug: string; content: ContactContent | null }>;
    SocialMediaLinks: React.ComponentType<{ content: SocialMediaContent }>;
    ProductCard: React.ComponentType<{ product: PublicProduct; slug: string }>;

    // Optional: full control of the homepage layout (bypasses section-by-section rendering)
    HomePage?: React.ComponentType<{ business: PublicBusiness; slug: string; config: WebsiteConfig }>;
}
