'use client';

import QuotesView from '@/services/modules/shipping-quotes/ui/components/QuotesView';
import { useOrdersBusiness } from '@/shared/contexts/orders-business-context';

export default function ShipmentQuotesPage() {
    const { selectedBusinessId } = useOrdersBusiness();

    return (
        <div className="flex flex-col h-full p-6 overflow-y-auto">
            <QuotesView selectedBusinessId={selectedBusinessId} />
        </div>
    );
}
