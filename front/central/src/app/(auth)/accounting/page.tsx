import { AccountingGate, ReportView } from '@/services/modules/accounting/ui';
import { getAccountingReportAction } from '@/services/modules/accounting/infra/actions';
import { firstDayOfCurrentMonth, lastDayOfCurrentMonth } from '@/services/modules/accounting/ui/format';

interface AccountingPageProps {
    searchParams: Promise<{ from?: string; to?: string }>;
}

export default async function AccountingPage({ searchParams }: AccountingPageProps) {
    const params = await searchParams;
    const from = params.from || firstDayOfCurrentMonth();
    const to = params.to || lastDayOfCurrentMonth();

    const result = await getAccountingReportAction(from, to);
    const report = result.success ? result.data : null;
    const error = result.success ? null : result.error;

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <AccountingGate>
                <ReportView report={report} error={error} from={from} to={to} />
            </AccountingGate>
        </div>
    );
}
