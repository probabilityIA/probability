import React from 'react';
import { OriginAddressManager } from '@/services/modules/shipments/ui/components/OriginAddressManager';

export default function OriginAddressesPage() {
    return (
        <div className="container mx-auto py-8">
            <OriginAddressManager />
        </div>
    );
}
