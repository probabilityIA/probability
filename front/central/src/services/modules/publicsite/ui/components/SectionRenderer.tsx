import Link from 'next/link';
import { PublicBusiness, WebsiteConfig, SectionKey } from '../../domain/types';
import { TemplateComponents } from '../templates/types';

interface SectionRendererProps {
    sectionKey: SectionKey;
    business: PublicBusiness;
    slug: string;
    config: WebsiteConfig | null;
    template: TemplateComponents;
}

const DEFAULT_VISIBILITY: Record<string, boolean> = {
    show_hero: true,
    show_featured_products: true,
    show_full_catalog: true,
    show_contact: true,
};

export function isSectionVisible(sectionKey: SectionKey, config: WebsiteConfig | null): boolean {
    const flag = `show_${sectionKey === 'featured_products' ? 'featured_products' : sectionKey}`;
    const show = config ? (config as unknown as Record<string, boolean>)[flag] : DEFAULT_VISIBILITY[flag];
    if (!show) return false;

    switch (sectionKey) {
        case 'about':
            return !!config?.about_content;
        case 'testimonials':
            return !!config?.testimonials_content;
        case 'location':
            return !!config?.location_content;
        case 'social_media':
            return !!config?.social_media_content;
        default:
            return true;
    }
}

export function SectionRenderer({ sectionKey, business, slug, config, template }: SectionRendererProps) {
    switch (sectionKey) {
        case 'hero':
            return (
                <template.HeroSection
                    content={config?.hero_content || null}
                    business={business}
                    slug={slug}
                />
            );
        case 'about':
            return config?.about_content ? <template.AboutSection content={config.about_content} /> : null;
        case 'featured_products':
            if (!business.featured_products || business.featured_products.length === 0) return null;
            return (
                <section className="py-16 px-4 max-w-7xl mx-auto">
                    <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-8 text-center">Productos Destacados</h2>
                    <template.FeaturedProducts products={business.featured_products} slug={slug} />
                    <div className="text-center mt-8">
                        <Link
                            href={`/tienda/${slug}/productos`}
                            className="inline-block px-6 py-3 rounded-lg text-white font-medium transition-colors"
                            style={{ backgroundColor: 'var(--brand-secondary)' }}
                        >
                            Ver todo el catálogo
                        </Link>
                    </div>
                </section>
            );
        case 'testimonials':
            return config?.testimonials_content ? <template.TestimonialsSection content={config.testimonials_content} /> : null;
        case 'location':
            return config?.location_content ? <template.LocationSection content={config.location_content} /> : null;
        case 'social_media':
            return config?.social_media_content ? <template.SocialMediaLinks content={config.social_media_content} /> : null;
        case 'contact':
            return <template.ContactSection slug={slug} content={config?.contact_content || null} />;
        case 'full_catalog':
            return (
                <section className="py-16 text-center" style={{ backgroundColor: 'var(--brand-primary)' }}>
                    <h2 className="text-3xl font-bold text-white mb-4">¿Listo para hacer tu pedido?</h2>
                    <p className="text-white/80 mb-8 max-w-2xl mx-auto">
                        Explora nuestro catálogo completo y realiza tu pedido en línea
                    </p>
                    <Link
                        href={`/login?redirect=/storefront/catalogo&business_code=${slug}`}
                        className="inline-block px-8 py-4 bg-white rounded-lg font-bold text-lg transition-transform hover:scale-105"
                        style={{ color: 'var(--brand-primary)' }}
                    >
                        Hacer Pedido
                    </Link>
                </section>
            );
        default:
            return null;
    }
}

export function resolveSectionsOrder(config: WebsiteConfig | null, defaultOrder: SectionKey[]): SectionKey[] {
    const stored = config?.sections_order;
    if (!stored || !Array.isArray(stored) || stored.length === 0) return defaultOrder;
    const valid = stored.filter((k): k is SectionKey => defaultOrder.includes(k as SectionKey));
    const missing = defaultOrder.filter(k => !valid.includes(k));
    return [...valid, ...missing];
}
