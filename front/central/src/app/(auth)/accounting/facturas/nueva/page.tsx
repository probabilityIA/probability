import { AccountingGate, InvoiceForm } from '@/services/modules/accounting/ui';
import {
    getAccountingConceptsAction,
    getAccountingTaxesAction,
} from '@/services/modules/accounting/infra/actions';

export default async function NuevaFacturaPage() {
    const [conceptsResult, taxesResult] = await Promise.all([
        getAccountingConceptsAction(),
        getAccountingTaxesAction(),
    ]);

    const concepts = conceptsResult.success ? conceptsResult.data : [];
    const taxes = taxesResult.success ? taxesResult.data : [];
    const error = !conceptsResult.success
        ? conceptsResult.error
        : !taxesResult.success
            ? taxesResult.error
            : null;

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <AccountingGate>
                <InvoiceForm concepts={concepts} taxes={taxes} error={error} />
            </AccountingGate>
        </div>
    );
}
