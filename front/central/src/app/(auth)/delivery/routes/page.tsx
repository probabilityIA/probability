'use client';

import { RouteManager } from '@/services/modules/routes/ui';
import { useDeliveryBusiness } from '@/shared/contexts/delivery-business-context';

export default function RoutesPage() {
    const { selectedBusinessId } = useDeliveryBusiness();
    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <RouteManager selectedBusinessId={selectedBusinessId} />
        </div>
    );
}
