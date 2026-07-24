import { AccountingGate, EntriesView } from '@/services/modules/accounting/ui';
import {
    getAccountingConceptsAction,
    getAccountingEntriesAction,
} from '@/services/modules/accounting/infra/actions';
import { ConceptKind } from '@/services/modules/accounting/domain/types';

interface MovimientosPageProps {
    searchParams: Promise<{
        page?: string;
        page_size?: string;
        from?: string;
        to?: string;
        concept_id?: string;
        kind?: string;
    }>;
}

export default async function MovimientosPage({ searchParams }: MovimientosPageProps) {
    const params = await searchParams;
    const page = Math.max(Number(params.page) || 1, 1);
    const pageSize = Math.max(Number(params.page_size) || 10, 1);
    const from = params.from || '';
    const to = params.to || '';
    const conceptId = Number(params.concept_id) || null;
    const kind = params.kind === 'INCOME' || params.kind === 'EXPENSE' ? (params.kind as ConceptKind) : null;

    const [entriesResult, conceptsResult] = await Promise.all([
        getAccountingEntriesAction({
            page,
            page_size: pageSize,
            from: from || undefined,
            to: to || undefined,
            concept_id: conceptId || undefined,
            kind: kind || undefined,
        }),
        getAccountingConceptsAction(),
    ]);

    const entries = entriesResult.success
        ? entriesResult.data
        : { data: [], total: 0, page, page_size: pageSize, total_pages: 0 };
    const concepts = conceptsResult.success ? conceptsResult.data : [];
    const error = !entriesResult.success
        ? entriesResult.error
        : !conceptsResult.success
            ? conceptsResult.error
            : null;

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <AccountingGate>
                <EntriesView
                    entries={entries}
                    concepts={concepts}
                    filters={{ page, page_size: pageSize, from, to, concept_id: conceptId, kind }}
                    error={error}
                />
            </AccountingGate>
        </div>
    );
}
