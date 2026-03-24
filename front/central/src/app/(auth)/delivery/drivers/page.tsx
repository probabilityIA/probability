'use client';

import { DriverManager } from '@/services/modules/drivers/ui';
import { useDeliveryBusiness } from '@/shared/contexts/delivery-business-context';

export default function DriversPage() {
    const { selectedBusinessId } = useDeliveryBusiness();
    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <DriverManager selectedBusinessId={selectedBusinessId} />
        </div>
    );
}
