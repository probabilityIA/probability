'use client';

import ShipmentList from '@/services/modules/shipments/ui/components/ShipmentList';

export default function ShipmentsPage() {
    return (
        <div className="flex flex-col h-full">
            <div className="mb-4 flex-shrink-0">
                <h1 className="text-2xl font-bold text-gray-900">Envíos</h1>
                <p className="text-sm text-gray-500 mt-1">
                    Gestiona y rastrea los envíos de tus órdenes
                </p>
            </div>
            <div className="flex-1 min-h-0">
                <ShipmentList />
            </div>
        </div>
    );
}
