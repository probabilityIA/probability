import { AccountingGate, InvoicesView } from '@/services/modules/accounting/ui';
import {
    getAccountingConceptsAction,
    getAccountingInvoiceAction,
    getAccountingInvoicesAction,
    getAccountingServicesAction,
    getAccountingTaxesAction,
} from '@/services/modules/accounting/infra/actions';
import { InvoiceStatus } from '@/services/modules/accounting/domain/types';

const VALID_STATUSES: InvoiceStatus[] = ['DRAFT', 'SENT', 'PAID', 'CANCELLED'];

interface FacturasPageProps {
    searchParams: Promise<{
        page?: string;
        page_size?: string;
        status?: string;
        business_id?: string;
        is_test?: string;
        nueva?: string;
        ver?: string;
        editar?: string;
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
    const isTest = params.is_test === 'true' || params.is_test === 'false'
        ? (params.is_test as 'true' | 'false')
        : null;
    const openCreate = params.nueva === '1';
    const detailId = Number(params.ver) || 0;
    const openEditDetail = params.editar === '1';

    const [invoicesResult, conceptsResult, taxesResult, servicesResult, detailResult] = await Promise.all([
        getAccountingInvoicesAction({
            page,
            page_size: pageSize,
            status: status || undefined,
            business_id: businessId || undefined,
            is_test: isTest || undefined,
        }),
        getAccountingConceptsAction(),
        getAccountingTaxesAction(),
        getAccountingServicesAction(),
        detailId > 0 ? getAccountingInvoiceAction(detailId) : Promise.resolve(null),
    ]);

    const invoices = invoicesResult.success
        ? invoicesResult.data
        : { data: [], total: 0, page, page_size: pageSize, total_pages: 0 };
    const concepts = conceptsResult.success ? conceptsResult.data : [];
    const taxes = taxesResult.success ? taxesResult.data : [];
    const services = servicesResult.success ? servicesResult.data : [];
    const detailInvoice = detailResult?.success ? detailResult.data : null;
    const error = !invoicesResult.success
        ? invoicesResult.error
        : !conceptsResult.success
            ? conceptsResult.error
            : !taxesResult.success
                ? taxesResult.error
                : detailResult && !detailResult.success
                    ? detailResult.error
                    : null;

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <AccountingGate>
                <InvoicesView
                    invoices={invoices}
                    concepts={concepts}
                    taxes={taxes}
                    services={services}
                    filters={{ page, page_size: pageSize, status, business_id: businessId, is_test: isTest }}
                    detailInvoice={detailInvoice}
                    openEditDetail={openEditDetail}
                    error={error}
                    openCreate={openCreate}
                />
            </AccountingGate>
        </div>
    );
}
