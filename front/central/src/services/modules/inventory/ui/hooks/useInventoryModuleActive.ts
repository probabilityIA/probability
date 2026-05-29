'use client';

import { useEffect, useState } from 'react';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import type { Integration } from '@/services/integrations/core/domain/types';

interface UseInventoryModuleActiveResult {
    isActive: boolean;
    loading: boolean;
}

export function useInventoryModuleActive(businessId?: number | null): UseInventoryModuleActiveResult {
    const [isActive, setIsActive] = useState(true);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!businessId) {
            setLoading(false);
            setIsActive(true);
            return;
        }

        let cancelled = false;
        setLoading(true);
        getIntegrationsAction({ business_id: businessId, page_size: 100 } as any)
            .then(res => {
                if (cancelled) return;
                if (res?.success && Array.isArray(res.data)) {
                    const inventory = (res.data as Integration[]).find(
                        i => i.integration_type?.code === 'inventory',
                    );
                    setIsActive(Boolean(inventory?.is_active));
                } else {
                    setIsActive(false);
                }
            })
            .catch(() => {
                if (!cancelled) setIsActive(false);
            })
            .finally(() => {
                if (!cancelled) setLoading(false);
            });

        return () => {
            cancelled = true;
        };
    }, [businessId]);

    return { isActive, loading };
}
