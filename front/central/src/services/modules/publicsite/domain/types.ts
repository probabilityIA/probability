export interface PublicBusiness {
    id: number;
    name: string;
    code: string;
    description: string;
    logo_url: string;
    primary_color: string;
    secondary_color: string;
    tertiary_color: string;
    quaternary_color: string;
    navbar_image_url: string;
    website_config: WebsiteConfig | null;
    featured_products: PublicProduct[];
}

export interface WebsiteConfig {
    template: string;
    show_hero: boolean;
    show_about: boolean;
    show_featured_products: boolean;
    show_full_catalog: boolean;
    show_testimonials: boolean;
    show_location: boolean;
    show_contact: boolean;
    show_social_media: boolean;
    show_whatsapp: boolean;
    hero_content: HeroContent | null;
    about_content: AboutContent | null;
    testimonials_content: Testimonial[] | null;
    location_content: LocationContent | null;
    contact_content: ContactContent | null;
    social_media_content: SocialMediaContent | null;
    whatsapp_content: WhatsAppContent | null;
}

export interface HeroContent {
    title?: string;
    subtitle?: string;
    cta_text?: string;
    background_image?: string;
}

export interface AboutContent {
    text?: string;
    image?: string;
    mission?: string;
    vision?: string;
}

export interface Testimonial {
    name: string;
    text: string;
    rating?: number;
    avatar?: string;
}

export interface LocationContent {
    lat?: number;
    lng?: number;
    address?: string;
    hours?: string;
}

export interface ContactContent {
    email?: string;
    phone?: string;
    form_enabled?: boolean;
    contacts?: { name: string; role: string; phone: string }[];
}

export interface SocialMediaContent {
    facebook?: string;
    instagram?: string;
    twitter?: string;
    tiktok?: string;
}

export interface WhatsAppContent {
    number?: string;
    message?: string;
    show_floating_button?: boolean;
}

export interface PublicProduct {
    id: string;
    name: string;
    description: string;
    short_description: string;
    price: number;
    compare_at_price?: number;
    currency: string;
    image_url: string;
    images?: string[];
    sku: string;
    stock_quantity: number;
    category: string;
    brand: string;
    is_featured: boolean;
    created_at: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface ContactFormDTO {
    name: string;
    email?: string;
    phone?: string;
    message: string;
}
