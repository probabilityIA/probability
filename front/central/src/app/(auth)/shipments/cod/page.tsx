'use client';

import CODShipmentList from '@/services/modules/shipments/ui/components/CODShipmentList';
import { useOrdersBusiness } from '@/shared/contexts/orders-business-context';

export default function CODShipmentsPage() {
    const { selectedBusinessId } = useOrdersBusiness();

    return (
        <div className="flex flex-col h-full p-6">
            <div className="flex-1 min-h-0">
                <CODShipmentList selectedBusinessId={selectedBusinessId} />
            </div>
        </div>
    );
}
