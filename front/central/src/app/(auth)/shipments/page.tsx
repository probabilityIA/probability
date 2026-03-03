'use client';

import ShipmentList from '@/services/modules/shipments/ui/components/ShipmentList';
import { useOrdersBusiness } from '@/shared/contexts/orders-business-context';

export default function ShipmentsPage() {
    const { selectedBusinessId } = useOrdersBusiness();

    return (
        <div className="flex flex-col h-full p-6">
            <div className="flex-1 min-h-0">
                <ShipmentList selectedBusinessId={selectedBusinessId} />
            </div>
        </div>
    );
}
