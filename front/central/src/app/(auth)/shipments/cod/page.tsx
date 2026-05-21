'use client';

import CodReportView from '@/services/modules/codreport/ui/components/CodReportView';
import { useOrdersBusiness } from '@/shared/contexts/orders-business-context';

export default function CODShipmentsPage() {
    const { selectedBusinessId } = useOrdersBusiness();

    return (
        <div className="flex flex-col h-full p-6">
            <div className="flex-1 min-h-0">
                <CodReportView selectedBusinessId={selectedBusinessId} />
            </div>
        </div>
    );
}
