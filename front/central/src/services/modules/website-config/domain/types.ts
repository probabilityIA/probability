export interface WebsiteConfigData {
    id: number;
    business_id: number;
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
    hero_content: Record<string, any> | null;
    about_content: Record<string, any> | null;
    testimonials_content: Record<string, any>[] | null;
    location_content: Record<string, any> | null;
    contact_content: Record<string, any> | null;
    social_media_content: Record<string, any> | null;
    whatsapp_content: Record<string, any> | null;
}

export interface UpdateWebsiteConfigDTO {
    template?: string;
    show_hero?: boolean;
    show_about?: boolean;
    show_featured_products?: boolean;
    show_full_catalog?: boolean;
    show_testimonials?: boolean;
    show_location?: boolean;
    show_contact?: boolean;
    show_social_media?: boolean;
    show_whatsapp?: boolean;
    hero_content?: Record<string, any>;
    about_content?: Record<string, any>;
    testimonials_content?: Record<string, any>[];
    location_content?: Record<string, any>;
    contact_content?: Record<string, any>;
    social_media_content?: Record<string, any>;
    whatsapp_content?: Record<string, any>;
}
