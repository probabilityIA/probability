import { AccountingGate, InvoicesView } from '@/services/modules/accounting/ui';
import { getAccountingInvoicesAction } from '@/services/modules/accounting/infra/actions';
import { InvoiceStatus } from '@/services/modules/accounting/domain/types';

const VALID_STATUSES: InvoiceStatus[] = ['DRAFT', 'SENT', 'PAID', 'CANCELLED'];

interface FacturasPageProps {
    searchParams: Promise<{
        page?: string;
        page_size?: string;
        status?: string;
        business_id?: string;
    }>;
}

export default async function FacturasPage({ searchParams }: FacturasPageProps) {
    const params = await searchParams;
    const page = Math.max(Number(params.page) || 1, 1);
    const pageSize = Math.max(Number(params.page_size) || 10, 1);
    const status = VALID_STATUSES.includes(params.status as InvoiceStatus)
        ? (params.status as InvoiceStatus)
        : null;
    const businessId = Number(params.business_id) || null;

    const invoicesResult = await getAccountingInvoicesAction({
        page,
        page_size: pageSize,
        status: status || undefined,
        business_id: businessId || undefined,
    });

    const invoices = invoicesResult.success
        ? invoicesResult.data
        : { data: [], total: 0, page, page_size: pageSize, total_pages: 0 };
    const error = invoicesResult.success ? null : invoicesResult.error;

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <AccountingGate>
                <InvoicesView
                    invoices={invoices}
                    filters={{ page, page_size: pageSize, status, business_id: businessId }}
                    error={error}
                />
            </AccountingGate>
        </div>
    );
}
