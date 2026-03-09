import { notFound } from 'next/navigation';
import { getPublicBusinessAction } from '@/services/modules/publicsite/infra/actions';
import { getTemplate } from '@/services/modules/publicsite/ui/templates/registry';
import Link from 'next/link';

interface PageProps {
    params: Promise<{ slug: string }>;
}

export default async function TiendaPage({ params }: PageProps) {
    const { slug } = await params;
    const business = await getPublicBusinessAction(slug);
    if (!business) return notFound();

    const config = business.website_config;
    const template = getTemplate(config?.template || 'default');

    // If the template provides a full HomePage component, use it
    if (template.HomePage) {
        const HomePage = template.HomePage;
        return <HomePage business={business} slug={slug} config={config!} />;
    }

    // Default: section-by-section rendering
    const defaults = {
        show_hero: true,
        show_featured_products: true,
        show_full_catalog: true,
        show_contact: true,
    };

    const show = (key: string) => config ? (config as any)[key] : (defaults as any)[key];

    return (
        <>
            {show('show_hero') && (
                <template.HeroSection
                    content={config?.hero_content || null}
                    business={business}
                    slug={slug}
                />
            )}

            {show('show_about') && config?.about_content && (
                <template.AboutSection content={config.about_content} />
            )}

            {show('show_featured_products') && business.featured_products?.length > 0 && (
                <section className="py-16 px-4 max-w-7xl mx-auto">
                    <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">Productos Destacados</h2>
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
            )}

            {show('show_testimonials') && config?.testimonials_content && (
                <template.TestimonialsSection content={config.testimonials_content} />
            )}

            {show('show_location') && config?.location_content && (
                <template.LocationSection content={config.location_content} />
            )}

            {show('show_social_media') && config?.social_media_content && (
                <template.SocialMediaLinks content={config.social_media_content} />
            )}

            {show('show_contact') && (
                <template.ContactSection slug={slug} content={config?.contact_content || null} />
            )}

            {/* CTA Section */}
            {show('show_full_catalog') && (
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
            )}
        </>
    );
}
