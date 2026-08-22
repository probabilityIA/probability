'use client';

import type { Integration } from '@/services/integrations/core/domain/types';
import { SyncActivityProvider } from '../sync-activity-context';
import { InventoryCompareModal } from './InventoryCompareModal';

interface InventoryCompareStandaloneProps {
    integration: Integration;
    businessId: number | null;
    onClose: () => void;
}

export function InventoryCompareStandalone({ integration, businessId, onClose }: InventoryCompareStandaloneProps) {
    return (
        <SyncActivityProvider
            integrations={[integration]}
            businessId={businessId}
            view="diagrama"
            onViewChange={() => { }}
        >
            <InventoryCompareModal
                integration={integration}
                businessId={businessId}
                onClose={onClose}
            />
        </SyncActivityProvider>
    );
}
