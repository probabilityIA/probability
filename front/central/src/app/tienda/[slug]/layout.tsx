import { notFound } from 'next/navigation';
import { getPublicBusinessAction } from '@/services/modules/publicsite/infra/actions';
import { getTemplate } from '@/services/modules/publicsite/ui/templates/registry';
import { CartProvider } from '@/services/modules/publicsite/ui/cart/cart-context';
import { CartWidget } from '@/services/modules/publicsite/ui/cart/CartWidget';

interface LayoutProps {
    children: React.ReactNode;
    params: Promise<{ slug: string }>;
}

export default async function TiendaLayout({ children, params }: LayoutProps) {
    const { slug } = await params;
    const business = await getPublicBusinessAction(slug);
    if (!business) return notFound();

    const config = business.website_config;
    const template = getTemplate(config?.template || 'default');

    const Nav = template.Nav;
    const Footer = template.Footer;
    const WhatsApp = template.WhatsAppButton;

    return (
        <div
            style={{
                '--brand-primary': business.primary_color || '#1f2937',
                '--brand-secondary': business.secondary_color || '#3b82f6',
                '--brand-tertiary': business.tertiary_color || '#10b981',
                '--brand-quaternary': business.quaternary_color || '#fbbf24',
            } as React.CSSProperties}
        >
            <CartProvider slug={slug}>
                <Nav business={business} />
                <main className="min-h-screen">
                    {children}
                </main>
                <Footer business={business} />
                {config?.show_whatsapp && config.whatsapp_content && (
                    <WhatsApp content={config.whatsapp_content} />
                )}
                <CartWidget slug={slug} />
            </CartProvider>
        </div>
    );
}
