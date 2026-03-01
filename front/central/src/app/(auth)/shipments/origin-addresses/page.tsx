'use client';

import React from 'react';
import { OriginAddressManager } from '@/services/modules/shipments/ui/components/OriginAddressManager';
import { useOrdersBusiness } from '@/shared/contexts/orders-business-context';

export default function OriginAddressesPage() {
    const { selectedBusinessId } = useOrdersBusiness();

    return (
        <div className="container mx-auto py-8">
            <OriginAddressManager selectedBusinessId={selectedBusinessId} />
        </div>
    );
}
