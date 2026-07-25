import { AccountingGate, ConfigView } from '@/services/modules/accounting/ui';
import {
    getAccountingConceptsAction,
    getAccountingDianConfigAction,
    getAccountingTaxesAction,
} from '@/services/modules/accounting/infra/actions';

export default async function ConfiguracionPage() {
    const [conceptsResult, taxesResult, dianResult] = await Promise.all([
        getAccountingConceptsAction(),
        getAccountingTaxesAction(),
        getAccountingDianConfigAction(),
    ]);

    const concepts = conceptsResult.success ? conceptsResult.data : [];
    const taxes = taxesResult.success ? taxesResult.data : [];
    const dianConfig = dianResult.success ? dianResult.data : null;
    const error = !conceptsResult.success
        ? conceptsResult.error
        : !taxesResult.success
            ? taxesResult.error
            : null;

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <AccountingGate>
                <ConfigView concepts={concepts} taxes={taxes} dianConfig={dianConfig} error={error} />
            </AccountingGate>
        </div>
    );
}
