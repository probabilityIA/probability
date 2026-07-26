'use client';

import { CustomerManager } from '@/services/modules/customers/ui';
import { useSelectedBusiness } from '@/shared/contexts/selected-business-context';

export default function CustomersPage() {
    const { selectedBusinessId, setSelectedBusinessId } = useSelectedBusiness();

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <CustomerManager selectedBusinessId={selectedBusinessId} onBusinessChange={setSelectedBusinessId} />
        </div>
    );
}
