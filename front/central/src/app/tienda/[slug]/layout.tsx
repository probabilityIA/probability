import { notFound } from 'next/navigation';
import { getPublicBusinessAction } from '@/services/modules/publicsite/infra/actions';
import { getTemplate } from '@/services/modules/publicsite/ui/templates/registry';
import { CartProvider } from '@/services/modules/publicsite/ui/cart/cart-context';
import { TiendaPricesProvider } from '@/services/modules/publicsite/ui/session/prices-context';
import { CartWidget } from '@/services/modules/publicsite/ui/cart/CartWidget';
import { buildThemeStyle } from '@/services/modules/publicsite/ui/theme';
import { AnnouncementBar } from '@/services/modules/publicsite/ui/components/PublicNav';

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
    const navbarContent = config?.navbar_content;
    const isSidebar = navbarContent?.position === 'left' || navbarContent?.position === 'right';

    const themeStyle = buildThemeStyle(theme, {
        primary: business.primary_color || '#1f2937',
        secondary: business.secondary_color || '#3b82f6',
        tertiary: business.tertiary_color || '#10b981',
        quaternary: business.quaternary_color || '#fbbf24',
    });

    return (
        <div style={themeStyle}>
            <AnnouncementBar navbarContent={navbarContent} />
            <CartProvider slug={slug}>
                <TiendaPricesProvider slug={slug}>
                {isSidebar ? (
                    <div className="flex min-h-screen">
                        {navbarContent?.position === 'left' && <Nav business={business} navbarContent={navbarContent} />}
                        <div className="flex-1 flex flex-col">
                            <main className="flex-1">
                                {children}
                            </main>
                            <Footer business={business} />
                        </div>
                        {navbarContent?.position === 'right' && <Nav business={business} navbarContent={navbarContent} />}
                    </div>
                ) : (
                    <>
                        <Nav business={business} navbarContent={navbarContent} />
                        <main className="min-h-screen">
                            {children}
                        </main>
                        <Footer business={business} />
                    </>
                )}
                {config?.show_whatsapp && config.whatsapp_content && (
                    <WhatsApp content={config.whatsapp_content} />
                )}
                <CartWidget slug={slug} />
                </TiendaPricesProvider>
            </CartProvider>
        </div>
    );
}
