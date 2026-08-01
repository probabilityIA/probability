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

export type SectionKey =
    | 'hero'
    | 'about'
    | 'featured_products'
    | 'testimonials'
    | 'location'
    | 'social_media'
    | 'contact'
    | 'full_catalog';

export const DEFAULT_SECTIONS_ORDER: SectionKey[] = [
    'hero',
    'about',
    'featured_products',
    'testimonials',
    'location',
    'social_media',
    'contact',
    'full_catalog',
];

export interface ThemeContent {
    primary?: string;
    secondary?: string;
    tertiary?: string;
    quaternary?: string;
    background?: string;
    apply_background?: boolean;
    apply_navbar?: boolean;
    background_image?: string;
    background_overlay?: number;
}

export interface NavbarLink {
    label: string;
    href: string;
    icon?: string;
    icon_type?: 'symbol' | 'image';
    icon_image?: string;
}

export interface NavbarContent {
    position?: 'top' | 'left' | 'right';
    alignment?: 'left' | 'center' | 'right';
    collapsible?: boolean;
    links?: NavbarLink[];
    announcement_text?: string;
    announcement_color?: string;
    announcement_animation?: 'none' | 'fade' | 'slide' | 'marquee';
}

export interface WebsiteConfig {
    template: string;
    sections_order?: SectionKey[] | null;
    theme_content?: ThemeContent | null;
    navbar_content?: NavbarContent | null;
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

export interface HeroPosition {
    x: number;
    y: number;
}

export interface HeroBlock {
    type: 'text' | 'image' | 'button' | 'empty';
    title?: string;
    subtitle?: string;
    title_size?: number;
    subtitle_size?: number;
    image?: string;
    button_text?: string;
    button_link?: string;
    position?: HeroPosition;
}

export interface HeroContent {
    title?: string;
    subtitle?: string;
    title_size?: number;
    subtitle_size?: number;
    cta_text?: string;
    background_image?: string;
    layout_rows?: number;
    layout_columns?: number;
    column_widths?: number[];
    blocks?: HeroBlock[];
    position?: HeroPosition;
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
    show_map?: boolean;
    show_directions?: boolean;
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

export interface TiendaSession {
    name: string;
    avatar_url?: string;
    email: string;
    phone: string;
    dni: string;
    customer_id: number;
    address: {
        street: string;
        city: string;
        state: string;
        country: string;
        postal_code: string;
    } | null;
}

export interface TiendaOrderSummary {
    id: string;
    order_number: string;
    total: number;
    currency: string;
    status: string;
    created_at: string;
    items_count: number;
}

export interface TiendaAddress {
    id: number;
    street: string;
    city: string;
    state: string;
    country: string;
    postal_code: string;
    is_primary: boolean;
}

export interface TiendaAddressInput {
    street: string;
    city?: string;
    state?: string;
    country?: string;
    postal_code?: string;
    is_primary?: boolean;
}

export interface TiendaProfileUpdate {
    name?: string;
    phone?: string;
    dni?: string;
}

export interface TiendaRegisterDTO {
    name: string;
    email: string;
    password: string;
    phone?: string;
    dni?: string;
}

export interface ContactFormDTO {
    name: string;
    email?: string;
    phone?: string;
    message: string;
}

export interface CartItem {
    product_id: string;
    name: string;
    sku: string;
    price: number;
    image_url: string;
    quantity: number;
}

export interface CheckoutAddressInput {
    street: string;
    city: string;
    state: string;
    country: string;
    postal_code?: string;
    instructions?: string;
}

export interface CreateCheckoutInput {
    items: { product_id: string; quantity: number }[];
    customer_name: string;
    customer_email?: string;
    customer_phone?: string;
    customer_dni?: string;
    payment_method?: string;
    address?: CheckoutAddressInput;
}

export interface CheckoutSession {
    payment_method?: string;
    reference: string;
    amount: number;
    currency: string;
    hash: string;
    public_key: string;
    redirection_url: string;
    is_sandbox: boolean;
    polling_enabled: boolean;
}
