import { CuentaForm } from '@/services/modules/publicsite/ui/session/CuentaForm';

interface PageProps {
    params: Promise<{ slug: string }>;
}

export default async function CuentaPage({ params }: PageProps) {
    const { slug } = await params;

    return (
        <div className="py-12 px-4">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-8 text-center">Tu cuenta</h1>
            <CuentaForm slug={slug} />
        </div>
    );
}
