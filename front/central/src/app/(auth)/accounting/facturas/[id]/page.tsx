import Link from 'next/link';
import { AccountingGate, InvoiceDetailView } from '@/services/modules/accounting/ui';
import { getAccountingInvoiceAction } from '@/services/modules/accounting/infra/actions';

interface FacturaDetallePageProps {
    params: Promise<{ id: string }>;
}

export default async function FacturaDetallePage({ params }: FacturaDetallePageProps) {
    const { id } = await params;
    const invoiceId = Number(id) || 0;
    const result = invoiceId > 0
        ? await getAccountingInvoiceAction(invoiceId)
        : { success: false as const, error: 'Factura no encontrada' };

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <AccountingGate>
                {result.success ? (
                    <InvoiceDetailView invoice={result.data} />
                ) : (
                    <div className="max-w-4xl space-y-4">
                        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 text-sm rounded-lg px-4 py-3">
                            {result.error}
                        </div>
                        <Link
                            href="/accounting/facturas"
                            className="inline-block px-4 py-2 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                        >
                            Volver a facturas
                        </Link>
                    </div>
                )}
            </AccountingGate>
        </div>
    );
}
