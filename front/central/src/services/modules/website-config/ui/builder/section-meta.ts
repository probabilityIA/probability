import { SectionKey, DEFAULT_SECTIONS_ORDER } from '@/services/modules/publicsite/domain/types';

export interface SectionMeta {
    key: SectionKey;
    label: string;
    toggleField: string;
    contentField: string | null;
}

export const SECTION_META: Record<SectionKey, SectionMeta> = {
    hero: { key: 'hero', label: 'Hero / Banner', toggleField: 'show_hero', contentField: 'hero_content' },
    about: { key: 'about', label: 'Sobre Nosotros', toggleField: 'show_about', contentField: 'about_content' },
    featured_products: { key: 'featured_products', label: 'Productos Destacados', toggleField: 'show_featured_products', contentField: null },
    testimonials: { key: 'testimonials', label: 'Testimonios', toggleField: 'show_testimonials', contentField: 'testimonials_content' },
    location: { key: 'location', label: 'Ubicacion', toggleField: 'show_location', contentField: 'location_content' },
    social_media: { key: 'social_media', label: 'Redes Sociales', toggleField: 'show_social_media', contentField: 'social_media_content' },
    contact: { key: 'contact', label: 'Contacto', toggleField: 'show_contact', contentField: 'contact_content' },
    full_catalog: { key: 'full_catalog', label: 'CTA Catalogo', toggleField: 'show_full_catalog', contentField: null },
};

export function normalizeOrder(stored: string[] | null | undefined): SectionKey[] {
    if (!stored || stored.length === 0) return [...DEFAULT_SECTIONS_ORDER];
    const valid = stored.filter((k): k is SectionKey => DEFAULT_SECTIONS_ORDER.includes(k as SectionKey));
    const missing = DEFAULT_SECTIONS_ORDER.filter(k => !valid.includes(k));
    return [...valid, ...missing];
}
