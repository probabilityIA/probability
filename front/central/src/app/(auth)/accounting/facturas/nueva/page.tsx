import { redirect } from 'next/navigation';

export default function NuevaFacturaPage() {
    redirect('/accounting/facturas?nueva=1');
}
