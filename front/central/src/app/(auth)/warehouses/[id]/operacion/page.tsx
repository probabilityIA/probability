'use client';

import { useParams } from 'next/navigation';
import { Alert } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import WarehouseOperationView from '@/services/modules/warehouses/ui/components/WarehouseOperationView';

export default function WarehouseOperationPage() {
    const params = useParams<{ id: string }>();
    const warehouseId = Number(params?.id);
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    if (isSuperAdmin && selectedBusinessId === null) {
        return (
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900 px-6 py-8">
                <Alert type="info">Selecciona un negocio para ver la vista operativa.</Alert>
            </div>
        );
    }

    return <WarehouseOperationView warehouseId={warehouseId} businessId={businessId} />;
}
