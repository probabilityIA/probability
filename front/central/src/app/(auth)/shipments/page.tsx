'use client';

import ShipmentList from '@/services/modules/shipments/ui/components/ShipmentList';

export default function ShipmentsPage() {
    return (
        <div className="flex flex-col h-full p-6">
            <div className="flex-1 min-h-0">
                <ShipmentList />
            </div>
        </div>
    );
}
