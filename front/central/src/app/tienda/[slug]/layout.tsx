import { notFound } from 'next/navigation';
import { getPublicBusinessAction } from '@/services/modules/publicsite/infra/actions';
import { getTemplate } from '@/services/modules/publicsite/ui/templates/registry';
import { CartProvider } from '@/services/modules/publicsite/ui/cart/cart-context';
import { CartWidget } from '@/services/modules/publicsite/ui/cart/CartWidget';
import { buildThemeStyle } from '@/services/modules/publicsite/ui/theme';

export const revalidate = 60;

export function generateStaticParams(): Array<{ slug: string }> {
    return [];
}

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
    const theme = config?.theme_content;

    const Nav = template.Nav;
    const Footer = template.Footer;
    const WhatsApp = template.WhatsAppButton;

    const themeStyle = buildThemeStyle(theme, {
        primary: business.primary_color || '#1f2937',
        secondary: business.secondary_color || '#3b82f6',
        tertiary: business.tertiary_color || '#10b981',
        quaternary: business.quaternary_color || '#fbbf24',
    });

    return (
        <div style={themeStyle}>
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
