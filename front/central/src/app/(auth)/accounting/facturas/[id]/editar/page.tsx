import { redirect } from 'next/navigation';

interface EditarFacturaPageProps {
    params: Promise<{ id: string }>;
}

export default async function EditarFacturaPage({ params }: EditarFacturaPageProps) {
    const { id } = await params;
    const invoiceId = Number(id) || 0;
    if (invoiceId <= 0) {
        redirect('/accounting/facturas');
    }
    redirect(`/accounting/facturas/${invoiceId}?editar=1`);
}
