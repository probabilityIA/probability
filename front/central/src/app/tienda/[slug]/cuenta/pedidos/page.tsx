import { MisPedidos } from '@/services/modules/publicsite/ui/session/MisPedidos';

interface PageProps {
    params: Promise<{ slug: string }>;
}

export default async function MisPedidosPage({ params }: PageProps) {
    const { slug } = await params;

    return (
        <div className="py-12 px-4">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-8 text-center">Mis pedidos</h1>
            <MisPedidos slug={slug} />
        </div>
    );
}
