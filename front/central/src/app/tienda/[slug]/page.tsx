import { notFound } from 'next/navigation';
import { getPublicBusinessAction } from '@/services/modules/publicsite/infra/actions';
import { getTemplate } from '@/services/modules/publicsite/ui/templates/registry';
import { SectionRenderer, isSectionVisible, resolveSectionsOrder } from '@/services/modules/publicsite/ui/components/SectionRenderer';
import { DEFAULT_SECTIONS_ORDER } from '@/services/modules/publicsite/domain/types';

export const revalidate = 60;

interface PageProps {
    params: Promise<{ slug: string }>;
}

export default async function TiendaPage({ params }: PageProps) {
    const { slug } = await params;
    const business = await getPublicBusinessAction(slug);
    if (!business) return notFound();

    const config = business.website_config;
    const template = getTemplate(config?.template || 'default');

    if (template.HomePage) {
        const HomePage = template.HomePage;
        return <HomePage business={business} slug={slug} config={config!} />;
    }

    const order = resolveSectionsOrder(config, DEFAULT_SECTIONS_ORDER);

    return (
        <>
            {order.filter(key => isSectionVisible(key, config)).map(key => (
                <SectionRenderer
                    key={key}
                    sectionKey={key}
                    business={business}
                    slug={slug}
                    config={config}
                    template={template}
                />
            ))}
        </>
    );
}
