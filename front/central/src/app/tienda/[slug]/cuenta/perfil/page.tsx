import { PerfilForm } from '@/services/modules/publicsite/ui/session/PerfilForm';

interface PageProps {
    params: Promise<{ slug: string }>;
}

export default async function PerfilPage({ params }: PageProps) {
    const { slug } = await params;

    return (
        <div className="py-12 px-4">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-8 text-center">Mi perfil</h1>
            <PerfilForm slug={slug} />
        </div>
    );
}
