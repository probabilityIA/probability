import { PublicBusiness, HeroContent, HeroPosition, NavbarContent, AboutContent, Testimonial, LocationContent, ContactContent, SocialMediaContent, WhatsAppContent, PublicProduct, WebsiteConfig } from '../../domain/types';

export interface TemplateComponents {
    Nav: React.ComponentType<{ business: PublicBusiness; navbarContent?: NavbarContent | null }>;
    Footer: React.ComponentType<{ business: PublicBusiness }>;
    WhatsAppButton: React.ComponentType<{ content: WhatsAppContent }>;

    HeroSection: React.ComponentType<{
        content: HeroContent | null;
        business: PublicBusiness;
        slug: string;
        editable?: boolean;
        onContentPositionChange?: (position: HeroPosition) => void;
        onBlockPositionChange?: (index: number, position: HeroPosition) => void;
    }>;
    AboutSection: React.ComponentType<{ content: AboutContent }>;
    FeaturedProducts: React.ComponentType<{ products: PublicProduct[]; slug: string }>;
    TestimonialsSection: React.ComponentType<{ content: Testimonial[] }>;
    LocationSection: React.ComponentType<{ content: LocationContent }>;
    ContactSection: React.ComponentType<{ slug: string; content: ContactContent | null }>;
    SocialMediaLinks: React.ComponentType<{ content: SocialMediaContent }>;
    ProductCard: React.ComponentType<{ product: PublicProduct; slug: string }>;

    HomePage?: React.ComponentType<{ business: PublicBusiness; slug: string; config: WebsiteConfig }>;
}
