'use client';

import ShipmentList from '@/services/modules/shipments/ui/components/ShipmentList';

export default function ShipmentsPage() {
    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Envíos</h1>
                    <p className="text-sm text-gray-500">
                        Gestiona y rastrea los envíos de tus órdenes
                    </p>
                </div>
            </div>

            <ShipmentList />
        </div>
    );
}
