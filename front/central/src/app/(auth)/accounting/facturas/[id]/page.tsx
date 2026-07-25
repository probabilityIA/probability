import { redirect } from 'next/navigation';

interface FacturaDetallePageProps {
    params: Promise<{ id: string }>;
    searchParams: Promise<{ editar?: string }>;
}

export default async function FacturaDetallePage({ params, searchParams }: FacturaDetallePageProps) {
    const [{ id }, query] = await Promise.all([params, searchParams]);
    const invoiceId = Number(id) || 0;
    if (invoiceId <= 0) {
        redirect('/accounting/facturas');
    }
    redirect(`/accounting/facturas?ver=${invoiceId}${query.editar === '1' ? '&editar=1' : ''}`);
}
