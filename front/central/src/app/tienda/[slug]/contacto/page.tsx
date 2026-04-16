import { notFound } from 'next/navigation';
import { getPublicBusinessAction } from '@/services/modules/publicsite/infra/actions';
import { getTemplate } from '@/services/modules/publicsite/ui/templates/registry';

interface PageProps {
    params: Promise<{ slug: string }>;
}

export default async function ContactoPage({ params }: PageProps) {
    const { slug } = await params;
    const business = await getPublicBusinessAction(slug);
    if (!business) return notFound();

    const template = getTemplate(business.website_config?.template || 'default');
    const Contact = template.ContactSection;

    return (
        <div className="py-8 px-4 max-w-3xl mx-auto">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-8 text-center">Contáctanos</h1>
            <Contact slug={slug} content={business.website_config?.contact_content || null} />
        </div>
    );
}
